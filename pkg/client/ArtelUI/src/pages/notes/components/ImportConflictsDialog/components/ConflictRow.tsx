import { Button, Input } from "@vervstack/chures"

import { cn } from "@/app/utils/cn.ts"
import { ImportConflictAction, ImportResolutionInput } from "@/app/hooks/Notes.ts"
import cls from "@/pages/notes/components/ImportConflictsDialog/components/ConflictRow.module.css"

interface Props {
    path: string
    resolution: ImportResolutionInput
    onChange: (resolution: ImportResolutionInput) => void
}

export default function ConflictRow({ path, resolution, onChange }: Props) {
    return (
        <div className={cls.ConflictRowContainer}>
            <span className={cls.ConflictPath}>{path}</span>
            <Button
                variant={resolution.action === ImportConflictAction.SKIP ? "secondary" : "ghost"}
                className={cls.ActionBtn}
                onClick={() => onChange({ path, action: ImportConflictAction.SKIP })}
            >
                Skip
            </Button>
            <Button
                variant={resolution.action === ImportConflictAction.OVERWRITE ? "secondary" : "ghost"}
                className={cls.ActionBtn}
                onClick={() => onChange({ path, action: ImportConflictAction.OVERWRITE })}
            >
                Overwrite
            </Button>
            <Button
                variant={resolution.action === ImportConflictAction.RENAME ? "secondary" : "ghost"}
                className={cls.ActionBtn}
                onClick={() => {
                    const renameTo = resolution.renameTo ?? path
                    onChange({ path, action: ImportConflictAction.RENAME, renameTo })
                }}
            >
                Rename
            </Button>
            {resolution.action === ImportConflictAction.RENAME && (
                <Input
                    className={cn(cls.RenameInputWrapper)}
                    inputClassName={cls.RenameInput}
                    type="text"
                    value={resolution.renameTo ?? ""}
                    setValue={v => onChange({ ...resolution, renameTo: v })}
                    spellCheck={false}
                />
            )}
        </div>
    )
}
