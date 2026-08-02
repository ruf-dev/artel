import {PublicDocsNoteItem} from "@/app/api/artel"

// Minimal local reimplementation of
// pages/notes/components/NotesSidebar/processes/notesTreeHelpers.ts's buildFolderTree /
// sortNotesByName / getNoteName — duplicated here per the colocation rule (that file lives
// under pages/notes/**, local to NotesPage) rather than cross-imported.

export interface DocsFolderNode {
    name: string
    path: string
    children: DocsFolderNode[]
}

function byName(a: { name: string }, b: { name: string }): number {
    return a.name.localeCompare(b.name)
}

function sortTree(nodes: DocsFolderNode[]): DocsFolderNode[] {
    for (const node of nodes) {
        sortTree(node.children)
    }
    return nodes.sort(byName)
}

export function buildFolderTree(folders: string[]): DocsFolderNode[] {
    const nodeMap = new Map<string, DocsFolderNode>()

    const allPaths = new Set<string>()
    for (const folder of folders) {
        const parts = folder.split("/")
        for (let i = 1; i <= parts.length; i++) {
            allPaths.add(parts.slice(0, i).join("/"))
        }
    }

    for (const path of allPaths) {
        const parts = path.split("/")
        nodeMap.set(path, { name: parts[parts.length - 1], path, children: [] })
    }

    const roots: DocsFolderNode[] = []
    for (const [path, node] of nodeMap) {
        const lastSlash = path.lastIndexOf("/")
        if (lastSlash === -1) {
            roots.push(node)
        } else {
            nodeMap.get(path.slice(0, lastSlash))?.children.push(node)
        }
    }

    return sortTree(roots)
}

export function getNoteName(note: PublicDocsNoteItem): string {
    if (!note.path) return "Untitled"
    const parts = note.path.split("/")
    const filename = parts[parts.length - 1] || note.path
    return filename.endsWith(".md") ? filename.slice(0, -3) : filename
}

export function sortNotesByName(notes: PublicDocsNoteItem[]): PublicDocsNoteItem[] {
    return [...notes].sort((a, b) => getNoteName(a).localeCompare(getNoteName(b)))
}

export function getDirectNotes(folderPath: string, notes: PublicDocsNoteItem[]): PublicDocsNoteItem[] {
    const prefix = folderPath + "/"
    const direct = notes.filter(n => {
        if (!n.path || !n.path.startsWith(prefix)) return false
        return !n.path.slice(prefix.length).includes("/")
    })
    return sortNotesByName(direct)
}
