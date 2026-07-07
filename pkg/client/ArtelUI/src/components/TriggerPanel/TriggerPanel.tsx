import {useEffect, useRef, useState} from "react"
import type {RefObject} from "react"
import {createPortal} from "react-dom"

import {Button, ConfirmDialog, Dropdown} from "@vervstack/chures"
import type {DropdownOption} from "@vervstack/chures"
import cls from "@/components/TriggerPanel/TriggerPanel.module.css"

import {SchemaNode, SchemaProperty, TractCondition, Trigger, TriggerSource} from "@/processes/Tracts.ts"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"

import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import TemplateInput from "@/components/TemplateInput/TemplateInput.tsx"
import ManageGitlabDialog from "@/dialogs/ManageGitlabDialog/ManageGitlabDialog.tsx"

// -- Category/provider display helpers (shared by the source picker and TriggerRow) --

function categoryLabel(category: string): string {
    switch (category) {
        case "gitlab":
            return "GitLab hooks"
        case "generic":
            return "Generic webhook"
        default:
            return category
    }
}

const PROVIDER_ENUM_BY_KEY: Record<string, ExternalProvider> = {
    gitlab: ExternalProvider.EXTERNAL_PROVIDER_GITLAB,
}

function providerLabel(provider: string): string {
    if (provider === "gitlab") return "GitLab"
    return provider ? provider.charAt(0).toUpperCase() + provider.slice(1) : ""
}

function providerEnumFor(provider: string): ExternalProvider | undefined {
    return PROVIDER_ENUM_BY_KEY[provider]
}

interface Props {
    tractUuid: string
    linkedTriggerSummaries: { uuid: string; name: string; kind: string; source: string }[]
}

export default function TriggerPanel({tractUuid, linkedTriggerSummaries}: Props) {
    const {triggers, triggerSources, fetchTriggers, fetchTriggerSources, unlinkTrigger, setTriggerEnabled, rotateTriggerToken} = useTracts()
    const {OpenDialog, CloseDialog} = useDialog()
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
                onClose={CloseDialog}
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
                onClose={CloseDialog}
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
                <span className={cls.PanelTitle}>Trigger</span>
                {linked.length === 0 && (
                    <Button variant="ghost" onClick={() => OpenDialog(<AddTriggerDialog tractUuid={tractUuid} linkedUuids={linkedUuids}/>)}>
                        + Add trigger
                    </Button>
                )}
            </div>
            {linked.length === 0 && (
                <p className={cls.Empty}>No triggers linked — use "Run" in the Runs panel to fire manually, or add a trigger.</p>
            )}
            {linked.map(t => (
                <TriggerRow
                    key={t.uuid}
                    trigger={t}
                    triggerSources={triggerSources}
                    onUnlink={() => handleUnlink(t.uuid)}
                    onToggle={enabled => setTriggerEnabled(t.uuid, enabled).catch(err => bakeError("Failed to update trigger", err))}
                    onRotate={() => handleRotate(t.uuid)}
                />
            ))}
        </div>
    )
}

function TriggerRow({trigger, triggerSources, onUnlink, onToggle, onRotate}: {
    trigger: Trigger
    triggerSources: TriggerSource[]
    onUnlink: () => void
    onToggle: (enabled: boolean) => void
    onRotate: () => void
}) {
    const [copied, setCopied] = useState(false)
    const preset = triggerSources.find(s => s.key === trigger.source)
    const sharedProvider = trigger.kind === "webhook" ? preset?.provider ?? "" : ""
    const webhookUrl = trigger.kind === "webhook" && !sharedProvider ? `${window.location.origin}/tract/hook/${trigger.triggerUuid}` : ""

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
            {trigger.kind === "webhook" && !sharedProvider && <span className={cls.Url}>{webhookUrl}</span>}
            <div className={cls.RowActions}>
                {trigger.kind === "webhook" && (
                    sharedProvider ? (
                        <span className={cls.Chip}>Shared · {providerLabel(sharedProvider)} connection</span>
                    ) : (
                        <>
                            <Button variant="ghost" onClick={handleCopy}>{copied ? "Copied!" : "Copy URL"}</Button>
                            <Button variant="ghost" onClick={onRotate}>Rotate token</Button>
                        </>
                    )
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
    const {triggers, triggerSources, createTrigger, linkTrigger, fetchTriggers, fetchTriggerSources} = useTracts()
    const {connections} = useExternalConnections()
    const bakeError = useBakeError()

    const [step, setStep] = useState<AddStep>("mode")
    const [saving, setSaving] = useState(false)

    // Refetch on every open — this dialog mounts fresh each time OpenDialog is called
    // (see app/hooks/Dialog.ts), so a plain mount effect refreshes stale presets/triggers
    // instead of reusing whatever TriggerPanel's own mount fetch loaded at page-load time.
    useEffect(() => {
        void fetchTriggerSources()
        void fetchTriggers()
    }, [fetchTriggerSources, fetchTriggers])

    // create-trigger fields
    const [name, setName] = useState("")
    const [kind, setKind] = useState<"webhook" | "manual">("webhook")
    const [source, setSource] = useState("generic")
    const [fields, setFields] = useState<SchemaFieldRow[]>([{name: "", type: "string", description: "", required: false}])

    // link-existing fields
    const [selectedTriggerId, setSelectedTriggerId] = useState("")
    const [filters, setFilters] = useState<TractCondition[]>([])

    const selectedSource = kind === "webhook" ? triggerSources.find(s => s.key === source) : undefined
    const hasProviderConnection = selectedSource?.provider
        ? connections.some(c => c.provider === providerEnumFor(selectedSource.provider))
        : true

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
        const payloadSchema = kind === "webhook" && source !== "generic"
            ? selectedSource?.payloadSchema ?? {properties: {}}
            : buildSchemaFromFields()
        createTrigger(name.trim(), kind, kind === "webhook" ? source : "generic", {}, payloadSchema)
            .then(result => linkTrigger(result.trigger.uuid, tractUuid, []).then(() => result))
            .then(result => {
                CloseDialog()
                if (kind === "webhook" && result.webhookToken) {
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
                <DialogHeaderWithBack title="New trigger" onBack={() => setStep("mode")}/>
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
                        <SourcePicker triggerSources={triggerSources} source={source} onSourceChange={setSource}/>
                    )}
                    {(kind === "manual" || source === "generic") && (
                        <SchemaBuilder fields={fields} onChange={setFields}/>
                    )}
                    {kind === "webhook" && source !== "generic" && selectedSource && (
                        <PresetDetails
                            preset={selectedSource}
                            hasProviderConnection={hasProviderConnection}
                            onConnectProvider={() => OpenDialog(<ManageGitlabDialog/>)}
                        />
                    )}
                </div>
                <div className={`${cls.DialogActions} ${cls.DialogActionsEnd}`}>
                    <Button variant="ghost" onClick={CloseDialog} disabled={saving}>Cancel</Button>
                    <Button variant="primary" onClick={handleCreate} disabled={saving || !name.trim() || !hasProviderConnection}>
                        {saving ? "Creating…" : "Create & link"}
                    </Button>
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
                <DialogHeaderWithBack title="Link existing trigger" onBack={() => setStep("mode")}/>
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
                <div className={`${cls.DialogActions} ${cls.DialogActionsEnd}`}>
                    <Button variant="ghost" onClick={CloseDialog} disabled={saving}>Cancel</Button>
                    <Button variant="primary" onClick={handleLink} disabled={saving || !selectedTriggerId}>
                        {saving ? "Linking…" : "Link"}
                    </Button>
                </div>
            </div>
        )
    }

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Add trigger</h2>
            <div className={cls.Body}>
                <Button variant="ghost" className={cls.OptionRow} onClick={() => setStep("create")}>
                    <span className={cls.Name}>Create new trigger</span>
                </Button>
                <Button variant="ghost" className={cls.OptionRow} onClick={() => setStep("link")}>
                    <span className={cls.Name}>Link existing trigger</span>
                </Button>
            </div>
            <div className={`${cls.DialogActions} ${cls.DialogActionsEnd}`}>
                <Button variant="ghost" onClick={CloseDialog}>Cancel</Button>
            </div>
        </div>
    )
}

function DialogHeaderWithBack({title, onBack}: { title: string; onBack: () => void }) {
    return (
        <div className={cls.DialogHeader}>
            <Button variant="ghost" className={cls.BackBtn} onClick={onBack} aria-label="Back">
                <BackChevronIcon className={cls.BackIcon}/>
            </Button>
            <h2 className={cls.DialogTitle}>{title}</h2>
        </div>
    )
}

function BackChevronIcon({className}: { className?: string }) {
    return (
        <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={1.6} strokeLinecap="round" strokeLinejoin="round">
            <polyline points="15 18 9 12 15 6"/>
        </svg>
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

// Two-level webhook source picker: pick a category first, then (if the category has more
// than one preset) pick a preset within it. Categories/presets with a single option resolve
// immediately without a second pick.
function SourcePicker({triggerSources, source, onSourceChange}: {
    triggerSources: TriggerSource[]
    source: string
    onSourceChange: (source: string) => void
}) {
    const categories = Array.from(new Set(triggerSources.map(s => s.category)))
    const selected = triggerSources.find(s => s.key === source)
    const [pendingCategory, setPendingCategory] = useState<string | null>(null)
    const category = pendingCategory ?? selected?.category ?? ""
    const presetsInCategory = triggerSources.filter(s => s.category === category)

    const [categoryOpen, setCategoryOpen] = useState(false)
    const [presetOpen, setPresetOpen] = useState(false)
    const categoryAnchorRef = useRef<HTMLDivElement>(null)
    const presetAnchorRef = useRef<HTMLDivElement>(null)

    function pickCategory(opt: DropdownOption) {
        const key = typeof opt === "string" ? opt : opt.id
        setCategoryOpen(false)
        const presets = triggerSources.filter(s => s.category === key)
        if (presets.length === 1) {
            setPendingCategory(null)
            onSourceChange(presets[0].key)
        } else {
            setPendingCategory(key)
            onSourceChange("")
        }
    }

    function pickPreset(opt: DropdownOption) {
        const key = typeof opt === "string" ? opt : opt.id
        setPresetOpen(false)
        setPendingCategory(null)
        onSourceChange(key)
    }

    return (
        <>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Category</span>
                <div className={cls.PickerAnchor} ref={categoryAnchorRef}>
                    <Button type="button" variant="secondary" className={cls.PickerButton} onClick={() => setCategoryOpen(o => !o)}>
                        {category ? categoryLabel(category) : "Choose category…"}
                    </Button>
                </div>
                {categoryOpen && (
                    <FloatingDropdown
                        anchorRef={categoryAnchorRef}
                        options={categories.map(c => ({id: c, name: categoryLabel(c)}))}
                        onPick={pickCategory}
                        onClose={() => setCategoryOpen(false)}
                        placeholder="Search categories…"
                    />
                )}
            </label>
            {category && presetsInCategory.length > 1 && (
                <label className={cls.Field}>
                    <span className={cls.FieldLabel}>Preset</span>
                    <div className={cls.PickerAnchor} ref={presetAnchorRef}>
                        <Button type="button" variant="secondary" className={cls.PickerButton} onClick={() => setPresetOpen(o => !o)}>
                            {selected ? selected.label : "Choose preset…"}
                        </Button>
                    </div>
                    {presetOpen && (
                        <FloatingDropdown
                            anchorRef={presetAnchorRef}
                            options={presetsInCategory.map(p => ({id: p.key, name: p.label}))}
                            onPick={pickPreset}
                            onClose={() => setPresetOpen(false)}
                            placeholder="Search presets…"
                        />
                    )}
                </label>
            )}
        </>
    )
}

// Portals a chures <Dropdown> to document.body, positioned under its anchor via
// getBoundingClientRect — same "never use z-index" pattern as TemplateInput.tsx.
function FloatingDropdown({anchorRef, options, onPick, onClose, placeholder}: {
    anchorRef: RefObject<HTMLElement | null>
    options: DropdownOption[]
    onPick: (opt: DropdownOption) => void
    onClose: () => void
    placeholder?: string
}) {
    const [rect, setRect] = useState<{ top: number; left: number; width: number } | null>(null)

    useEffect(() => {
        function reposition() {
            const el = anchorRef.current
            if (!el) return
            const box = el.getBoundingClientRect()
            setRect({top: box.bottom + 4, left: box.left, width: box.width})
        }

        reposition()
        window.addEventListener("scroll", reposition, true)
        window.addEventListener("resize", reposition)
        return () => {
            window.removeEventListener("scroll", reposition, true)
            window.removeEventListener("resize", reposition)
        }
    }, [anchorRef])

    if (!rect) return null

    return createPortal(
        <div className={cls.DropdownPanel} style={{top: rect.top, left: rect.left, width: rect.width}}>
            <Dropdown options={options} onPick={onPick} onClose={onClose} anchorRef={anchorRef} placeholder={placeholder}/>
        </div>,
        document.body,
    )
}

// Read-only preview of a resolved preset: description + its payload schema fields (no inputs),
// plus a provider-connection nudge for provider-linked presets (e.g. gitlab_push).
function PresetDetails({preset, hasProviderConnection, onConnectProvider}: {
    preset: TriggerSource
    hasProviderConnection: boolean
    onConnectProvider: () => void
}) {
    const propertyEntries = Object.entries(preset.payloadSchema.properties)

    return (
        <div className={cls.Field}>
            <span className={cls.FieldLabel}>{preset.label}</span>
            <p className={cls.Notice}>{preset.description}</p>
            {propertyEntries.length > 0 && (
                <div className={cls.SchemaPreview}>
                    {propertyEntries.map(([propName, prop]) => (
                        <div className={cls.SchemaFieldRow} key={propName}>
                            <span className={cls.SchemaFieldName}>{propName}</span>
                            <span className={cls.SchemaFieldTypeSelect}>{prop.type}</span>
                            {prop.description && <span className={cls.Notice}>{prop.description}</span>}
                        </div>
                    ))}
                </div>
            )}
            {preset.provider && (
                hasProviderConnection ? (
                    <p className={cls.Notice}>Uses your connected {providerLabel(preset.provider)}.</p>
                ) : (
                    <div className={cls.ProviderNotice}>
                        <p className={cls.Notice}>Connect {providerLabel(preset.provider)} first to use this trigger.</p>
                        <Button variant="secondary" onClick={onConnectProvider}>Connect {providerLabel(preset.provider)}</Button>
                    </div>
                )
            )}
        </div>
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
