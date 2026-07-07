import {Button} from "@vervstack/chures"

import {emptySchemaField, SchemaFieldRow} from "@/dialogs/AddTriggerDialog/AddTriggerDialogContext.ts"
import SchemaFieldRowEditor from "@/dialogs/AddTriggerDialog/components/SchemaFieldRowEditor/SchemaFieldRowEditor.tsx"

export default function SchemaFieldList({fields, onChange}: { fields: SchemaFieldRow[]; onChange: (fields: SchemaFieldRow[]) => void }) {
    function update(i: number, patch: Partial<SchemaFieldRow>) {
        onChange(fields.map((f, fi) => fi === i ? {...f, ...patch} : f))
    }

    return (
        <>
            {fields.map((f, i) => (
                <SchemaFieldRowEditor
                    key={i}
                    field={f}
                    onChange={patch => update(i, patch)}
                    onRemove={() => onChange(fields.filter((_, fi) => fi !== i))}
                />
            ))}
            <Button variant="ghost" onClick={() => onChange([...fields, emptySchemaField()])}>
                + Add field
            </Button>
        </>
    )
}
