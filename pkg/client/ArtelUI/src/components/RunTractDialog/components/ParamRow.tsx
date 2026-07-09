import cls from "@/components/RunTractDialog/RunTractDialog.module.css"
import {SchemaProperty} from "@/processes/Tracts.ts"
import ParamInput from "@/components/RunTractDialog/components/ParamInput.tsx"

interface Props {
    name: string
    def: SchemaProperty
    required: boolean
    value: string
    onChange: (name: string, value: string) => void
}

export default function ParamRow({name, def, required, value, onChange}: Props) {
    return (
        <div className={cls.ParamRow}>
            <span className={cls.ParamName}>{name}{required ? " *" : ""}</span>
            {def.description && <span className={cls.ParamDesc}>{def.description}</span>}
            <ParamInput def={def} value={value} onChange={v => onChange(name, v)}/>
        </div>
    )
}
