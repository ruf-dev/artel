import { useEffect, useRef, useState } from "react"

import cls from "./SaveStatusIndicator.module.css"

export type SaveStatus = 'idle' | 'dirty' | 'saving' | 'saved' | 'error'

interface SaveStatusIndicatorProps {
    status: SaveStatus
    errorMessage?: string
}

function SpinnerIcon() {
    return (
        <svg viewBox="0 0 12 12" width="10" height="10" className={cls.Spinner}>
            <circle cx="6" cy="6" r="4.5" fill="none" stroke="rgba(255,255,255,0.18)" strokeWidth="1.5" />
            <path d="M6 1.5 A4.5 4.5 0 0 1 10.5 6" fill="none" stroke="rgba(255,255,255,0.55)" strokeWidth="1.5" strokeLinecap="round" />
        </svg>
    )
}

function CheckIcon() {
    return (
        <svg viewBox="0 0 12 12" width="10" height="10">
            <path d="M2 6 L5 9 L10 3" fill="none" stroke="rgba(255,255,255,0.45)" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
    )
}

function ErrorDotIcon() {
    return (
        <svg viewBox="0 0 8 8" width="6" height="6">
            <circle cx="4" cy="4" r="3" fill="#c1253d" />
        </svg>
    )
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
