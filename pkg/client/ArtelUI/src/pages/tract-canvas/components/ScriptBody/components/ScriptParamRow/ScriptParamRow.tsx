import {Button, Input} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/components/ScriptBody/components/ScriptParamRow/ScriptParamRow.module.css"
import {ScriptParam} from "@/processes/Tracts.ts"
import TemplateInput, {TemplateSource} from "@/components/TemplateInput/TemplateInput.tsx"
import {CloseIcon} from "@/pages/tract-canvas/icons/CloseIcon/CloseIcon.tsx"

const PARAM_TYPES: ScriptParam["type"]["type"][] = ["string", "number", "integer", "boolean", "array", "object"]

interface Props {
    param: ScriptParam
    onChange: (patch: Partial<ScriptParam>) => void
    onRemove: () => void
    /** Present for an input param row — renders a TemplateInput binding this param's value
     * to a template expression. Absent for output params, which are declared only. */
    binding?: {
        value: string
        onChange: (value: string) => void
        sources: TemplateSource[]
    }
}

/** ScriptParamRow renders one editable {name, type, description-less} row shared by
 * ScriptBody's Inputs and Outputs sections — the only difference is whether `binding` is
 * passed (input params bind to a template expr; output params are declared only). */
export default function ScriptParamRow({param, onChange, onRemove, binding}: Props) {
    return (
        <div className={cls.ScriptParamRowContainer}>
            <div className={cls.Head}>
                <Input
                    className={cls.NameInputWrapper}
                    inputClassName={cls.NameInput}
                    value={param.name}
                    setValue={v => onChange({name: v})}
                />
                <select
                    className={cls.TypeSelect}
                    value={param.type.type}
                    onChange={e => {
                        const nextType = e.target.value as ScriptParam["type"]["type"]
                        onChange({type: {...param.type, type: nextType}})
                    }}
                >
                    {PARAM_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
                </select>
                <Button variant="iconDanger" onClick={onRemove} aria-label={`Remove ${param.name}`}>
                    <CloseIcon/>
                </Button>
            </div>
            {binding && (
                <TemplateInput
                    value={binding.value}
                    onChange={binding.onChange}
                    sources={binding.sources}
                    placeholder={`bind ${param.name}`}
                />
            )}
        </div>
    )
}
