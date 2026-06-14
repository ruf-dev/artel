import { useState } from "react"
import cls from "./CreateNoteDialog.module.css"
import { useDialog, useDialogKeyboard } from "@/app/hooks/Dialog"

function ensureExtension(path: string): string {
    const filename = path.split("/").pop() ?? path
    return filename.includes(".") ? path : path + ".md"
}

interface Props {
    onConfirm: (path: string) => Promise<void>
}

export default function CreateNoteDialog({ onConfirm }: Props) {
    const { CloseDialog } = useDialog()
    const [path, setPath] = useState("")
    const [loading, setLoading] = useState(false)

    async function handleConfirm() {
        if (!path.trim()) return
        setLoading(true)
        try {
            await onConfirm(ensureExtension(path.trim()))
        } finally {
            setLoading(false)
            CloseDialog()
        }
    }

    const handleKeyDown = useDialogKeyboard(handleConfirm)

    return (
        <div className={cls.CreateNoteContainer} role="dialog" aria-modal="true">
            <h2 className={cls.CreateNoteTitle}>New Note</h2>
            <p className={cls.CreateNoteLabel}>Path</p>
            <input
                className={cls.CreateNoteInput}
                type="text"
                placeholder="folder/note-name"
                value={path}
                onChange={e => setPath(e.target.value)}
                onKeyDown={handleKeyDown}
                disabled={loading}
                autoFocus
                spellCheck={false}
            />
            <div className={cls.CreateNoteActions}>
                <button
                    className={cls.BtnCancel}
                    type="button"
                    onClick={CloseDialog}
                    disabled={loading}
                >
                    Cancel
                </button>
                <button
                    className={cls.BtnConfirm}
                    type="button"
                    onClick={handleConfirm}
                    disabled={loading || !path.trim()}
                >
                    {loading ? "…" : "Create"}
                </button>
            </div>
        </div>
    )
}
