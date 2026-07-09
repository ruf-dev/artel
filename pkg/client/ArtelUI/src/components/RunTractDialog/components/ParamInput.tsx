import cls from "@/components/RunTractDialog/RunTractDialog.module.css"
import {SchemaProperty} from "@/processes/Tracts.ts"
import Input from "@/components/shared/Input/Input.tsx"

interface Props {
    def: SchemaProperty
    value: string
    onChange: (value: string) => void
}

export default function ParamInput({def, value, onChange}: Props) {
    if (def.enum?.length) {
        return (
            <select className={cls.ParamField} value={value} onChange={e => onChange(e.target.value)}>
                <option value="">—</option>
                {def.enum.map(v => <option key={v} value={v}>{v}</option>)}
            </select>
        )
    }
    if (def.type === "integer" || def.type === "number") {
        return <Input className={cls.ParamField} type="number" value={value} onChange={e => onChange(e.currentTarget.value)}/>
    }
    if (def.type === "boolean") {
        return <Input type="checkbox" checked={value === "true"} onChange={e => onChange(String(e.currentTarget.checked))}/>
    }
    return <Input className={cls.ParamField} type="text" value={value} onChange={e => onChange(e.currentTarget.value)}/>
}
