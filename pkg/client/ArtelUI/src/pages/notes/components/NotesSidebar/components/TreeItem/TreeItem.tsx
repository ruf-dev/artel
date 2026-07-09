import {Button} from "@vervstack/chures"

import ArrowIcon from "@/pages/notes/components/NotesSidebar/components/icons/ArrowIcon.tsx"
import FileIcon from "@/pages/notes/components/NotesSidebar/components/icons/FileIcon.tsx"
import FolderIcon from "@/pages/notes/components/NotesSidebar/components/icons/FolderIcon.tsx"
import cls from "@/pages/notes/components/NotesSidebar/components/TreeItem/TreeItem.module.css"

interface TreeItemProps {
    name: string
    subtitle?: string
    active?: boolean
    depth?: number
    isFolder?: boolean
    isOpen?: boolean
    onClick?: () => void
    onAddInFolder?: () => void
}

export default function TreeItem(props: TreeItemProps) {
    const { name, subtitle, active } = props
    const { isFolder, isOpen, onClick, onAddInFolder } = props
    const depth = props.depth ?? 0
    const paddingLeft = 1.12 + depth * 0.84
    const rowClass = `${cls.TreeItemRowContainer}${active ? ` ${cls.TreeItemRowActive}` : ""}`

    return (
        <div className={rowClass} style={{ padding: `0.28rem 1.12rem 0.28rem ${paddingLeft}rem` }} onClick={onClick}>
            {isFolder ? <ArrowIcon open={!!isOpen} /> : <span className={cls.ArrowSpacer} />}
            {isFolder ? <FolderIcon /> : <FileIcon />}
            <div className={cls.TreeItemMain}>
                <span className={cls.TreeItemLabel}>{name}</span>
                {subtitle && <span className={cls.TreeItemSubtitle}>{subtitle}</span>}
            </div>
            {isFolder && (
                <Button
                    variant="ghost"
                    className={cls.FolderAddBtn}
                    onClick={e => { e.stopPropagation(); onAddInFolder?.() }}
                    title="New note here"
                >
                    <svg viewBox="0 0 12 12" width={11} height={11} fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round">
                        <path d="M6 1v10M1 6h10" />
                    </svg>
                </Button>
            )}
        </div>
    )
}
