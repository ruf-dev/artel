import { useState } from "react"
import { Button, Input } from "@vervstack/chures"

import { useDialog } from "@/app/hooks/Dialog.ts"
import DropZone from "@/components/atoms/DropZone/DropZone.tsx"
import cls from "@/pages/notes/components/ImportZipDialog/ImportZipDialog.module.css"

interface Props {
    onConfirm: (destPath: string, zipData: Uint8Array) => Promise<void>
}

export default function ImportZipDialog({ onConfirm }: Props) {
    const { CloseDialog } = useDialog()
    const [file, setFile] = useState<File | null>(null)
    const [destPath, setDestPath] = useState("")
    const [destPathTouched, setDestPathTouched] = useState(false)
    const [loading, setLoading] = useState(false)

    function handleFileSelected(selected: File) {
        setFile(selected)
        if (!destPathTouched) {
            setDestPath(selected.name.replace(/\.zip$/i, ""))
        }
    }

    function handleDestPathChange(value: string) {
        setDestPathTouched(true)
        setDestPath(value)
    }

    function handleImport() {
        if (!file) return
        setLoading(true)
        file.arrayBuffer()
            .then(buf => onConfirm(destPath, new Uint8Array(buf)))
            .finally(() => setLoading(false))
    }

    return (
        <div className={cls.ImportZipContainer} role="dialog" aria-modal="true">
            <h2 className={cls.ImportZipTitle}>Import from .zip</h2>
            <DropZone accept=".zip" selectedFileName={file?.name} onFileSelected={handleFileSelected} />
            <p className={cls.ImportZipLabel}>Destination path (empty = vault root)</p>
            <Input
                inputClassName={cls.ImportZipInput}
                type="text"
                value={destPath}
                setValue={handleDestPathChange}
                disabled={loading}
                placeholder="e.g. Imported/Work"
                spellCheck={false}
            />
            <div className={cls.ImportZipActions}>
                <Button variant="ghost" onClick={CloseDialog} disabled={loading}>
                    Cancel
                </Button>
                <Button variant="primary" onClick={handleImport} disabled={loading || !file}>
                    {loading ? "…" : "Import"}
                </Button>
            </div>
        </div>
    )
}
