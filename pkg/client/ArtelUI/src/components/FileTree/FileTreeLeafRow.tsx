import {FileTreeItem, getItemName} from "@/components/FileTree/fileTree.ts"
import FileTreeRow, {FileTreeRowDragProps} from "@/components/FileTree/FileTreeRow.tsx"

// A leaf (file) row: resolves the display name, active/highlighted state via predicates
// and the optional per-item drag props, then renders a FileTreeRow. Extracted so the
// root-level list in FileTree.tsx and the in-folder list in FileTreeNode.tsx share one
// implementation (and to keep both call sites under the jsx-max-depth limit).
export interface FileTreeLeafRowProps<T extends FileTreeItem> {
    item: T
    depth: number
    isActive?: (path: string) => boolean
    isHighlighted?: (path: string) => boolean
    onSelectItem: (path: string) => void
    itemDragProps?: (path: string, isFolder: boolean) => FileTreeRowDragProps
}

export default function FileTreeLeafRow<T extends FileTreeItem>(props: FileTreeLeafRowProps<T>) {
    const {item, depth, isActive, isHighlighted, onSelectItem, itemDragProps} = props
    const path = item.path
    const drag = path ? itemDragProps?.(path, false) : undefined

    return (
        <FileTreeRow
            name={getItemName(item)}
            path={path}
            depth={depth}
            active={!!path && !!isActive?.(path)}
            highlighted={!!path && !!isHighlighted?.(path)}
            onClick={() => path && onSelectItem(path)}
            draggable={drag?.draggable}
            onDragStart={drag?.onDragStart}
            onDropTarget={drag?.onDropTarget}
        />
    )
}
