// TODO: chures has no multiline variant yet, drop this wrapper once it does
import {useLayoutEffect, useRef} from "react"

import {cn} from "@/app/utils/cn.ts"
import cls from "@/components/atoms/Textarea/Textarea.module.css"

interface Props {
    value: string
    setValue: (v: string) => void
    disabled?: boolean
    placeholder?: string
    className?: string
    rows?: number
    onKeyDown?: (e: React.KeyboardEvent<HTMLTextAreaElement>) => void
    autoGrow?: boolean
    maxHeightRem?: number
    minHeightRem?: number
}

const DEFAULT_MAX_HEIGHT_REM = 11

export default function Textarea(props: Props) {
    const ref = useRef<HTMLTextAreaElement>(null)
    const maxHeightRem = props.maxHeightRem ?? DEFAULT_MAX_HEIGHT_REM

    useLayoutEffect(() => {
        if (!props.autoGrow) return
        const el = ref.current
        if (!el) return
        const rootFontSize = parseFloat(getComputedStyle(document.documentElement).fontSize) || 16

        // Empty state gets a fixed height instead of the usual scrollHeight measurement —
        // font line-box metrics for the composer's font are inconsistent enough that
        // measuring the "empty" height isn't reliably the same across engines/loads, which
        // broke vertical-centering against the absolutely-positioned send button. A caller
        // that cares (ChatComposer) sizes minHeightRem to match that button's own height.
        if (!props.value.trim() && props.minHeightRem !== undefined) {
            el.style.height = `${props.minHeightRem * rootFontSize}px`
            return
        }
        el.style.height = "auto"
        el.style.height = `${Math.min(el.scrollHeight, maxHeightRem * rootFontSize)}px`
    }, [props.value, props.autoGrow, maxHeightRem, props.minHeightRem])

    return (
        <div className={cn(cls.TextareaContainer, props.autoGrow && cls.AutoGrow, props.className)}>
            <textarea
                ref={ref}
                value={props.value}
                onChange={e => props.setValue(e.target.value)}
                onKeyDown={props.onKeyDown}
                disabled={props.disabled}
                placeholder={props.placeholder}
                rows={props.rows ?? (props.autoGrow ? 1 : 4)}
                style={props.autoGrow ? {maxHeight: `${maxHeightRem}rem`} : undefined}
            />
        </div>
    )
}
