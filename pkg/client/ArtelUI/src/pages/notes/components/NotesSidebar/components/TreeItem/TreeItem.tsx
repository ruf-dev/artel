import {DragEvent, useState} from "react"
import {Button} from "@vervstack/chures"

import { cn } from "@/app/utils/cn.ts"
import KebabMenu from "@/components/atoms/KebabMenu/KebabMenu.tsx"
import ArrowIcon from "@/pages/notes/components/NotesSidebar/components/icons/ArrowIcon.tsx"
import FileIcon from "@/pages/notes/components/NotesSidebar/components/icons/FileIcon.tsx"
import FolderIcon from "@/pages/notes/components/NotesSidebar/components/icons/FolderIcon.tsx"
import CopyIcon from "@/pages/notes/components/icons/CopyIcon.tsx"
import CheckIcon from "@/pages/notes/components/icons/CheckIcon.tsx"
import cls from "@/pages/notes/components/NotesSidebar/components/TreeItem/TreeItem.module.css"

interface TreeItemProps {
    name: string
    subtitle?: string
    path?: string
    active?: boolean
    highlighted?: boolean
    depth?: number
    isFolder?: boolean
    isOpen?: boolean
    onClick?: () => void
    onAddInFolder?: () => void
    onCopyPath?: () => void
    onDownloadFolder?: () => void
    onDeleteFolder?: () => void
    onDragStart?: (e: DragEvent<HTMLDivElement>) => void
    onDropTarget?: (e: DragEvent<HTMLDivElement>) => void
}

export default function TreeItem(props: TreeItemProps) {
    const { name, subtitle, path, active, highlighted } = props
    const { isFolder, isOpen, onClick, onAddInFolder, onCopyPath } = props
    const { onDownloadFolder, onDeleteFolder, onDragStart, onDropTarget } = props
    const [copied, setCopied] = useState(false)
    const [isDragOver, setIsDragOver] = useState(false)
    const depth = props.depth ?? 0
    const paddingLeft = 1.12 + depth * 0.84
    const acceptsDrop = !!isFolder && !!onDropTarget
    const rowClass = cn(
        cls.TreeItemRowContainer,
        active && cls.TreeItemRowActive,
        highlighted && cls.TreeItemRowHighlight,
        isDragOver && cls.TreeItemRowDropTarget
    )

    function handleDragOver(e: DragEvent<HTMLDivElement>) {
        if (!acceptsDrop) return
        e.preventDefault()
        setIsDragOver(true)
    }

    function handleDragLeave() {
        setIsDragOver(false)
    }

    function handleDrop(e: DragEvent<HTMLDivElement>) {
        if (!acceptsDrop) return
        e.preventDefault()
        setIsDragOver(false)
        onDropTarget?.(e)
    }

    return (
        <div
            className={rowClass}
            style={{ padding: `0.28rem 1.12rem 0.28rem ${paddingLeft}rem` }}
            onClick={onClick}
            data-path={path}
            draggable={!!path}
            onDragStart={onDragStart}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
        >
            {isFolder ? <ArrowIcon open={!!isOpen} /> : <span className={cls.ArrowSpacer} />}
            {isFolder ? <FolderIcon /> : <FileIcon />}
            <div className={cls.TreeItemMain}>
                <span className={cls.TreeItemLabel}>{name}</span>
                {subtitle && <span className={cls.TreeItemSubtitle}>{subtitle}</span>}
            </div>
            {isFolder && onCopyPath && (
                <Button
                    variant="ghost"
                    className={cls.CopyPathBtn}
                    onClick={e => {
                        e.stopPropagation()
                        onCopyPath()
                        setCopied(true)
                        setTimeout(() => setCopied(false), 1500)
                    }}
                    title={copied ? "Copied!" : "Copy path"}
                >
                    {copied ? <CheckIcon /> : <CopyIcon />}
                </Button>
            )}
            {isFolder && (
                <Button
                    variant="ghost"
                    className={cls.FolderAddBtn}
                    onClick={e => { e.stopPropagation(); onAddInFolder?.() }}
                    title="New note here"
                >
                    <svg
                        viewBox="0 0 12 12" width={11} height={11} fill="none" stroke="currentColor"
                        strokeWidth="1.6" strokeLinecap="round"
                    >
                        <path d="M6 1v10M1 6h10" />
                    </svg>
                </Button>
            )}
            {isFolder && (onDownloadFolder || onDeleteFolder) && (
                <KebabMenu
                    title="Folder actions"
                    items={[
                        ...(onDownloadFolder ? [{ label: "Download as .zip", onClick: onDownloadFolder }] : []),
                        ...(onDeleteFolder ? [{ label: "Delete folder", onClick: onDeleteFolder, danger: true }] : []),
                    ]}
                />
            )}
        </div>
    )
}
