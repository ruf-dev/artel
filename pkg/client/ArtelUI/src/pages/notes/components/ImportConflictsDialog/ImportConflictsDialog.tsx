import { useState } from "react"
import { Button } from "@vervstack/chures"

import { useDialog } from "@/app/hooks/Dialog.ts"
import { ImportConflictAction, ImportResolutionInput } from "@/app/hooks/Notes.ts"
import ConflictRow from "@/pages/notes/components/ImportConflictsDialog/components/ConflictRow.tsx"
import cls from "@/pages/notes/components/ImportConflictsDialog/ImportConflictsDialog.module.css"

interface Props {
    conflicts: string[]
    onConfirm: (resolutions: ImportResolutionInput[]) => Promise<void>
}

function initialResolutions(conflicts: string[]): Record<string, ImportResolutionInput> {
    const entries = conflicts.map(path => [path, { path, action: ImportConflictAction.SKIP }] as const)
    return Object.fromEntries(entries)
}

export default function ImportConflictsDialog({ conflicts, onConfirm }: Props) {
    const { CloseDialog } = useDialog()
    const [resolutions, setResolutions] = useState(() => initialResolutions(conflicts))
    const [loading, setLoading] = useState(false)

    function handleChange(path: string, resolution: ImportResolutionInput) {
        setResolutions(prev => ({ ...prev, [path]: resolution }))
    }

    function handleConfirm() {
        setLoading(true)
        onConfirm(Object.values(resolutions)).finally(() => setLoading(false))
    }

    return (
        <div className={cls.ImportConflictsContainer} role="dialog" aria-modal="true">
            <h2 className={cls.ImportConflictsTitle}>Resolve conflicts</h2>
            <p className={cls.ImportConflictsWarning}>
                {conflicts.length} item(s) already exist at the destination. Renaming does not update any
                links or the note index pointing at the old name — you may need to fix those manually.
            </p>
            <div className={cls.ConflictList}>
                {conflicts.map(path => (
                    <ConflictRow
                        key={path}
                        path={path}
                        resolution={resolutions[path]}
                        onChange={r => handleChange(path, r)}
                    />
                ))}
            </div>
            <div className={cls.ImportConflictsActions}>
                <Button variant="ghost" onClick={CloseDialog} disabled={loading}>
                    Cancel
                </Button>
                <Button variant="primary" onClick={handleConfirm} disabled={loading}>
                    {loading ? "…" : "Resolve & Import"}
                </Button>
            </div>
        </div>
    )
}
