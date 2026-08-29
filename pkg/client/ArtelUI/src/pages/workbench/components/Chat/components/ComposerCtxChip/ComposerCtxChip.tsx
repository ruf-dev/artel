import {useState} from "react"
import {Button} from "@vervstack/chures"
import {MorphIcon} from "morphicons/react"
import {Trash2, X} from "lucide"

import {getItemName} from "@/components/FileTree/fileTree.ts"
import FileIcon from "@/components/FileTree/icons/FileIcon.tsx"
import cls from "@/pages/workbench/components/Chat/components/ComposerCtxChip/ComposerCtxChip.module.css"

interface Props {
    path: string
    onRemove: (path: string) => void
}

// One removable "attached vault file" chip shown in the composer's context row.
// Shows the file basename with the full vault-relative path as the title; the
// trash glyph morphs to an × on hover/focus as the detach affordance — mirrors
// NewChatButton's chures Button + MorphIcon shape. Matches ComposerChipRow's chip
// visual language (pill border, muted fg) — but those chips open Tweaks, this one
// is unrelated.
export default function ComposerCtxChip({path, onRemove}: Props) {
    const [hovered, setHovered] = useState(false)

    return (
        <span className={cls.ComposerCtxChipContainer} title={path}>
            <FileIcon/>
            <span className={cls.Label}>{getItemName({path})}</span>
            <Button
                variant="unstyled"
                className={cls.Remove}
                onClick={() => onRemove(path)}
                aria-label={`Remove ${path}`}
                onMouseEnter={() => setHovered(true)}
                onMouseLeave={() => setHovered(false)}
                onFocus={() => setHovered(true)}
                onBlur={() => setHovered(false)}
            >
                <MorphIcon
                    icon={hovered ? X : Trash2}
                    size={14}
                    strokeWidth={1.6}
                    className={cls.RemoveIcon}
                />
            </Button>
        </span>
    )
}
