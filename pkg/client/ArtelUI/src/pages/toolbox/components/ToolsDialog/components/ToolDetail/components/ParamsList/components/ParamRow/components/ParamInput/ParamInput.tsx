import {Input} from "@vervstack/chures"

import {ToolParamDef} from "@/app/api/artel/mcp_keys.pb.ts"
import cls from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/ParamsList/components/ParamRow/components/ParamInput/ParamInput.module.css"

export default function ParamInput({def, value, onChange}: {
    def: ToolParamDef
    value: string
    onChange: (value: string) => void
}) {
    if (def.enumParam) {
        return (
            <select className={cls.ParamField} value={value} onChange={e => onChange(e.target.value)}>
                <option value="">—</option>
                {(def.enumParam.values ?? []).map(v => (
                    <option key={v} value={v}>{v}</option>
                ))}
            </select>
        )
    }
    if (def.integerParam) {
        return <Input className={cls.ParamField} type="number" value={value} setValue={onChange} placeholder="integer"/>
    }
    return <Input className={cls.ParamField} type="text" value={value} setValue={onChange} placeholder="string"/>
}
