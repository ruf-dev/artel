import cls from "@/components/RunTractDialog/RunTractDialog.module.css"
import {SchemaNode} from "@/processes/Tracts.ts"
import ParamRow from "@/components/RunTractDialog/components/ParamRow.tsx"

interface Props {
    schema: SchemaNode
    values: Record<string, string>
    onChange: (name: string, value: string) => void
}

export default function ParamsList({schema, values, onChange}: Props) {
    const entries = Object.entries(schema.properties)
    if (entries.length === 0) return null
    return (
        <div className={cls.ParamsSection}>
            {entries.map(([name, def]) => (
                <ParamRow
                    key={name}
                    name={name}
                    def={def}
                    required={schema.required?.includes(name) ?? false}
                    value={values[name] ?? ""}
                    onChange={onChange}
                />
            ))}
        </div>
    )
}
