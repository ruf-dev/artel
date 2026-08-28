import { PublicDocsNoteItem } from "@/app/api/artel"
import FileTree from "@/components/FileTree/FileTree.tsx"
import cls from "@/pages/docs/components/DocsSidebar/DocsSidebar.module.css"

interface DocsSidebarProps {
    vaultName: string | null
    folders: string[]
    notes: PublicDocsNoteItem[]
    selectedPath: string | null
    onSelectNote: (path: string) => void
}

// Public read-only tree: the shared FileTree component drives folder expand/collapse and note
// selection. No drag-and-drop, no folder actions — those belong to the authenticated /notes
// editor. PublicDocsNoteItem ({ path?, mtime? }) already satisfies FileTree's FileTreeItem.
export default function DocsSidebar(
    { vaultName, folders, notes, selectedPath, onSelectNote }: DocsSidebarProps
) {
    return (
        <div className={cls.DocsSidebarContainer}>
            <div className={cls.Header}>{vaultName ?? "Docs"}</div>
            <div className={cls.ScrollArea}>
                <FileTree
                    folders={folders}
                    items={notes}
                    isActive={path => path === selectedPath}
                    onSelectItem={onSelectNote}
                />
            </div>
        </div>
    )
}
