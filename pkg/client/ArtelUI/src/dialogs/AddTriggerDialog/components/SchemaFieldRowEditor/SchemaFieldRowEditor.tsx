import {Button} from "@vervstack/chures"
import cls from "@/dialogs/AddTriggerDialog/components/SchemaFieldRowEditor/SchemaFieldRowEditor.module.css"

import {cn} from "@/app/utils/cn.ts"
import {SchemaProperty} from "@/processes/Tracts.ts"
import {emptySchemaField, FIELD_TYPES, SchemaFieldRow} from "@/dialogs/AddTriggerDialog/AddTriggerDialogContext.ts"
import SchemaFieldList from "@/dialogs/AddTriggerDialog/components/SchemaFieldList/SchemaFieldList.tsx"

export default function SchemaFieldRowEditor({field, onChange, onRemove, hideName}: {
    field: SchemaFieldRow
    onChange: (patch: Partial<SchemaFieldRow>) => void
    onRemove?: () => void
    hideName?: boolean
}) {
    return (
        <div className={cls.SchemaFieldRowEditorContainer}>
            <div className={cls.SchemaFieldRow}>
                {!hideName && (
                    <input
                        className={cn(cls.TextInput, cls.SchemaFieldName)}
                        placeholder="field name"
                        value={field.name}
                        onChange={e => onChange({name: e.target.value})}
                    />
                )}
                <select
                    className={cn(cls.PlainSelect, cls.SchemaFieldTypeSelect)}
                    value={field.type}
                    onChange={e => onChange({type: e.target.value as SchemaProperty["type"]})}
                >
                    {FIELD_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
                </select>
                {onRemove && (
                    <Button variant="iconDanger" className={cls.RemoveRowBtn} onClick={onRemove} aria-label="Remove field">✕</Button>
                )}
            </div>
            {field.type === "object" && (
                <div className={cls.NestedFields}>
                    <SchemaFieldList fields={field.properties ?? []} onChange={properties => onChange({properties})}/>
                </div>
            )}
            {field.type === "array" && (
                <div className={cls.NestedFields}>
                    <SchemaFieldRowEditor
                        field={field.items ?? emptySchemaField()}
                        onChange={patch => onChange({items: {...(field.items ?? emptySchemaField()), ...patch}})}
                        hideName
                    />
                </div>
            )}
        </div>
    )
}
