import TemplateInput, {TemplateSource} from "@/components/TemplateInput/TemplateInput.tsx"
import cls from "@/components/TractStepTree/TractStepTree.module.css"
import {SchemaProperty} from "@/processes/Tracts.ts"

interface Props {
    name: string
    def: SchemaProperty
    required: boolean
    value: string
    sources: TemplateSource[]
    onChange: (name: string, value: string) => void
}

export default function ParamRow({name, def, required, value, sources, onChange}: Props) {
    return (
        <div className={cls.ParamRow}>
            <span className={cls.ParamName}>{name}{required ? " *" : ""}</span>
            {def.description && <span className={cls.ParamDesc}>{def.description}</span>}
            {def.enum?.length ? (
                <select className={cls.PlainSelect} value={value} onChange={e => onChange(name, e.target.value)}>
                    <option value="">—</option>
                    {def.enum.map(v => <option key={v} value={v}>{v}</option>)}
                </select>
            ) : (
                <TemplateInput value={value} onChange={v => onChange(name, v)} sources={sources} placeholder={def.type}/>
            )}
        </div>
    )
}
