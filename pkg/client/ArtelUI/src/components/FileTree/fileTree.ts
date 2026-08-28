// Shared, generic file-tree helpers.
//
// This supersedes the two page-local copies that Stage 2 removes once /docs and /notes
// migrate onto this component:
//   - pages/docs/processes/docsTree.ts
//   - pages/notes/components/NotesSidebar/processes/notesTreeHelpers.ts
// The algorithm is lifted verbatim from docsTree.ts (which already matched
// notesTreeHelpers.ts) and generalised over the item type via `T extends FileTreeItem`
// so docs (PublicDocsNoteItem), notes (NoteItem) and a future workbench consumer can
// all share it.

export interface FileTreeItem {
    path?: string
    mtime?: string
}

export interface FolderNode {
    name: string
    path: string
    children: FolderNode[]
}

function byName(a: {name: string}, b: {name: string}): number {
    return a.name.localeCompare(b.name)
}

function sortTree(nodes: FolderNode[]): FolderNode[] {
    for (const node of nodes) {
        sortTree(node.children)
    }
    return nodes.sort(byName)
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
        nodeMap.set(path, {name: parts[parts.length - 1], path, children: []})
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

    return sortTree(roots)
}

export function getItemName(item: FileTreeItem): string {
    if (!item.path) return "Untitled"
    const parts = item.path.split("/")
    const filename = parts[parts.length - 1] || item.path
    return filename.endsWith(".md") ? filename.slice(0, -3) : filename
}

export function sortItemsByName<T extends FileTreeItem>(items: T[]): T[] {
    return [...items].sort((a, b) => getItemName(a).localeCompare(getItemName(b)))
}

export function getDirectNotes<T extends FileTreeItem>(folderPath: string, items: T[]): T[] {
    const prefix = folderPath + "/"
    const direct = items.filter(n => {
        if (!n.path || !n.path.startsWith(prefix)) return false
        return !n.path.slice(prefix.length).includes("/")
    })
    return sortItemsByName(direct)
}

// Every ancestor folder path of a note path, outermost first
// (e.g. "a/b/c/note.md" -> ["a", "a/b", "a/b/c"]). Notes uses this to auto-expand
// folders on its reveal-in-tree behaviour; kept here so that logic can move onto the
// shared component too.
export function getAncestorFolderPaths(notePath: string): string[] {
    const parts = notePath.split("/")
    const ancestors: string[] = []
    for (let i = 1; i < parts.length; i++) {
        ancestors.push(parts.slice(0, i).join("/"))
    }
    return ancestors
}
