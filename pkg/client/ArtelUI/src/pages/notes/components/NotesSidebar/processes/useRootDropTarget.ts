import { DragEvent, useState } from "react"

// Local drag-over visual state + handlers for the sticky "All Notes" header, which doubles as
// the vault-root drop target for the notes sidebar's drag-and-drop move. Split out of
// FolderSection.tsx to keep that component under the max-lines-per-function limit.
export function useRootDropTarget(onFolderDrop: (targetFolderPath: string) => (e: DragEvent<HTMLDivElement>) => void) {
    const [isDragOver, setIsDragOver] = useState(false)

    function onDragOver(e: DragEvent<HTMLDivElement>) {
        e.preventDefault()
        setIsDragOver(true)
    }

    function onDragLeave() {
        setIsDragOver(false)
    }

    function onDrop(e: DragEvent<HTMLDivElement>) {
        e.preventDefault()
        setIsDragOver(false)
        onFolderDrop("")(e)
    }

    return { isDragOver, onDragOver, onDragLeave, onDrop }
}
