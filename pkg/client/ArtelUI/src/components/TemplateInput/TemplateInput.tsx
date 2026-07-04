import {useMemo, useRef, useState} from "react"

import cls from "@/components/TemplateInput/TemplateInput.module.css"

import {SchemaNode, SchemaProperty} from "@/processes/Tracts.ts"
import {extractRefs, validateRef} from "@/processes/tractTemplate.ts"

import Button from "@/components/shared/Button/Button.tsx"

export interface TemplateSource {
    id: string
    label: string
    schema?: SchemaNode
}

interface FlatEntry {
    key: string
    ref: string
    label: string
    type?: string
    description?: string
    depth: number
}

function flattenProperty(baseRef: string, name: string, prop: SchemaProperty, depth: number, out: FlatEntry[]) {
    const ref = `${baseRef}.${name}`
    out.push({key: ref, ref, label: name, type: prop.type, description: prop.description, depth})

    if (prop.type === "object" && prop.properties) {
        for (const [childName, childProp] of Object.entries(prop.properties)) {
            flattenProperty(ref, childName, childProp, depth + 1, out)
        }
    } else if (prop.type === "array" && prop.items) {
        const arrRef = `${ref}[0]`
        out.push({key: arrRef, ref: arrRef, label: "[0]", type: prop.items.type, description: prop.items.description, depth: depth + 1})
        if (prop.items.type === "object" && prop.items.properties) {
            for (const [childName, childProp] of Object.entries(prop.items.properties)) {
                flattenProperty(arrRef, childName, childProp, depth + 2, out)
            }
        }
    }
}

function flattenSource(source: TemplateSource): FlatEntry[] {
    const out: FlatEntry[] = []
    if (source.schema) {
        for (const [name, prop] of Object.entries(source.schema.properties)) {
            flattenProperty(source.id, name, prop, 1, out)
        }
    }
    return out
}

interface Props {
    value: string
    onChange: (value: string) => void
    sources: TemplateSource[]
    placeholder?: string
}

export default function TemplateInput({value, onChange, sources, placeholder}: Props) {
    const inputRef = useRef<HTMLInputElement>(null)
    const [open, setOpen] = useState(false)
    const [query, setQuery] = useState("")
    const [triggerStart, setTriggerStart] = useState<number | null>(null)
    const [activeIndex, setActiveIndex] = useState(0)

    const sourceIds = useMemo(() => sources.map(s => s.id), [sources])

    const {invalid, malformed} = useMemo(() => {
        try {
            const refs = extractRefs(value)
            return {invalid: refs.filter(r => !validateRef(r, sourceIds)), malformed: false}
        } catch {
            return {invalid: [], malformed: true}
        }
    }, [value, sourceIds])

    const filtered = useMemo(() => {
        const q = query.toLowerCase()
        return sources.map(source => ({
            source,
            entries: flattenSource(source).filter(e => !q || e.ref.toLowerCase().includes(q)),
            rootMatches: !q || source.id.toLowerCase().includes(q),
        })).filter(g => g.rootMatches || g.entries.length > 0)
    }, [sources, query])

    const flatList = useMemo(() => {
        const list: { ref: string; label: string }[] = []
        for (const g of filtered) {
            if (g.rootMatches) list.push({ref: g.source.id, label: g.source.label})
            for (const e of g.entries) list.push({ref: e.ref, label: e.label})
        }
        return list
    }, [filtered])

    function openDropdown(start: number | null) {
        setTriggerStart(start)
        setQuery("")
        setActiveIndex(0)
        setOpen(true)
    }

    function closeDropdown() {
        setOpen(false)
        setTriggerStart(null)
    }

    function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
        const next = e.target.value
        onChange(next)

        const caret = e.target.selectionStart ?? next.length
        const upToCaret = next.slice(0, caret)
        const lastOpen = upToCaret.lastIndexOf("{{")
        const lastClose = upToCaret.lastIndexOf("}}")
        if (lastOpen !== -1 && lastOpen > lastClose) {
            setTriggerStart(lastOpen)
            setQuery(upToCaret.slice(lastOpen + 2).trim())
            setActiveIndex(0)
            setOpen(true)
        } else if (open && triggerStart !== null) {
            closeDropdown()
        }
    }

    function insertRef(ref: string) {
        const input = inputRef.current
        const caret = input?.selectionStart ?? value.length
        const before = triggerStart !== null ? value.slice(0, triggerStart) : value.slice(0, caret)
        const after = value.slice(caret)
        const inserted = `{{ ${ref} }}`
        const nextValue = before + inserted + after
        onChange(nextValue)
        closeDropdown()

        const nextCaret = before.length + inserted.length
        requestAnimationFrame(() => {
            input?.focus()
            input?.setSelectionRange(nextCaret, nextCaret)
        })
    }

    function handleKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
        if (!open) return
        if (e.key === "ArrowDown") {
            e.preventDefault()
            setActiveIndex(i => Math.min(i + 1, flatList.length - 1))
        } else if (e.key === "ArrowUp") {
            e.preventDefault()
            setActiveIndex(i => Math.max(i - 1, 0))
        } else if (e.key === "Enter") {
            e.preventDefault()
            const entry = flatList[activeIndex]
            if (entry) insertRef(entry.ref)
        } else if (e.key === "Escape") {
            e.stopPropagation()
            closeDropdown()
        }
    }

    const isInvalid = malformed || invalid.length > 0
    const tooltip = malformed
        ? "Malformed template expression"
        : invalid.length > 0
            ? `Unknown reference: ${invalid.join(", ")}`
            : undefined

    return (
        <div className={cls.Wrap}>
            <div className={cls.InputRow}>
                <input
                    ref={inputRef}
                    className={`${cls.Input} ${isInvalid ? cls.InputInvalid : ""}`}
                    value={value}
                    onChange={handleChange}
                    onKeyDown={handleKeyDown}
                    onBlur={() => setTimeout(closeDropdown, 150)}
                    placeholder={placeholder}
                    spellCheck={false}
                    data-tooltip-id={tooltip ? "root-tooltip" : undefined}
                    data-tooltip-content={tooltip}
                />
                <Button
                    type="button"
                    variant="ghost"
                    className={cls.AffixBtn}
                    onClick={() => openDropdown(null)}
                >
                    {"{}"}
                </Button>
            </div>
            {open && (
                <div className={cls.Dropdown}>
                    <SourceGroups
                        groups={filtered}
                        flatList={flatList}
                        activeIndex={activeIndex}
                        onSelect={insertRef}
                    />
                </div>
            )}
        </div>
    )
}

function SourceGroups({groups, flatList, activeIndex, onSelect}: {
    groups: { source: TemplateSource; entries: FlatEntry[]; rootMatches: boolean }[]
    flatList: { ref: string; label: string }[]
    activeIndex: number
    onSelect: (ref: string) => void
}) {
    if (groups.length === 0) {
        return <div className={cls.Empty}>No matching references.</div>
    }

    function indexOf(ref: string): number {
        return flatList.findIndex(f => f.ref === ref)
    }

    return (
        <>
            {groups.map(g => (
                <div className={cls.Group} key={g.source.id}>
                    {g.rootMatches && (
                        <div
                            className={`${cls.GroupLabel} ${indexOf(g.source.id) === activeIndex ? cls.EntryActive : ""}`}
                            onMouseDown={e => {
                                e.preventDefault()
                                onSelect(g.source.id)
                            }}
                        >
                            {g.source.label}
                        </div>
                    )}
                    {g.entries.map(e => (
                        <div
                            key={e.key}
                            className={`${cls.Entry} ${indexOf(e.ref) === activeIndex ? cls.EntryActive : ""}`}
                            style={{paddingLeft: `${0.5 + e.depth * 0.75}rem`}}
                            onMouseDown={ev => {
                                ev.preventDefault()
                                onSelect(e.ref)
                            }}
                            data-tooltip-id={e.description ? "root-tooltip" : undefined}
                            data-tooltip-content={e.description}
                        >
                            <span className={cls.EntryLabel}>{e.label}</span>
                            {e.type && <span className={cls.EntryType}>{e.type}</span>}
                        </div>
                    ))}
                </div>
            ))}
        </>
    )
}
