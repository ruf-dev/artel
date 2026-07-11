import {type ReactNode} from "react"
import {Button, Input} from "@vervstack/chures"

import cls from "@/dialogs/AddTriggerDialog/components/SchemaFieldRowEditor/SchemaFieldRowEditor.module.css"
import {cn} from "@/app/utils/cn.ts"
import {SchemaProperty} from "@/processes/Tracts.ts"
import {emptySchemaField, FIELD_TYPES, SchemaFieldRow} from "@/dialogs/AddTriggerDialog/addTriggerDialogContext.ts"

export default function SchemaFieldRowEditor({field, onChange, onRemove, hideName, renderNestedList}: {
    field: SchemaFieldRow
    onChange: (patch: Partial<SchemaFieldRow>) => void
    onRemove?: () => void
    hideName?: boolean
    renderNestedList: (fields: SchemaFieldRow[], onChange: (fields: SchemaFieldRow[]) => void) => ReactNode
}) {
    return (
        <div className={cls.SchemaFieldRowEditorContainer}>
            <div className={cls.SchemaFieldRow}>
                {!hideName && (
                    <Input
                        className={cn(cls.TextInput, cls.SchemaFieldName)}
                        placeholder="field name"
                        value={field.name}
                        setValue={(newValue) => onChange({name: newValue})}
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
                    <Button
                        variant="iconDanger"
                        className={cls.RemoveRowBtn}
                        onClick={onRemove}
                        aria-label="Remove field"
                    >
                        ✕
                    </Button>
                )}
            </div>
            {field.type === "object" && (
                <div className={cls.NestedFields}>
                    {renderNestedList(field.properties ?? [], properties => onChange({properties}))}
                </div>
            )}
            {field.type === "array" && (
                <div className={cls.NestedFields}>
                    <SchemaFieldRowEditor
                        field={field.items ?? emptySchemaField()}
                        onChange={patch => onChange({items: {...(field.items ?? emptySchemaField()), ...patch}})}
                        hideName
                        renderNestedList={renderNestedList}
                    />
                </div>
            )}
        </div>
    )
}
