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
}

const DEFAULT_MAX_HEIGHT_REM = 11

export default function Textarea(props: Props) {
    const ref = useRef<HTMLTextAreaElement>(null)
    const maxHeightRem = props.maxHeightRem ?? DEFAULT_MAX_HEIGHT_REM

    useLayoutEffect(() => {
        if (!props.autoGrow) return
        const el = ref.current
        if (!el) return
        el.style.height = "auto"
        const rootFontSize = parseFloat(getComputedStyle(document.documentElement).fontSize) || 16
        el.style.height = `${Math.min(el.scrollHeight, maxHeightRem * rootFontSize)}px`
    }, [props.value, props.autoGrow, maxHeightRem])

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
