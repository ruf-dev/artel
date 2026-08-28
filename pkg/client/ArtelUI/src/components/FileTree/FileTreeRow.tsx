import {DragEvent, ReactNode, useState} from "react"

import {cn} from "@/app/utils/cn.ts"
import ArrowIcon from "@/components/FileTree/icons/ArrowIcon.tsx"
import FolderIcon from "@/components/FileTree/icons/FolderIcon.tsx"
import FileIcon from "@/components/FileTree/icons/FileIcon.tsx"
import cls from "@/components/FileTree/FileTreeRow.module.css"

// The drag-related props a caller wires per row (notes does; docs/workbench omit).
// Bundled so `itemDragProps(path, isFolder)` can return them in one object.
export interface FileTreeRowDragProps {
    draggable?: boolean
    onDragStart?: (e: DragEvent<HTMLDivElement>) => void
    onDropTarget?: (e: DragEvent<HTMLDivElement>) => void
}

export interface FileTreeRowProps {
    name: string
    path?: string
    subtitle?: string
    depth?: number
    isFolder?: boolean
    isOpen?: boolean
    active?: boolean
    highlighted?: boolean
    onClick?: () => void
    // Right-aligned folder-action slot (kebab / copy-path / add-in-folder). Notes passes
    // these; docs and workbench pass nothing. Kept a raw ReactNode slot so this row stays
    // a tier-3 leaf with no KebabMenu import.
    trailing?: ReactNode
    draggable?: boolean
    onDragStart?: (e: DragEvent<HTMLDivElement>) => void
    // When set together with isFolder, the row accepts a drop and shows a drop-target
    // highlight while a drag hovers over it.
    onDropTarget?: (e: DragEvent<HTMLDivElement>) => void
}

export default function FileTreeRow(props: FileTreeRowProps) {
    const {name, path, subtitle, depth, isFolder, isOpen} = props
    const {active, highlighted, onClick, trailing, draggable} = props
    const {onDragStart, onDropTarget} = props
    const [isDragOver, setIsDragOver] = useState(false)
    const paddingLeft = 1.12 + (depth ?? 0) * 0.84
    const acceptsDrop = !!isFolder && !!onDropTarget
    const rowClass = cn(
        cls.FileTreeRowContainer,
        active && cls.FileTreeRowActive,
        highlighted && cls.FileTreeRowHighlight,
        isDragOver && cls.FileTreeRowDropTarget,
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
            style={{padding: `0.28rem 1.12rem 0.28rem ${paddingLeft}rem`}}
            onClick={onClick}
            data-path={path}
            draggable={draggable}
            onDragStart={onDragStart}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
        >
            {isFolder ? <ArrowIcon open={!!isOpen}/> : <span className={cls.ArrowSpacer}/>}
            {isFolder ? <FolderIcon/> : <FileIcon/>}
            <div className={cls.Main}>
                <span className={cls.Label}>{name}</span>
                {subtitle ? <span className={cls.Subtitle}>{subtitle}</span> : null}
            </div>
            {trailing}
        </div>
    )
}
