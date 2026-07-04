import {useEffect, useState} from "react"

import cls from "@/components/TriggerPanel/TriggerPanel.module.css"

import {SchemaNode, SchemaProperty, TractCondition, Trigger} from "@/processes/Tracts.ts"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

import Button from "@/components/shared/Button/Button.tsx"
import ConfirmDialog from "@/components/ConfirmDialog/ConfirmDialog.tsx"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import TemplateInput from "@/components/TemplateInput/TemplateInput.tsx"

interface Props {
    tractUuid: string
    linkedTriggerSummaries: { uuid: string; name: string; kind: string; source: string }[]
}

export default function TriggerPanel({tractUuid, linkedTriggerSummaries}: Props) {
    const {triggers, fetchTriggers, fetchTriggerSources, unlinkTrigger, setTriggerEnabled, rotateTriggerToken} = useTracts()
    const {OpenDialog} = useDialog()
    const bakeError = useBakeError()

    useEffect(() => {
        void fetchTriggers()
        void fetchTriggerSources()
    }, [fetchTriggers, fetchTriggerSources])

    const linkedUuids = new Set(linkedTriggerSummaries.map(t => t.uuid))
    const linked = triggers.filter(t => linkedUuids.has(t.uuid))

    function handleUnlink(triggerUuid: string) {
        OpenDialog(
            <ConfirmDialog
                title="Unlink trigger"
                message="This tract will no longer start on this trigger's events."
                confirmLabel="Unlink"
                danger
                onConfirm={() => unlinkTrigger(triggerUuid, tractUuid).catch(err => bakeError("Failed to unlink trigger", err))}
            />
        )
    }

    function handleRotate(triggerUuid: string) {
        OpenDialog(
            <ConfirmDialog
                title="Rotate webhook token"
                message="The old webhook URL stops working immediately."
                confirmLabel="Rotate"
                danger
                onConfirm={() => rotateTriggerToken(triggerUuid)
                    .then(result => {
                        setTimeout(() => OpenDialog(<TokenRevealDialog webhookUrl={result.webhookUrl} webhookToken={result.webhookToken}/>), 0)
                    })
                    .catch(err => bakeError("Failed to rotate token", err))}
            />
        )
    }

    return (
        <div className={cls.Panel}>
            <div className={cls.PanelHeader}>
                <span className={cls.PanelTitle}>Triggers</span>
                <Button variant="ghost" onClick={() => OpenDialog(<AddTriggerDialog tractUuid={tractUuid} linkedUuids={linkedUuids}/>)}>
                    + Add trigger
                </Button>
            </div>
            {linked.length === 0 && (
                <p className={cls.Empty}>No triggers linked — use "Run" in the Runs panel to fire manually, or add a trigger.</p>
            )}
            {linked.map(t => (
                <TriggerRow
                    key={t.uuid}
                    trigger={t}
                    onUnlink={() => handleUnlink(t.uuid)}
                    onToggle={enabled => setTriggerEnabled(t.uuid, enabled).catch(err => bakeError("Failed to update trigger", err))}
                    onRotate={() => handleRotate(t.uuid)}
                />
            ))}
        </div>
    )
}

function TriggerRow({trigger, onUnlink, onToggle, onRotate}: {
    trigger: Trigger
    onUnlink: () => void
    onToggle: (enabled: boolean) => void
    onRotate: () => void
}) {
    const [copied, setCopied] = useState(false)
    const webhookUrl = trigger.kind === "webhook" ? `${window.location.origin}/tract/hook/${trigger.triggerUuid}` : ""

    function handleCopy() {
        navigator.clipboard.writeText(webhookUrl).then(() => {
            setCopied(true)
            setTimeout(() => setCopied(false), 2000)
        })
    }

    return (
        <div className={cls.Row}>
            <span className={cls.Name}>{trigger.name}</span>
            <span className={cls.Chip}>{trigger.kind}</span>
            <span className={cls.Chip}>{trigger.source}</span>
            {trigger.kind === "webhook" && <span className={cls.Url}>{webhookUrl}</span>}
            <div className={cls.RowActions}>
                {trigger.kind === "webhook" && (
                    <>
                        <Button variant="ghost" onClick={handleCopy}>{copied ? "Copied!" : "Copy URL"}</Button>
                        <Button variant="ghost" onClick={onRotate}>Rotate token</Button>
                    </>
                )}
                <Button variant="ghost" onClick={() => onToggle(!trigger.enabled)}>{trigger.enabled ? "Enabled" : "Disabled"}</Button>
                <Button variant="iconDanger" onClick={onUnlink} aria-label="Unlink trigger">✕</Button>
            </div>
        </div>
    )
}

function TokenRevealDialog({webhookUrl, webhookToken}: { webhookUrl: string; webhookToken: string }) {
    const {CloseDialog} = useDialog()
    const [copied, setCopied] = useState(false)

    function handleCopy() {
        navigator.clipboard.writeText(webhookToken).then(() => {
            setCopied(true)
            setTimeout(() => setCopied(false), 2000)
        })
    }

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Webhook token</h2>
            <p className={cls.Notice}>
                This token is shown once. Store it now — send it as <code>X-Tract-Token</code> (or{" "}
                <code>X-Gitlab-Token</code>) on webhook deliveries to <code>{webhookUrl}</code>.
            </p>
            <div className={cls.TokenBox}>{webhookToken}</div>
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={handleCopy}>{copied ? "Copied!" : "Copy token"}</Button>
                <Button variant="primary" onClick={CloseDialog}>Done</Button>
            </div>
        </div>
    )
}

// -- Add trigger wizard --

type AddStep = "mode" | "create" | "link"
type SchemaFieldRow = { name: string; type: SchemaProperty["type"]; description: string; required: boolean }
const FIELD_TYPES: SchemaProperty["type"][] = ["string", "integer", "number", "boolean", "array", "object"]

function AddTriggerDialog({tractUuid, linkedUuids}: { tractUuid: string; linkedUuids: Set<string> }) {
    const {CloseDialog, OpenDialog} = useDialog()
    const {triggers, triggerSources, createTrigger, linkTrigger} = useTracts()
    const bakeError = useBakeError()

    const [step, setStep] = useState<AddStep>("mode")
    const [saving, setSaving] = useState(false)

    // create-trigger fields
    const [name, setName] = useState("")
    const [kind, setKind] = useState<"webhook" | "manual">("webhook")
    const [source, setSource] = useState("generic")
    const [fields, setFields] = useState<SchemaFieldRow[]>([{name: "", type: "string", description: "", required: false}])

    // link-existing fields
    const [selectedTriggerId, setSelectedTriggerId] = useState("")
    const [filters, setFilters] = useState<TractCondition[]>([])

    function buildSchemaFromFields(): SchemaNode {
        const properties: Record<string, SchemaProperty> = {}
        const required: string[] = []
        for (const f of fields) {
            if (!f.name) continue
            properties[f.name] = {type: f.type, description: f.description || undefined}
            if (f.required) required.push(f.name)
        }
        return {properties, required}
    }

    function handleCreate() {
        if (!name.trim()) return
        setSaving(true)
        const payloadSchema = source === "generic" ? buildSchemaFromFields() : {properties: {}}
        createTrigger(name.trim(), kind, kind === "webhook" ? source : "generic", {}, payloadSchema)
            .then(result => linkTrigger(result.trigger.uuid, tractUuid, []).then(() => result))
            .then(result => {
                CloseDialog()
                if (kind === "webhook") {
                    setTimeout(() => OpenDialog(<TokenRevealDialog webhookUrl={result.webhookUrl} webhookToken={result.webhookToken}/>), 0)
                }
            })
            .catch(err => bakeError("Failed to create trigger", err))
            .finally(() => setSaving(false))
    }

    function handleLink() {
        if (!selectedTriggerId) return
        setSaving(true)
        linkTrigger(selectedTriggerId, tractUuid, filters)
            .then(CloseDialog)
            .catch(err => bakeError("Failed to link trigger", err))
            .finally(() => setSaving(false))
    }

    if (step === "create") {
        return (
            <div className={cls.DialogContainer} role="dialog" aria-modal="true">
                <h2 className={cls.DialogTitle}>New trigger</h2>
                <div className={cls.Body}>
                    <label className={cls.Field}>
                        <span className={cls.FieldLabel}>Name</span>
                        <input className={cls.TextInput} value={name} onChange={e => setName(e.target.value)} autoFocus/>
                    </label>
                    <label className={cls.Field}>
                        <span className={cls.FieldLabel}>Kind</span>
                        <KindSelect value={kind} onChange={setKind}/>
                    </label>
                    {kind === "webhook" && (
                        <label className={cls.Field}>
                            <span className={cls.FieldLabel}>Source preset</span>
                            <select className={cls.PlainSelect} value={source} onChange={e => setSource(e.target.value)}>
                                {triggerSources.map(s => <option key={s.key} value={s.key}>{s.key} — {s.description}</option>)}
                                {triggerSources.length === 0 && <option value="generic">generic</option>}
                            </select>
                        </label>
                    )}
                    {(kind === "manual" || source === "generic") && (
                        <SchemaBuilder fields={fields} onChange={setFields}/>
                    )}
                </div>
                <div className={cls.DialogActions}>
                    <Button variant="ghost" onClick={() => setStep("mode")} disabled={saving}>Back</Button>
                    <div className={cls.ActionsRight}>
                        <Button variant="ghost" onClick={CloseDialog} disabled={saving}>Cancel</Button>
                        <Button variant="primary" onClick={handleCreate} disabled={saving || !name.trim()}>
                            {saving ? "Creating…" : "Create & link"}
                        </Button>
                    </div>
                </div>
            </div>
        )
    }

    if (step === "link") {
        const linkable = triggers.filter(t => !linkedUuids.has(t.uuid))
        const selected = triggers.find(t => t.uuid === selectedTriggerId)
        const filterSources = [{id: "trigger", label: "trigger", schema: selected?.payloadSchema}]

        return (
            <div className={cls.DialogContainer} role="dialog" aria-modal="true">
                <h2 className={cls.DialogTitle}>Link existing trigger</h2>
                <div className={cls.Body}>
                    {linkable.map(t => (
                        <SelectOption key={t.uuid} label={`${t.name} (${t.kind}/${t.source})`} selected={t.uuid === selectedTriggerId} onSelect={() => setSelectedTriggerId(t.uuid)}/>
                    ))}
                    {linkable.length === 0 && <p className={cls.Notice}>No other triggers available — create a new one instead.</p>}
                    {selected && (
                        <div className={cls.Field}>
                            <span className={cls.FieldLabel}>Filters (AND, optional — gate which deliveries start this tract)</span>
                            {filters.map((f, i) => (
                                <div className={cls.SchemaFieldRow} key={i}>
                                    <TemplateInput value={f.left} onChange={v => setFilters(fs => fs.map((x, xi) => xi === i ? {...x, left: v} : x))} sources={filterSources} placeholder="left"/>
                                    <select className={cls.PlainSelect} value={f.op} onChange={e => setFilters(fs => fs.map((x, xi) => xi === i ? {...x, op: e.target.value as TractCondition["op"]} : x))}>
                                        {["==", "!=", ">", "<", ">=", "<=", "contains", "glob", "regex"].map(op => <option key={op} value={op}>{op}</option>)}
                                    </select>
                                    <TemplateInput value={f.right} onChange={v => setFilters(fs => fs.map((x, xi) => xi === i ? {...x, right: v} : x))} sources={filterSources} placeholder="right"/>
                                    <Button variant="iconDanger" className={cls.RemoveRowBtn} onClick={() => setFilters(fs => fs.filter((_, xi) => xi !== i))} aria-label="Remove filter">✕</Button>
                                </div>
                            ))}
                            <Button variant="ghost" onClick={() => setFilters(fs => [...fs, {left: "", op: "==", right: ""}])}>+ Add filter</Button>
                        </div>
                    )}
                </div>
                <div className={cls.DialogActions}>
                    <Button variant="ghost" onClick={() => setStep("mode")} disabled={saving}>Back</Button>
                    <div className={cls.ActionsRight}>
                        <Button variant="ghost" onClick={CloseDialog} disabled={saving}>Cancel</Button>
                        <Button variant="primary" onClick={handleLink} disabled={saving || !selectedTriggerId}>
                            {saving ? "Linking…" : "Link"}
                        </Button>
                    </div>
                </div>
            </div>
        )
    }

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Add trigger</h2>
            <div className={cls.Body}>
                <button type="button" className={cls.Row} onClick={() => setStep("create")}>
                    <span className={cls.Name}>Create new trigger</span>
                </button>
                <button type="button" className={cls.Row} onClick={() => setStep("link")}>
                    <span className={cls.Name}>Link existing trigger</span>
                </button>
            </div>
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={CloseDialog}>Cancel</Button>
            </div>
        </div>
    )
}

function KindSelect({value, onChange}: { value: "webhook" | "manual"; onChange: (v: "webhook" | "manual") => void }) {
    return (
        <select className={cls.PlainSelect} value={value} onChange={e => onChange(e.target.value as "webhook" | "manual")}>
            <option value="webhook">Webhook</option>
            <option value="manual">Manual</option>
        </select>
    )
}

function SchemaBuilder({fields, onChange}: { fields: SchemaFieldRow[]; onChange: (fields: SchemaFieldRow[]) => void }) {
    function update(i: number, patch: Partial<SchemaFieldRow>) {
        onChange(fields.map((f, fi) => fi === i ? {...f, ...patch} : f))
    }

    return (
        <div className={cls.Field}>
            <span className={cls.FieldLabel}>Input schema (fields the trigger accepts)</span>
            {fields.map((f, i) => (
                <div className={cls.SchemaFieldRow} key={i}>
                    <input
                        className={`${cls.TextInput} ${cls.SchemaFieldName}`}
                        placeholder="field name"
                        value={f.name}
                        onChange={e => update(i, {name: e.target.value})}
                    />
                    <select
                        className={`${cls.PlainSelect} ${cls.SchemaFieldTypeSelect}`}
                        value={f.type}
                        onChange={e => update(i, {type: e.target.value as SchemaProperty["type"]})}
                    >
                        {FIELD_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
                    </select>
                    <Button variant="iconDanger" className={cls.RemoveRowBtn} onClick={() => onChange(fields.filter((_, fi) => fi !== i))} aria-label="Remove field">✕</Button>
                </div>
            ))}
            <Button variant="ghost" onClick={() => onChange([...fields, {name: "", type: "string", description: "", required: false}])}>
                + Add field
            </Button>
        </div>
    )
}
