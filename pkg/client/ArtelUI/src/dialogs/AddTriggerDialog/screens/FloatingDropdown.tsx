import {useEffect, useState} from "react"
import type {RefObject} from "react"
import {createPortal} from "react-dom"

import {Dropdown} from "@vervstack/chures"
import type {DropdownOption} from "@vervstack/chures"
import cls from "@/dialogs/AddTriggerDialog/AddTriggerDialog.module.css"

// Portals a chures <Dropdown> to document.body, positioned under its anchor via
// getBoundingClientRect — same "never use z-index" pattern as TemplateInput.tsx.
export default function FloatingDropdown({anchorRef, options, onPick, onClose, placeholder}: {
    anchorRef: RefObject<HTMLElement | null>
    options: DropdownOption[]
    onPick: (opt: DropdownOption) => void
    onClose: () => void
    placeholder?: string
}) {
    const [rect, setRect] = useState<{ top: number; left: number; width: number } | null>(null)

    useEffect(() => {
        function reposition() {
            const el = anchorRef.current
            if (!el) return
            const box = el.getBoundingClientRect()
            setRect({top: box.bottom + 4, left: box.left, width: box.width})
        }

        reposition()
        window.addEventListener("scroll", reposition, true)
        window.addEventListener("resize", reposition)
        return () => {
            window.removeEventListener("scroll", reposition, true)
            window.removeEventListener("resize", reposition)
        }
    }, [anchorRef])

    if (!rect) return null

    return createPortal(
        <div className={cls.DropdownPanel} style={{top: rect.top, left: rect.left, width: rect.width}}>
            <Dropdown options={options} onPick={onPick} onClose={onClose} anchorRef={anchorRef} placeholder={placeholder}/>
        </div>,
        document.body,
    )
}
