import {Input, Toggle} from "@vervstack/chures"

import cls from "@/components/RunTractDialog/RunTractDialog.module.css"
import {SchemaProperty} from "@/processes/Tracts.ts"

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
        return <Input inputClassName={cls.ParamField} type="number" value={value} setValue={onChange}/>
    }
    if (def.type === "boolean") {
        return <Toggle checked={value === "true"} onChange={v => onChange(String(v))}/>
    }
    return <Input inputClassName={cls.ParamField} type="text" value={value} setValue={onChange}/>
}
