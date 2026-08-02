import { useState } from "react"

import { PublicDocsNoteItem } from "@/app/api/artel"
import DocsFolderNode from "@/pages/docs/components/DocsFolderNode/DocsFolderNode.tsx"
import DocsTreeItem from "@/pages/docs/components/DocsTreeItem/DocsTreeItem.tsx"
import { buildFolderTree, getNoteName, sortNotesByName } from "@/pages/docs/processes/docsTree.ts"
import cls from "@/pages/docs/components/DocsSidebar/DocsSidebar.module.css"

interface DocsSidebarProps {
    vaultName: string | null
    folders: string[]
    notes: PublicDocsNoteItem[]
    selectedPath: string | null
    onSelectNote: (path: string) => void
}

export default function DocsSidebar(
    { vaultName, folders, notes, selectedPath, onSelectNote }: DocsSidebarProps
) {
    const [openFolders, setOpenFolders] = useState<Set<string>>(new Set())

    function toggleFolder(path: string) {
        setOpenFolders(prev => {
            const next = new Set(prev)
            if (next.has(path)) {
                next.delete(path)
            } else {
                next.add(path)
            }
            return next
        })
    }

    const tree = buildFolderTree(folders)
    const rootNotes = sortNotesByName(notes.filter(n => n.path && !n.path.includes("/")))

    return (
        <div className={cls.DocsSidebarContainer}>
            <div className={cls.Header}>{vaultName ?? "Docs"}</div>
            <div className={cls.ScrollArea}>
                {tree.map(node => (
                    <DocsFolderNode
                        key={node.path}
                        node={node}
                        notes={notes}
                        openFolders={openFolders}
                        selectedPath={selectedPath}
                        depth={0}
                        onToggle={toggleFolder}
                        onSelectNote={onSelectNote}
                    />
                ))}
                {rootNotes.map(note => (
                    <DocsTreeItem
                        key={note.path}
                        name={getNoteName(note)}
                        path={note.path}
                        active={selectedPath === note.path}
                        onClick={() => note.path && onSelectNote(note.path)}
                    />
                ))}
            </div>
        </div>
    )
}
