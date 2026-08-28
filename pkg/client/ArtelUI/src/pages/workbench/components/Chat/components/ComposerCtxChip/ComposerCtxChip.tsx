import {Button} from "@vervstack/chures"

import {getItemName} from "@/components/FileTree/fileTree.ts"
import FileIcon from "@/components/FileTree/icons/FileIcon.tsx"
import CloseIcon from "@/icons/common/CloseIcon.tsx"
import cls from "@/pages/workbench/components/Chat/components/ComposerCtxChip/ComposerCtxChip.module.css"

interface Props {
    path: string
    onRemove: (path: string) => void
}

// One removable "attached vault file" chip shown in the composer's context row.
// Shows the file basename with the full vault-relative path as the title; the ×
// detaches it. Matches ComposerChipRow's chip visual language (pill border, muted
// fg) — but those chips open Tweaks, this one is unrelated.
export default function ComposerCtxChip({path, onRemove}: Props) {
    return (
        <span className={cls.ComposerCtxChipContainer} title={path}>
            <FileIcon/>
            <span className={cls.Label}>{getItemName({path})}</span>
            <Button
                variant="unstyled"
                className={cls.Remove}
                onClick={() => onRemove(path)}
                aria-label={`Remove ${path}`}
            >
                <CloseIcon className={cls.RemoveIcon}/>
            </Button>
        </span>
    )
}
