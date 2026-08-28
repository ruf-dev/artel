import {ReactNode} from "react"

import FileTreeLeafRow from "@/components/FileTree/FileTreeLeafRow.tsx"
import FileTreeRow, {FileTreeRowDragProps} from "@/components/FileTree/FileTreeRow.tsx"
import {FileTreeItem, FolderNode, getDirectNotes} from "@/components/FileTree/fileTree.ts"

export interface FileTreeNodeProps<T extends FileTreeItem> {
    node: FolderNode
    items: T[]
    openFolders: Set<string>
    depth: number
    // Replaces docs' `selectedPath === path`; a predicate so single- and multi-select
    // consumers (workbench) can both be served.
    isActive?: (path: string) => boolean
    isHighlighted?: (path: string) => boolean
    onToggleFolder: (path: string) => void
    onSelectItem: (path: string) => void
    // Notes passes its kebab / copy-path / add-in-folder cluster here.
    renderFolderTrailing?: (folderPath: string) => ReactNode
    // Optional; notes wires drag-and-drop, docs/workbench omit it.
    itemDragProps?: (path: string, isFolder: boolean) => FileTreeRowDragProps
}

export default function FileTreeNode<T extends FileTreeItem>(props: FileTreeNodeProps<T>) {
    const {node, items, openFolders, depth, isActive, isHighlighted} = props
    const {onToggleFolder, onSelectItem, renderFolderTrailing, itemDragProps} = props
    const isOpen = openFolders.has(node.path)
    const directItems = getDirectNotes(node.path, items)
    const folderDrag = itemDragProps?.(node.path, true)

    return (
        <>
            <FileTreeRow
                name={node.name}
                path={node.path}
                depth={depth}
                isFolder
                isOpen={isOpen}
                onClick={() => onToggleFolder(node.path)}
                trailing={renderFolderTrailing?.(node.path)}
                draggable={folderDrag?.draggable}
                onDragStart={folderDrag?.onDragStart}
                onDropTarget={folderDrag?.onDropTarget}
            />
            {isOpen && (
                <>
                    {node.children.map(child => (
                        <FileTreeNode
                            key={child.path}
                            node={child}
                            items={items}
                            openFolders={openFolders}
                            depth={depth + 1}
                            isActive={isActive}
                            isHighlighted={isHighlighted}
                            onToggleFolder={onToggleFolder}
                            onSelectItem={onSelectItem}
                            renderFolderTrailing={renderFolderTrailing}
                            itemDragProps={itemDragProps}
                        />
                    ))}
                    {directItems.map(item => (
                        <FileTreeLeafRow
                            key={item.path}
                            item={item}
                            depth={depth + 1}
                            isActive={isActive}
                            isHighlighted={isHighlighted}
                            onSelectItem={onSelectItem}
                            itemDragProps={itemDragProps}
                        />
                    ))}
                </>
            )}
        </>
    )
}
