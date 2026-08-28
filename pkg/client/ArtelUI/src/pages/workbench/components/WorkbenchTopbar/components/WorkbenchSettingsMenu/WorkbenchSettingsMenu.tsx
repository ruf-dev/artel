import {useEffect, useRef, useState} from "react"
import {createPortal} from "react-dom"
import {Button, ConfirmDialog} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import {useWorkbenchMutations} from "@/app/hooks/Workbench.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import WorkbenchSettingsMenuItem from "@/pages/workbench/components/WorkbenchTopbar/components/WorkbenchSettingsMenu/components/WorkbenchSettingsMenuItem.tsx"
// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import cls from "@/pages/workbench/components/WorkbenchTopbar/components/WorkbenchSettingsMenu/WorkbenchSettingsMenu.module.css"

interface MenuRect {
    top?: number
    bottom?: number
    right: number
}

interface WorkbenchSettingsMenuProps {
    vaultId: string
}

export default function WorkbenchSettingsMenu({vaultId}: WorkbenchSettingsMenuProps) {
    const [menuOpen, setMenuOpen] = useState(false)
    const [menuRect, setMenuRect] = useState<MenuRect | null>(null)
    const wrapRef = useRef<HTMLDivElement>(null)
    const {CloseDialog, OpenDialog} = useDialog()
    const {remove} = useWorkbenchMutations(vaultId)
    const bakeError = useBakeError()

    // Portal to document.body and position via getBoundingClientRect instead of z-index,
    // per the project's no-z-index rule — same pattern as TopbarUserMenu.tsx.
    useEffect(() => {
        if (!menuOpen) return

        function reposition() {
            const el = wrapRef.current
            if (!el) return
            const box = el.getBoundingClientRect()
            setMenuRect({top: box.bottom + 8, right: window.innerWidth - box.right})
        }

        reposition()
        window.addEventListener("scroll", reposition, true)
        window.addEventListener("resize", reposition)
        return () => {
            window.removeEventListener("scroll", reposition, true)
            window.removeEventListener("resize", reposition)
        }
    }, [menuOpen])

    function handleDelete() {
        setMenuOpen(false)
        OpenDialog(
            <ConfirmDialog
                title="Delete Workbench"
                message="Delete this Workbench? This destroys its container and workspace — all
                    files and state saved in it are lost. This cannot be undone."
                confirmLabel="Delete"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() => remove()
                    .then(() => undefined)
                    .catch(e => bakeError("Failed to delete workbench", e))}
            />
        )
    }

    return (
        <div className={cls.WorkbenchSettingsMenuContainer} ref={wrapRef}>
            <Button
                variant="secondary"
                className={cls.TriggerButton}
                onClick={e => {
                    e.stopPropagation()
                    setMenuOpen(v => !v)
                }}
            >
                <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
                     strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                    <circle cx="12" cy="12" r="3"/>
                    <path d="M19.07 4.93a10 10 0 0 1 0 14.14M4.93 4.93a10 10 0 0 0 0 14.14"/>
                </svg>
            </Button>
            {menuOpen && menuRect && createPortal(
                <>
                    <div className={cls.Backdrop} onClick={() => setMenuOpen(false)}/>
                    <div
                        className={cls.Menu}
                        role="menu"
                        style={menuRect}
                    >
                        <WorkbenchSettingsMenuItem
                            label="Delete"
                            onClick={handleDelete}
                            danger
                            icon={(
                                <svg
                                    viewBox="0 0 24 24"
                                    width="15"
                                    height="15"
                                    fill="none"
                                    stroke="currentColor"
                                    strokeWidth="1.6"
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                >
                                    <polyline points="3 6 5 6 21 6"/>
                                    {/* eslint-disable-next-line max-len -- SVG path data can't be wrapped */}
                                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                                    <line x1="10" y1="11" x2="10" y2="17"/>
                                    <line x1="14" y1="11" x2="14" y2="17"/>
                                </svg>
                            )}
                        />
                    </div>
                </>,
                document.body,
            )}
        </div>
    )
}
