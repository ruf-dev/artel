import { DragEvent } from "react"

import { useNotes } from "@/app/hooks/Notes.ts"
import { useBakeError } from "@/app/hooks/useErrorToast.ts"

// Custom drag payload MIME types (rather than a single JSON blob) so drag-start can set two
// primitive fields without needing to parse anything back out on drop.
const PATH_TYPE = "application/x-artel-note-path"
const IS_FOLDER_TYPE = "application/x-artel-note-is-folder"

function getParentFolder(path: string): string {
    const lastSlash = path.lastIndexOf("/")
    return lastSlash === -1 ? "" : path.slice(0, lastSlash)
}

function getBasename(path: string): string {
    const parts = path.split("/")
    return parts[parts.length - 1] || path
}

// A drop is invalid (and must be silently ignored rather than sent to the API) when it's a
// no-op — dropping onto the item's current parent — or, for folders, when it would create a
// cycle: dropping a folder onto itself or one of its own descendants.
function isInvalidDrop(draggedPath: string, isFolder: boolean, targetFolderPath: string): boolean {
    if (targetFolderPath === getParentFolder(draggedPath)) return true
    if (!isFolder) return false
    return targetFolderPath === draggedPath || targetFolderPath.startsWith(`${draggedPath}/`)
}

// Orchestrates drag-and-drop moves for the Notes sidebar tree: computing the destination path,
// rejecting invalid drops client-side (self/parent/descendant), dispatching the store's
// moveEntry action, and surfacing failures (e.g. a name conflict at the destination) as an
// error toast. Instantiated once in NotesSidebar.tsx and threaded down as callback props,
// following this repo's rule that exported components receive callbacks rather than reaching
// into the store themselves.
export function useNotesDragAndDrop() {
    const moveEntry = useNotes(state => state.moveEntry)
    const bakeError = useBakeError()

    function onDragStart(path: string, isFolder: boolean) {
        return (e: DragEvent<HTMLDivElement>) => {
            e.dataTransfer.setData(PATH_TYPE, path)
            e.dataTransfer.setData(IS_FOLDER_TYPE, isFolder ? "1" : "")
            e.dataTransfer.effectAllowed = "move"
        }
    }

    function onDrop(targetFolderPath: string) {
        return (e: DragEvent<HTMLDivElement>) => {
            e.preventDefault()

            const draggedPath = e.dataTransfer.getData(PATH_TYPE)
            if (!draggedPath) return

            const isFolder = e.dataTransfer.getData(IS_FOLDER_TYPE) === "1"
            if (isInvalidDrop(draggedPath, isFolder, targetFolderPath)) return

            const newPath = targetFolderPath
                ? `${targetFolderPath}/${getBasename(draggedPath)}`
                : getBasename(draggedPath)
            if (newPath === draggedPath) return

            moveEntry(draggedPath, newPath, isFolder)
                .catch(err => bakeError(isFolder ? "Failed to move folder" : "Failed to move note", err))
        }
    }

    return { onDragStart, onDrop }
}
