import { MouseEvent, useState } from "react"
import { Button } from "@vervstack/chures"

import KebabMenu from "@/components/atoms/KebabMenu/KebabMenu.tsx"
import CopyIcon from "@/pages/notes/components/icons/CopyIcon.tsx"
import CheckIcon from "@/pages/notes/components/icons/CheckIcon.tsx"
import AddInFolderIcon from "@/pages/notes/components/NotesSidebar/components/icons/AddInFolderIcon.tsx"
import cls from "@/pages/notes/components/NotesSidebar/components/NotesFolderActions/NotesFolderActions.module.css"

// The /notes folder-row action cluster, extracted out of the old TreeItem so it can be passed
// to the shared FileTree via `renderFolderTrailing`. Copy-path (with a transient "Copied!"
// state) and add-note-here are hover-revealed by FileTreeRow via the
// `data-reveal-on-row-hover` hook; the kebab (download-zip / delete-folder) is always visible
// and only rendered when at least one of those handlers is supplied (search-tree reuse omits
// them).
interface NotesFolderActionsProps {
    folderPath: string
    onCreateNoteInFolder: (folderPath: string) => void
    onDownloadFolder?: (folderPath: string) => void
    onDeleteFolder?: (folderPath: string) => void
}

export default function NotesFolderActions(props: NotesFolderActionsProps) {
    const { folderPath, onCreateNoteInFolder, onDownloadFolder, onDeleteFolder } = props
    const [copied, setCopied] = useState(false)

    const menuItems = [
        ...(onDownloadFolder
            ? [{ label: "Download as .zip", onClick: () => onDownloadFolder(folderPath) }]
            : []),
        ...(onDeleteFolder
            ? [{ label: "Delete folder", onClick: () => onDeleteFolder(folderPath), danger: true }]
            : []),
    ]

    function handleCopyPath(e: MouseEvent<HTMLButtonElement>) {
        e.stopPropagation()
        void navigator.clipboard.writeText(folderPath)
        setCopied(true)
        setTimeout(() => setCopied(false), 1500)
    }

    function handleAdd(e: MouseEvent<HTMLButtonElement>) {
        e.stopPropagation()
        onCreateNoteInFolder(folderPath)
    }

    return (
        <div className={cls.NotesFolderActionsContainer}>
            <Button
                variant="ghost"
                className={cls.CopyPathBtn}
                onClick={handleCopyPath}
                title={copied ? "Copied!" : "Copy path"}
                data-reveal-on-row-hover=""
            >
                {copied ? <CheckIcon /> : <CopyIcon />}
            </Button>
            <Button
                variant="ghost"
                className={cls.FolderAddBtn}
                onClick={handleAdd}
                title="New note here"
                data-reveal-on-row-hover=""
            >
                <AddInFolderIcon />
            </Button>
            {menuItems.length > 0 && <KebabMenu title="Folder actions" items={menuItems} />}
        </div>
    )
}
