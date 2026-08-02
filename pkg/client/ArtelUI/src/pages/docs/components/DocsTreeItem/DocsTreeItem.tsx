import { cn } from "@/app/utils/cn.ts"
import ArrowIcon from "@/pages/docs/components/icons/ArrowIcon.tsx"
import FolderIcon from "@/pages/docs/components/icons/FolderIcon.tsx"
import FileIcon from "@/pages/docs/components/icons/FileIcon.tsx"
import cls from "@/pages/docs/components/DocsTreeItem/DocsTreeItem.module.css"

// Stripped-down read-only row for the /docs tree: no drag-and-drop, no kebab menu, no
// rename/download/delete actions — those belong to the authenticated /notes editor, not the
// public reader. Icons are small dependency-free SVGs duplicated locally (not cross-imported
// from pages/notes/components/, which is colocation-local to NotesPage per CLAUDE.md).

interface DocsTreeItemProps {
    name: string
    path?: string
    isFolder?: boolean
    isOpen?: boolean
    active?: boolean
    depth?: number
    onClick?: () => void
}

export default function DocsTreeItem(props: DocsTreeItemProps) {
    const { name, path, isFolder, isOpen, active } = props
    const { depth = 0, onClick } = props
    const paddingLeft = 1.12 + depth * 0.84

    return (
        <div
            className={cn(cls.DocsTreeItemContainer, active && cls.DocsTreeItemActive)}
            style={{ padding: `0.28rem 1.12rem 0.28rem ${paddingLeft}rem` }}
            onClick={onClick}
            data-path={path}
        >
            {isFolder ? <ArrowIcon open={!!isOpen} /> : <span className={cls.ArrowSpacer} />}
            {isFolder ? <FolderIcon /> : <FileIcon />}
            <span className={cls.Label}>{name}</span>
        </div>
    )
}
