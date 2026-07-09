import { useEffect, useRef, useState } from "react"

import cls from "@/pages/notes/components/SaveStatusIndicator/SaveStatusIndicator.module.css"
import SpinnerIcon from "@/pages/notes/components/icons/SpinnerIcon.tsx"
import CheckIcon from "@/pages/notes/components/icons/CheckIcon.tsx"
import ErrorDotIcon from "@/pages/notes/components/icons/ErrorDotIcon.tsx"

export type SaveStatus = 'idle' | 'dirty' | 'saving' | 'saved' | 'error'

interface SaveStatusIndicatorProps {
    status: SaveStatus
    errorMessage?: string
}

export default function SaveStatusIndicator({ status, errorMessage }: SaveStatusIndicatorProps) {
    const [fading, setFading] = useState(false)
    const [hidden, setHidden] = useState(false)
    const fadeTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
    const hideTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

    useEffect(() => {
        if (fadeTimerRef.current) clearTimeout(fadeTimerRef.current)
        if (hideTimerRef.current) clearTimeout(hideTimerRef.current)
        setFading(false)
        setHidden(false)

        if (status === 'saved') {
            fadeTimerRef.current = setTimeout(() => {
                setFading(true)
                hideTimerRef.current = setTimeout(() => setHidden(true), 600)
            }, 3000)
        }

        return () => {
            if (fadeTimerRef.current) clearTimeout(fadeTimerRef.current)
            if (hideTimerRef.current) clearTimeout(hideTimerRef.current)
        }
    }, [status])

    if (status === 'idle' || hidden) return null

    const style = fading ? { opacity: 0, transition: 'opacity 600ms ease' } : undefined

    return (
        <div className={cls.SaveStatusContainer} style={style}>
            {status === 'dirty' && (
                <span className={cls.TextDirty}>unsaved</span>
            )}
            {status === 'saving' && (
                <>
                    <SpinnerIcon />
                    <span className={cls.TextSaving}>saving…</span>
                </>
            )}
            {status === 'saved' && (
                <>
                    <CheckIcon />
                    <span className={cls.TextSaved}>saved</span>
                </>
            )}
            {status === 'error' && (
                <>
                    <ErrorDotIcon />
                    <span className={cls.TextError}>{errorMessage ?? 'save failed'}</span>
                </>
            )}
        </div>
    )
}
