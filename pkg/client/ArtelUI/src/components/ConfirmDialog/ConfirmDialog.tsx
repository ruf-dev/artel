import {useState} from "react"
import cls from "@/components/ConfirmDialog/ConfirmDialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"

interface Props {
    title: string
    message: string
    confirmLabel?: string
    cancelLabel?: string
    danger?: boolean
    onConfirm: () => void | Promise<void>
}

export default function ConfirmDialog({
    title,
    message,
    confirmLabel = "Confirm",
    cancelLabel = "Cancel",
    danger = false,
    onConfirm,
}: Props) {
    const {CloseDialog} = useDialog()
    const [loading, setLoading] = useState(false)

    async function handleConfirm() {
        setLoading(true)
        try {
            await onConfirm()
        } finally {
            setLoading(false)
            CloseDialog()
        }
    }

    return (
        <div className={cls.ConfirmContainer} role="dialog" aria-modal="true">
            <h2 className={cls.ConfirmTitle}>{title}</h2>
            <p className={cls.ConfirmMessage}>{message}</p>
            <div className={cls.ConfirmActions}>
                <button
                    className={cls.BtnCancel}
                    type="button"
                    onClick={CloseDialog}
                    disabled={loading}
                >
                    {cancelLabel}
                </button>
                <button
                    className={danger ? cls.BtnDanger : cls.BtnConfirm}
                    type="button"
                    onClick={handleConfirm}
                    disabled={loading}
                >
                    {loading ? "…" : confirmLabel}
                </button>
            </div>
        </div>
    )
}
