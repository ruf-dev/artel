import { useState } from "react"
import {Input} from "@vervstack/chures"

import { useNotes } from "@/app/hooks/Notes.ts"
import { useDialog } from "@/app/hooks/Dialog.ts"
import { useBakeError } from "@/app/hooks/useErrorToast.ts"
import CreateNoteDialog from "@/pages/notes/components/CreateNoteDialog/CreateNoteDialog.tsx"
import SearchIcon from "@/pages/notes/components/icons/SearchIcon.tsx"
import TreeItem from "@/pages/notes/components/NotesSidebar/components/TreeItem/TreeItem.tsx"
import FolderSection from "@/pages/notes/components/NotesSidebar/components/FolderSection/FolderSection.tsx"
import { getNoteName } from "@/pages/notes/components/NotesSidebar/processes/notesTreeHelpers.ts"
import cls from "@/pages/notes/components/NotesSidebar/NotesSidebar.module.css"


interface VaultOption {
    id: string
    name: string
}

interface NotesSidebarProps {
    vaults: VaultOption[]
}

export default function NotesSidebar({ vaults }: NotesSidebarProps) {
    const notesStore = useNotes()
    const { OpenDialog } = useDialog()
    const bakeError = useBakeError()
    const [searchQuery, setSearchQuery] = useState("")

    function handleVaultChange(e: React.ChangeEvent<HTMLSelectElement>) {
        void notesStore.selectVault(e.target.value)
    }

    function handleSelectNote(vid: string, path: string) {
        void notesStore.selectNote(vid, path)
    }

    function handleCreateNote(folderPath?: string) {
        OpenDialog(
            <CreateNoteDialog
                initialPath={folderPath ? folderPath + "/" : ""}
                folders={notesStore.folders}
                onConfirm={async (path: string) => {
                    try {
                        await notesStore.createNote(path)
                    } catch (err) {
                        bakeError("Failed to create note", err)
                    }
                }}
            />
        )
    }

    const vaultId = notesStore.vaultId

    return (
        <div className={cls.NotesSidebarContainer}>
            <div className={cls.VaultPickerWrapper}>
                <select
                    className={cls.VaultSelect}
                    value={vaultId ?? ""}
                    onChange={handleVaultChange}
                >
                    <option value="" disabled className={cls.VaultSelectPlaceholder}>
                        Select vault…
                    </option>
                    {vaults.map(v => (
                        <option key={v.id} value={v.id}>{v.name}</option>
                    ))}
                </select>
            </div>
            <div className={cls.SearchWrapper}>
                <div className={cls.SearchBar}>
                    <SearchIcon className={cls.SearchIcon}/>
                    <Input
                        className={cls.SearchInput}
                        placeholder="Search notes…"
                        value={searchQuery}
                        setValue={setSearchQuery}
                    />
                </div>
            </div>
            <div className={cls.ScrollArea}>
                {!vaultId && (
                    <div className={cls.EmptyVaultHint}>Select a vault to browse notes</div>
                )}
                {vaultId && searchQuery.trim() && notesStore.notes
                    .filter(n => getNoteName(n).toLowerCase().includes(searchQuery.trim().toLowerCase()))
                    .map(note => {
                        const dir = note.path ? note.path.split("/").slice(0, -1).join("/") : undefined
                        return (
                            <TreeItem
                                key={note.path}
                                name={getNoteName(note)}
                                subtitle={dir || undefined}
                                active={notesStore.selectedPath === note.path}
                                onClick={() => note.path && handleSelectNote(vaultId, note.path)}
                            />
                        )
                    })
                }
                {vaultId && !searchQuery.trim() && (
                    <FolderSection
                        folders={notesStore.folders}
                        notes={notesStore.notes}
                        selectedPath={notesStore.selectedPath}
                        vaultId={vaultId}
                        onSelectNote={handleSelectNote}
                        onCreateNote={folderPath => handleCreateNote(folderPath)}
                    />
                )}
            </div>
        </div>
    )
}
