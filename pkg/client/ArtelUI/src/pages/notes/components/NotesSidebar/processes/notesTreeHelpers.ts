import { NoteItem } from "@/app/hooks/Notes.ts"

export interface FolderNode {
    name: string
    path: string
    children: FolderNode[]
}

export function buildFolderTree(folders: string[]): FolderNode[] {
    const nodeMap = new Map<string, FolderNode>()

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

    const roots: FolderNode[] = []
    for (const [path, node] of nodeMap) {
        const lastSlash = path.lastIndexOf("/")
        if (lastSlash === -1) {
            roots.push(node)
        } else {
            nodeMap.get(path.slice(0, lastSlash))?.children.push(node)
        }
    }

    return roots
}

export function getNoteName(note: NoteItem): string {
    if (!note.path) return "Untitled"
    const parts = note.path.split("/")
    const filename = parts[parts.length - 1] || note.path
    return filename.endsWith(".md") ? filename.slice(0, -3) : filename
}

export function getDirectNotes(folderPath: string, notes: NoteItem[]): NoteItem[] {
    const prefix = folderPath + "/"
    return notes.filter(n => {
        if (!n.path || !n.path.startsWith(prefix)) return false
        return !n.path.slice(prefix.length).includes("/")
    })
}
