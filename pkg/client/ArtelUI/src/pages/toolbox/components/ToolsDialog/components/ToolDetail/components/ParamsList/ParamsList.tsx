import {ToolParamDef} from "@/app/api/artel/mcp_keys.pb.ts"
import ParamRow from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/ParamsList/components/ParamRow/ParamRow.tsx"
import cls from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/ParamsList/ParamsList.module.css"

export default function ParamsList({params, values, onChange}: {
    params: Record<string, ToolParamDef>
    values: Record<string, string>
    onChange: (name: string, value: string) => void
}) {
    const entries = Object.entries(params)
    if (entries.length === 0) return null
    return (
        <div className={cls.ParamsListContainer}>
            <span className={cls.SectionTitle}>Parameters</span>
            {entries.map(([name, def]) => (
                <ParamRow key={name} name={name} def={def} value={values[name] ?? ""} onChange={onChange}/>
            ))}
        </div>
    )
}
