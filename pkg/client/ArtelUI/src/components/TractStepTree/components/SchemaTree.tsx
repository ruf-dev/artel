import SchemaFieldRow from "@/components/TractStepTree/components/SchemaFieldRow.tsx"
import cls from "@/components/TractStepTree/TractStepTree.module.css"
import {SchemaNode} from "@/processes/Tracts.ts"

interface Props {
    schema: SchemaNode
}

export default function SchemaTree({schema}: Props) {
    return (
        <div className={cls.SchemaTree}>
            {schema.isArray && schema.items ? (
                <SchemaFieldRow name="[0]" prop={schema.items} depth={0}/>
            ) : (
                Object.entries(schema.properties).map(([name, prop]) => (
                    <SchemaFieldRow key={name} name={name} prop={prop} depth={0}/>
                ))
            )}
        </div>
    )
}
