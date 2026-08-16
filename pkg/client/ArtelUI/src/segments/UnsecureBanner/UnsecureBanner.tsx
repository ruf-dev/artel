import {useLayoutEffect, useRef, useState} from "react"
import {createPortal} from "react-dom"
import {Button, CheckmarkIcon} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import {useAppConfig} from "@/app/hooks/AppConfig.ts"
import cls from "@/segments/UnsecureBanner/UnsecureBanner.module.css"

const CREDS_KEY_COMMAND = "openssl rand -hex 32"
const HOVER_CLOSE_DELAY_MS = 250
const PANEL_VIEWPORT_MARGIN_PX = 8

// .CredsKeyPanel is centered on this value via `transform: translateX(-50%)` —
// clamp it (using the panel's *actual* measured width, not an assumed one, since
// rem-per-px varies with browser/OS font-size settings) so it stays within the
// viewport when the anchor sits near the left/right edge, instead of spilling
// off-screen.
function clampPanelCenter(anchorRect: DOMRect, panelWidth: number): number {
    const center = anchorRect.left + anchorRect.width / 2
    const half = panelWidth / 2
    const min = PANEL_VIEWPORT_MARGIN_PX + half
    const max = window.innerWidth - PANEL_VIEWPORT_MARGIN_PX - half
    return Math.min(Math.max(center, min), max)
}

export default function UnsecureBanner() {
    const noAuthEnabled = useAppConfig((state) => state.noAuthEnabled)
    const credsEncrypted = useAppConfig((state) => state.credsEncrypted)
    const [copied, setCopied] = useState(false)
    const [open, setOpen] = useState(false)
    const [anchorRect, setAnchorRect] = useState<DOMRect | null>(null)
    const [panelCenter, setPanelCenter] = useState<number | null>(null)
    const anchorRef = useRef<HTMLSpanElement>(null)
    const panelRef = useRef<HTMLDivElement>(null)
    const closeTimer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)

    // Reposition against the panel's real rendered width once it's mounted —
    // runs before paint, so there's no visible jump.
    useLayoutEffect(() => {
        if (!open || !anchorRect || !panelRef.current) return
        setPanelCenter(clampPanelCenter(anchorRect, panelRef.current.offsetWidth))
    }, [open, anchorRect])

    if (!noAuthEnabled && credsEncrypted) return null

    // Hover close is delayed (and cancelable from either the trigger or the
    // portaled panel below) instead of closing on the anchor's mouseleave —
    // otherwise crossing the gap between them re-triggers a close mid-transit.
    function handleEnter() {
        clearTimeout(closeTimer.current)
        if (anchorRef.current) setAnchorRect(anchorRef.current.getBoundingClientRect())
        setOpen(true)
    }

    function handleLeave() {
        closeTimer.current = setTimeout(() => setOpen(false), HOVER_CLOSE_DELAY_MS)
    }

    function handleCopy() {
        navigator.clipboard.writeText(CREDS_KEY_COMMAND).then(() => {
            setCopied(true)
            setTimeout(() => setCopied(false), 2000)
        })
    }

    return (
        <>
            {noAuthEnabled && (
                <div className={cls.UnsecureBannerContainer}>
                    running unsecure instance — authentication is disabled
                </div>
            )}
            {!credsEncrypted && (
                <div className={cls.UnsecureBannerContainer}>
                    running in insecure mode — all credentials are stored as plain text. generate an
                    encryption key and set it as ENVIRONMENT_CREDS_ENCRYPTION_KEY
                    <span
                        ref={anchorRef}
                        className={cls.CredsKeyHelpWrap}
                        onMouseEnter={handleEnter}
                        onMouseLeave={handleLeave}
                    >
                        <Button
                            variant="ghost"
                            className={cls.CredsKeyHelp}
                            aria-label="How do I generate an encryption key?"
                        >?</Button>
                    </span>
                    {open && anchorRect && createPortal(
                        <div
                            ref={panelRef}
                            className={cls.CredsKeyPanel}
                            style={{
                                top: anchorRect.bottom + 8,
                                left: panelCenter ?? anchorRect.left + anchorRect.width / 2,
                            }}
                            onMouseEnter={handleEnter}
                            onMouseLeave={handleLeave}
                        >
                            <p className={cls.TooltipText}>
                                Generate a key, then set it as <code>ENVIRONMENT_CREDS_ENCRYPTION_KEY</code>
                            </p>
                            <div className={cls.TerminalBox}>
                                <Button
                                    variant="ghost"
                                    className={cn(cls.CopyIconBtn, copied && cls.CopyIconBtnCopied)}
                                    onClick={handleCopy}
                                    aria-label="Copy command"
                                    title={copied ? "Copied!" : "Copy"}
                                >
                                    {copied ? <CheckmarkIcon size={12}/> : (
                                        <svg
                                            viewBox="0 0 24 24" width="12" height="12" fill="none"
                                            stroke="currentColor" strokeWidth="1.8" strokeLinecap="round"
                                            strokeLinejoin="round"
                                        >
                                            <rect x="9" y="9" width="13" height="13" rx="2"/>
                                            <path d="M5 15V5a2 2 0 0 1 2-2h10"/>
                                        </svg>
                                    )}
                                </Button>
                                <code className={cls.Command}>{CREDS_KEY_COMMAND}</code>
                            </div>
                        </div>,
                        document.body,
                    )}
                </div>
            )}
        </>
    )
}
