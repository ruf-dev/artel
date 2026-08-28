import {ReactNode, useState} from "react"

import FileTreeLeafRow from "@/components/FileTree/FileTreeLeafRow.tsx"
import FileTreeNode from "@/components/FileTree/FileTreeNode.tsx"
import {FileTreeRowDragProps} from "@/components/FileTree/FileTreeRow.tsx"
import {buildFolderTree, FileTreeItem, sortItemsByName} from "@/components/FileTree/fileTree.ts"
import cls from "@/components/FileTree/FileTree.module.css"

export interface FileTreeProps<T extends FileTreeItem> {
    folders: string[]
    items: T[]
    isActive?: (path: string) => boolean
    isHighlighted?: (path: string) => boolean
    onSelectItem: (path: string) => void
    renderFolderTrailing?: (folderPath: string) => ReactNode
    itemDragProps?: (path: string, isFolder: boolean) => FileTreeRowDragProps
    // Escape hatch for consumers that need to drive expand/collapse externally
    // (e.g. notes' auto-expand-on-reveal). Omit both for the uncontrolled default.
    controlledOpenFolders?: Set<string>
    onOpenFoldersChange?: (next: Set<string>) => void
}

export default function FileTree<T extends FileTreeItem>(props: FileTreeProps<T>) {
    const {folders, items, isActive, isHighlighted, onSelectItem} = props
    const {renderFolderTrailing, itemDragProps, controlledOpenFolders, onOpenFoldersChange} = props
    const [uncontrolledOpen, setUncontrolledOpen] = useState<Set<string>>(new Set())
    const openFolders = controlledOpenFolders ?? uncontrolledOpen

    function toggleFolder(path: string) {
        const next = new Set(openFolders)
        if (next.has(path)) {
            next.delete(path)
        } else {
            next.add(path)
        }
        if (!controlledOpenFolders) {
            setUncontrolledOpen(next)
        }
        onOpenFoldersChange?.(next)
    }

    const tree = buildFolderTree(folders)
    const rootItems = sortItemsByName(items.filter(i => i.path && !i.path.includes("/")))

    return (
        <div className={cls.FileTreeContainer}>
            {tree.map(node => (
                <FileTreeNode
                    key={node.path}
                    node={node}
                    items={items}
                    openFolders={openFolders}
                    depth={0}
                    isActive={isActive}
                    isHighlighted={isHighlighted}
                    onToggleFolder={toggleFolder}
                    onSelectItem={onSelectItem}
                    renderFolderTrailing={renderFolderTrailing}
                    itemDragProps={itemDragProps}
                />
            ))}
            {rootItems.map(item => (
                <FileTreeLeafRow
                    key={item.path}
                    item={item}
                    depth={0}
                    isActive={isActive}
                    isHighlighted={isHighlighted}
                    onSelectItem={onSelectItem}
                    itemDragProps={itemDragProps}
                />
            ))}
        </div>
    )
}
