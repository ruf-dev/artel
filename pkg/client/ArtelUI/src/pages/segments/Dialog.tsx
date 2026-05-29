import {useEffect} from "react"
import cls from "@/pages/segments/Dialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import ModalClose from "@/components/ModalClose/ModalClose"

export default function Dialog() {
    const {children, closable, CloseDialog} = useDialog()

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === "Escape") CloseDialog()
        }
        document.addEventListener("keydown", handleKeyDown)
        return () => document.removeEventListener("keydown", handleKeyDown)
    }, [CloseDialog])

    if (!children) return null

    return (
        <div
            className={cls.DialogContainer}
            onMouseDown={(e) => {
                if (e.target === e.currentTarget) CloseDialog()
            }}
        >
            <div className={cls.DialogWrapper}>
                {closable && (
                    <ModalClose className={cls.DialogCloseButton} onClick={CloseDialog}/>
                )}
                <div
                    onMouseDown={(e) => {
                        e.stopPropagation()
                    }}
                >
                    {
                        children.map((v, idx) => {
                            return (<div key={idx}>
                                {v}
                            </div>)
                        })
                    }
                </div>
            </div>
        </div>
    )
}
