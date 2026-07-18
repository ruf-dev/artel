import {Button, Input} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/components/ScriptBody/components/SchemaPropsEditor/SchemaPropsEditor.module.css"
import {SchemaProperty} from "@/processes/Tracts.ts"
import {CloseIcon} from "@/pages/tract-canvas/icons/CloseIcon/CloseIcon.tsx"

const PRIMITIVE_TYPES: SchemaProperty["type"][] = ["string", "number", "integer", "boolean"]

interface Props {
    properties: Record<string, SchemaProperty>
    onChange: (properties: Record<string, SchemaProperty>) => void
}

function uniquePropName(base: string, taken: Set<string>): string {
    if (!taken.has(base)) return base
    let i = 2
    while (taken.has(`${base}${i}`)) i++
    return `${base}${i}`
}

/** SchemaPropsEditor edits a flat, one-level object schema's fields — each field's own type is
 * always a primitive here, since a field's array/object shape isn't itself editable. Shared by
 * ScriptParamRow for both a directly object-typed param and an array param whose item type is
 * object. */
export default function SchemaPropsEditor({properties, onChange}: Props) {
    const entries = Object.entries(properties)

    function addProperty() {
        const name = uniquePropName("field", new Set(entries.map(([n]) => n)))
        onChange({...properties, [name]: {type: "string"}})
    }

    function renameProperty(index: number, name: string) {
        const next = entries.map(([n, p], i) => [i === index ? name : n, p] as [string, SchemaProperty])
        onChange(Object.fromEntries(next))
    }

    function retypeProperty(index: number, type: SchemaProperty["type"]) {
        const next = entries.map(([n, p], i) => [n, i === index ? {...p, type} : p] as [string, SchemaProperty])
        onChange(Object.fromEntries(next))
    }

    function removeProperty(index: number) {
        onChange(Object.fromEntries(entries.filter((_, i) => i !== index)))
    }

    return (
        <div className={cls.SchemaPropsEditorContainer}>
            {entries.map(([name, prop], i) => (
                <div className={cls.Row} key={i}>
                    <Input
                        className={cls.NameInputWrapper}
                        inputClassName={cls.NameInput}
                        value={name}
                        setValue={v => renameProperty(i, v)}
                    />
                    <select
                        className={cls.TypeSelect}
                        value={prop.type}
                        onChange={e => retypeProperty(i, e.target.value as SchemaProperty["type"])}
                    >
                        {PRIMITIVE_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
                    </select>
                    <Button variant="iconDanger" onClick={() => removeProperty(i)} aria-label={`Remove ${name}`}>
                        <CloseIcon/>
                    </Button>
                </div>
            ))}
            <Button variant="ghost" className={cls.AddBtn} onClick={addProperty}>
                + Add field
            </Button>
        </div>
    )
}
