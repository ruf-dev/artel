import {useEffect, useState, RefObject} from "react"

interface DropdownRect {
    top: number
    left: number
    width: number
}

export function useDropdownPosition(open: boolean, wrapRef: RefObject<HTMLDivElement | null>): DropdownRect | null {
    const [dropdownRect, setDropdownRect] = useState<DropdownRect | null>(null)

    useEffect(() => {
        if (!open) return

        function reposition() {
            const el = wrapRef.current
            if (!el) return
            const box = el.getBoundingClientRect()
            setDropdownRect({top: box.bottom + 4, left: box.left, width: box.width})
        }

        reposition()
        window.addEventListener("scroll", reposition, true)
        window.addEventListener("resize", reposition)
        return () => {
            window.removeEventListener("scroll", reposition, true)
            window.removeEventListener("resize", reposition)
        }
    }, [open, wrapRef])

    return dropdownRect
}
