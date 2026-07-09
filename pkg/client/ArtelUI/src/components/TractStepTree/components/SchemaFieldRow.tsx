import cls from "@/components/TractStepTree/TractStepTree.module.css"
import {SchemaProperty} from "@/processes/Tracts.ts"

interface Props {
    name: string
    prop: SchemaProperty
    depth: number
}

export default function SchemaFieldRow({name, prop, depth}: Props) {
    return (
        <>
            <div className={cls.SchemaField} style={{paddingLeft: `${depth * 0.75}rem`}}>
                <span className={cls.SchemaFieldName}>{name}</span>
                <span className={cls.SchemaFieldType}>{prop.type}</span>
            </div>
            {prop.type === "object" && prop.properties && Object.entries(prop.properties).map(([n, p]) => (
                <SchemaFieldRow key={n} name={n} prop={p} depth={depth + 1}/>
            ))}
            {prop.type === "array" && prop.items && (
                <SchemaFieldRow name="[0]" prop={prop.items} depth={depth + 1}/>
            )}
        </>
    )
}
