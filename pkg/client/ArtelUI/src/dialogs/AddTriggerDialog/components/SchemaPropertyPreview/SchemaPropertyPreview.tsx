import {Button} from "@vervstack/chures"

import cls from "@/dialogs/AddTriggerDialog/components/SchemaPropertyPreview/SchemaPropertyPreview.module.css"
import {cn} from "@/app/utils/cn.ts"
import {SchemaProperty} from "@/processes/Tracts.ts"
import {FIELD_TYPES} from "@/dialogs/AddTriggerDialog/addTriggerDialogContext.ts"

export default function SchemaPropertyPreview({name, prop, hideName}: { name: string; prop: SchemaProperty; hideName?: boolean }) {
    return (
        <div className={cls.SchemaPropertyPreviewContainer}>
            <div className={cls.SchemaFieldRow}>
                {!hideName && <input className={cn(cls.TextInput, cls.SchemaFieldName)} value={name} disabled/>}
                <select className={cn(cls.PlainSelect, cls.SchemaFieldTypeSelect)} value={prop.type} disabled>
                    {FIELD_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
                </select>
                {prop.description && (
                    <Button
                        type="button"
                        variant="ghost"
                        className={cls.SchemaFieldHelp}
                        data-tooltip-id="root-tooltip"
                        data-tooltip-content={prop.description}
                        aria-label="Field description"
                    >?</Button>
                )}
            </div>
            {prop.type === "object" && prop.properties && (
                <div className={cls.NestedFields}>
                    {Object.entries(prop.properties).map(([propName, nested]) => (
                        <SchemaPropertyPreview key={propName} name={propName} prop={nested}/>
                    ))}
                </div>
            )}
            {prop.type === "array" && prop.items && (
                <div className={cls.NestedFields}>
                    <SchemaPropertyPreview name="" prop={prop.items} hideName/>
                </div>
            )}
        </div>
    )
}
