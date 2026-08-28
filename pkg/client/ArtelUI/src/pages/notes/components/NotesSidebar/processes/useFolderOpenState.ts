import { useEffect, useState } from "react"

import { getAncestorFolderPaths } from "@/components/FileTree/fileTree.ts"

// Expand/collapse state for the /notes folder tree, plus auto-expand of every ancestor folder
// whenever `revealPath` changes (the reveal-in-tree behaviour). Split out of FolderSection.tsx
// to keep that component under the max-lines-per-function limit.
export function useFolderOpenState(revealPath: string | null) {
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

    useEffect(() => {
        if (!revealPath) return
        const ancestors = getAncestorFolderPaths(revealPath)
        if (ancestors.length === 0) return
        setOpenFolders(prev => {
            const next = new Set(prev)
            ancestors.forEach(a => next.add(a))
            return next
        })
    }, [revealPath])

    return { openFolders, toggleFolder }
}
