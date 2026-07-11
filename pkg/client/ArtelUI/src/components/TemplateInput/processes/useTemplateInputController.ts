import {useMemo, useRef, useState, RefObject, KeyboardEvent} from "react"

import {extractRefs, validateRef} from "@/processes/tractTemplate.ts"
import {
    buildFlatList,
    filterSources,
    FilteredGroup,
    FlatListItem,
    TemplateSource,
} from "@/components/TemplateInput/processes/templateRefs.ts"
import {useDropdownPosition} from "@/components/TemplateInput/processes/useDropdownPosition.ts"

interface Result {
    wrapRef: RefObject<HTMLDivElement | null>
    inputRef: RefObject<HTMLInputElement | null>
    open: boolean
    activeIndex: number
    filtered: FilteredGroup[]
    flatList: FlatListItem[]
    isInvalid: boolean
    tooltip: string | undefined
    dropdownRect: { top: number; left: number; width: number } | null
    openDropdown: (start: number | null) => void
    closeDropdown: () => void
    handleChange: (next: string) => void
    handleKeyDown: (e: KeyboardEvent<HTMLInputElement>) => void
    insertRef: (ref: string) => void
}

export function useTemplateInputController(
    sources: TemplateSource[], value: string, onChange: (value: string) => void): Result {
    const wrapRef = useRef<HTMLDivElement>(null)
    const inputRef = useRef<HTMLInputElement>(null)
    const [open, setOpen] = useState(false)
    const [query, setQuery] = useState("")
    const [triggerStart, setTriggerStart] = useState<number | null>(null)
    const [activeIndex, setActiveIndex] = useState(0)

    // Portal the dropdown to document.body and position it with getBoundingClientRect instead
    // of z-index — see the "Never use z-index" rule in CLAUDE.md. Being appended after #root in
    // the DOM, it paints above the app with no stacking-context gymnastics needed.
    const dropdownRect = useDropdownPosition(open, wrapRef)

    const sourceIds = useMemo(() => sources.map(s => s.id), [sources])

    const {invalid, malformed} = useMemo(() => {
        try {
            const refs = extractRefs(value)
            return {invalid: refs.filter(r => !validateRef(r, sourceIds)), malformed: false}
        } catch {
            return {invalid: [], malformed: true}
        }
    }, [value, sourceIds])

    const filtered = useMemo(() => filterSources(sources, query), [sources, query])
    const flatList = useMemo(() => buildFlatList(filtered), [filtered])

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

    function handleChange(next: string) {
        onChange(next)

        const caret = inputRef.current?.selectionStart ?? next.length
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

    function handleKeyDown(e: KeyboardEvent<HTMLInputElement>) {
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

    return {
        wrapRef,
        inputRef,
        open,
        activeIndex,
        filtered,
        flatList,
        isInvalid,
        tooltip,
        dropdownRect,
        openDropdown,
        closeDropdown,
        handleChange,
        handleKeyDown,
        insertRef,
    }
}
