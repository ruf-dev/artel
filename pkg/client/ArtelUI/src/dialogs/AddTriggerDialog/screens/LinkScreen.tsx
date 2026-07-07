import {Button} from "@vervstack/chures"
import cls from "@/dialogs/AddTriggerDialog/AddTriggerDialog.module.css"

import {cn} from "@/app/utils/cn.ts"
import {TractCondition, tractsService} from "@/processes/Tracts.ts"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import TemplateInput from "@/components/TemplateInput/TemplateInput.tsx"
import {useAddTriggerDialog} from "@/dialogs/AddTriggerDialog/AddTriggerDialogContext.ts"
import DialogHeaderWithClose from "@/dialogs/AddTriggerDialog/components/DialogHeaderWithClose.tsx"

export default function LinkScreen() {
    const {CloseDialog} = useDialog()
    const {triggers} = useTracts()
    const bakeError = useBakeError()

    const context = useAddTriggerDialog()

    const linkable = triggers.filter(t => !context.linkedUuids.has(t.uuid))
    const selected = triggers.find(t => t.uuid === context.selectedTriggerId)
    const filterSources = [{id: "trigger", label: "trigger", schema: selected?.payloadSchema}]

    function handleLink() {
        if (!context.selectedTriggerId) return
        context.setSaving(true)
        tractsService.linkTrigger(context.selectedTriggerId, context.tractUuid, context.filters)
            .then(() => useTracts.getState().fetch())
            .then(CloseDialog)
            .catch(err => bakeError("Failed to link trigger", err))
            .finally(() => context.setSaving(false))
    }

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <DialogHeaderWithClose title="Link existing trigger"/>
            <div className={cls.Body}>
                {linkable.map(t => (
                    <SelectOption key={t.uuid} label={`${t.name} (${t.kind}/${t.source})`} selected={t.uuid === context.selectedTriggerId} onSelect={() => context.setSelectedTriggerId(t.uuid)}/>
                ))}
                {linkable.length === 0 && <p className={cls.Notice}>No other triggers available — create a new one instead.</p>}
                {selected && (
                    <div className={cls.Field}>
                        <span className={cls.FieldLabel}>Filters (AND, optional — gate which deliveries start this tract)</span>
                        {context.filters.map((f, i) => (
                            <div className={cls.SchemaFieldRow} key={i}>
                                <TemplateInput value={f.left} onChange={v => context.setFilters(fs => fs.map((x, xi) => xi === i ? {...x, left: v} : x))} sources={filterSources} placeholder="left"/>
                                <select className={cls.PlainSelect} value={f.op} onChange={e => context.setFilters(fs => fs.map((x, xi) => xi === i ? {...x, op: e.target.value as TractCondition["op"]} : x))}>
                                    {["==", "!=", ">", "<", ">=", "<=", "contains", "glob", "regex"].map(op => <option key={op} value={op}>{op}</option>)}
                                </select>
                                <TemplateInput value={f.right} onChange={v => context.setFilters(fs => fs.map((x, xi) => xi === i ? {...x, right: v} : x))} sources={filterSources} placeholder="right"/>
                                <Button variant="iconDanger" className={cls.RemoveRowBtn} onClick={() => context.setFilters(fs => fs.filter((_, xi) => xi !== i))} aria-label="Remove filter">✕</Button>
                            </div>
                        ))}
                        <Button variant="ghost" onClick={() => context.setFilters(fs => [...fs, {left: "", op: "==", right: ""}])}>+ Add filter</Button>
                    </div>
                )}
            </div>
            <div className={cn(cls.DialogActions, cls.DialogActionsEnd)}>
                <Button variant="ghost" onClick={() => context.setStep("mode")} disabled={context.saving}>Back</Button>
                <Button variant="primary" onClick={handleLink} disabled={context.saving || !context.selectedTriggerId}>
                    {context.saving ? "Linking…" : "Link"}
                </Button>
            </div>
        </div>
    )
}
