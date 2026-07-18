import {Button, Input} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/components/ScriptBody/components/ScriptParamRow/ScriptParamRow.module.css"
import {ScriptParam, SchemaProperty} from "@/processes/Tracts.ts"
import TemplateInput, {TemplateSource} from "@/components/TemplateInput/TemplateInput.tsx"
import {CloseIcon} from "@/pages/tract-canvas/icons/CloseIcon/CloseIcon.tsx"
import SchemaPropsEditor
    from "@/pages/tract-canvas/components/ScriptBody/components/SchemaPropsEditor/SchemaPropsEditor.tsx"

type ParamType = ScriptParam["type"]["type"]

const PARAM_TYPES: ParamType[] = ["string", "number", "integer", "boolean", "array", "object"]
const ITEM_TYPES: ParamType[] = ["string", "number", "integer", "boolean", "object"]

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

/** ScriptParamRow renders one editable {name, type, shape} row shared by ScriptBody's Inputs
 * and Outputs sections — the only difference is whether `binding` is passed (input params
 * bind to a template expr; output params are declared only). When type is "array", an item
 * type selector appears; when the type (or array item type) is "object", a one-level
 * SchemaPropsEditor appears beneath it to declare that object's fields. */
export default function ScriptParamRow({param, onChange, onRemove, binding}: Props) {
    const type = param.type

    function retype(nextType: ParamType) {
        const next: SchemaProperty = {type: nextType, description: type.description}
        if (nextType === "array") next.items = type.items ?? {type: "string"}
        if (nextType === "object") next.properties = type.properties ?? {}
        onChange({type: next})
    }

    function retypeItems(nextItemType: ParamType) {
        const items: SchemaProperty = nextItemType === "object"
            ? {type: "object", properties: type.items?.properties ?? {}}
            : {type: nextItemType}
        onChange({type: {...type, items}})
    }

    function changeItemProperties(properties: Record<string, SchemaProperty>) {
        onChange({type: {...type, items: {...type.items, type: "object", properties}}})
    }

    function changeProperties(properties: Record<string, SchemaProperty>) {
        onChange({type: {...type, properties}})
    }

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
                    value={type.type}
                    onChange={e => retype(e.target.value as ParamType)}
                >
                    {PARAM_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
                </select>
                <Button variant="iconDanger" onClick={onRemove} aria-label={`Remove ${param.name}`}>
                    <CloseIcon/>
                </Button>
            </div>
            {type.type === "array" && (
                <div className={cls.Head}>
                    <span className={cls.ItemTypeLabel}>item type</span>
                    <select
                        className={cls.TypeSelect}
                        value={type.items?.type ?? "string"}
                        onChange={e => retypeItems(e.target.value as ParamType)}
                    >
                        {ITEM_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
                    </select>
                </div>
            )}
            {type.type === "array" && type.items?.type === "object" && (
                <SchemaPropsEditor
                    properties={type.items.properties ?? {}}
                    onChange={changeItemProperties}
                />
            )}
            {type.type === "object" && (
                <SchemaPropsEditor
                    properties={type.properties ?? {}}
                    onChange={changeProperties}
                />
            )}
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
