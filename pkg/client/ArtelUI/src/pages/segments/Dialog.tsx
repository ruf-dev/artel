import {useEffect} from "react"
import {createPortal} from "react-dom"
import {ModalClose} from "@vervstack/chures"

import cls from "@/pages/segments/Dialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"

// Module-level lazy-created portal target. Ensures the dialog is always mounted
// to a dedicated body child so it paints above any other portaled overlays
// (WorkbenchTweaksPanel, etc.) without z-index.
let dialogRootElement: HTMLDivElement | null = null

function getDialogRoot(): HTMLDivElement {
    if (!dialogRootElement && typeof document !== "undefined") {
        dialogRootElement = document.createElement("div")
        dialogRootElement.setAttribute("data-dialog-root", "")
        document.body.appendChild(dialogRootElement)
    }
    return dialogRootElement!
}

export default function Dialog() {
    const {children, closable, CloseDialog} = useDialog()

    useEffect(() => {
        function handleKeyDown(e: KeyboardEvent) {
            if (e.key === "Escape") CloseDialog()
        }
        document.addEventListener("keydown", handleKeyDown)
        return () => document.removeEventListener("keydown", handleKeyDown)
    }, [CloseDialog])

    // Re-append the portal target to body whenever the dialog opens, so it stays
    // the last body child and paints above overlays portaled earlier (the
    // WorkbenchTweaksPanel, floating menus). appendChild on an existing child
    // moves it to last.
    useEffect(() => {
        if (!children) return
        document.body.appendChild(getDialogRoot())
    }, [children])

    // Clean up portal target on unmount.
    useEffect(() => {
        return () => {
            if (dialogRootElement) {
                dialogRootElement.remove()
                dialogRootElement = null
            }
        }
    }, [])

    if (!children) return null

    const dialogContent = (
        <div
            className={cls.DialogContainer}
            onMouseDown={(e) => {
                if (e.target === e.currentTarget) CloseDialog()
            }}
        >
            <div
                className={cls.DialogWrapper}>
                <div
                    className={cls.DialogBackground}
                    onMouseDown={(e) => {
                        e.stopPropagation()
                    }}
                >
                    {
                        children.map((v, idx) => {
                            return (
                                <div key={idx} className={cls.DialogChildWrapper}>
                                    {v}
                                </div>)
                        })
                    }
                </div>
                {closable && (
                    <ModalClose className={cls.DialogCloseButton} onClick={CloseDialog}/>
                )}
            </div>
        </div>
    )

    return createPortal(dialogContent, getDialogRoot())
}
