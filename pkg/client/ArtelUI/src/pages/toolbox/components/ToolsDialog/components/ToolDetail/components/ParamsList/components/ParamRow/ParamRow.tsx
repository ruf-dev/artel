import {ToolParamDef} from "@/app/api/artel/mcp_keys.pb.ts"
import ParamInput from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/ParamsList/components/ParamRow/components/ParamInput/ParamInput.tsx"
import cls from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/ParamsList/components/ParamRow/ParamRow.module.css"

export default function ParamRow({name, def, value, onChange}: {
    name: string
    def: ToolParamDef
    value: string
    onChange: (name: string, value: string) => void
}) {
    return (
        <div className={cls.ParamRowContainer}>
            <span className={cls.ParamName}>{name}</span>
            {def.description && <span className={cls.ParamDesc}>{def.description}</span>}
            <ParamInput def={def} value={value} onChange={v => onChange(name, v)}/>
        </div>
    )
}
