import {useState} from "react"
import {Button} from "@vervstack/chures"
import cls from "@/dialogs/AddTriggerDialog/AddTriggerDialog.module.css"

import {SchemaFieldRow, fieldsToSchemaNode} from "@/dialogs/AddTriggerDialog/AddTriggerDialogContext.ts"
import {useTriggerSources} from "@/app/hooks/Tracts.ts"
import CodeIcon from "@/icons/common/CodeIcon.tsx"
import JsonView from "@/dialogs/AddTriggerDialog/widgets/JsonView/JsonView.tsx"
import SchemaFieldList from "@/dialogs/AddTriggerDialog/components/SchemaFieldList/SchemaFieldList.tsx"
import SchemaPropertyPreview from "@/dialogs/AddTriggerDialog/components/SchemaPropertyPreview/SchemaPropertyPreview.tsx"

interface SchemaBuilderProps {
    // Trigger source key whose payload schema should be fetched and previewed read-only.
    // Empty when there is no existing schema to resolve, i.e. the user is building one by hand.
    schemaId: string
    fields: SchemaFieldRow[]
    onChange: (fields: SchemaFieldRow[]) => void
}

export default function SchemaBuilder({schemaId, fields, onChange}: SchemaBuilderProps) {
    const {triggerSources} = useTriggerSources()
    const preset = schemaId ? triggerSources.find(s => s.key === schemaId) : undefined
    const [view, setView] = useState<"form" | "json">("form")

    const schema = preset ? preset.payloadSchema : fieldsToSchemaNode(fields)

    return (
        <div className={cls.Field}>
            <div className={cls.FieldHeader}>
                <span className={cls.FieldLabel}>Input schema</span>
                <Button
                    variant="ghost"
                    className={cls.JsonToggleBtn}
                    onClick={() => setView(v => v === "form" ? "json" : "form")}
                    aria-label={view === "form" ? "View schema as JSON" : "View schema as form"}
                    aria-pressed={view === "json"}
                >
                    <CodeIcon className={cls.JsonToggleIcon}/>
                </Button>
            </div>
            {view === "json" ? (
                <JsonView value={schema}/>
            ) : preset ? (
                Object.entries(preset.payloadSchema.properties).map(([propName, prop]) => (
                    <SchemaPropertyPreview key={propName} name={propName} prop={prop}/>
                ))
            ) : (
                <SchemaFieldList fields={fields} onChange={onChange}/>
            )}
        </div>
    )
}
