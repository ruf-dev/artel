import { Button } from "@vervstack/chures"

import { useNotes } from "@/app/hooks/Notes.ts"
import { useDialog } from "@/app/hooks/Dialog.ts"
import { useBakeError } from "@/app/hooks/useErrorToast.ts"
import CreateNoteDialog from "@/pages/notes/components/CreateNoteDialog/CreateNoteDialog.tsx"
import FolderSection from "@/pages/notes/components/NotesSidebar/components/FolderSection/FolderSection.tsx"
import NotesSearchBar from "@/pages/notes/components/NotesSidebar/components/NotesSearchBar/NotesSearchBar.tsx"
import SearchResultsList from "@/pages/notes/components/NotesSidebar/components/SearchResultsList/SearchResultsList.tsx"
import { useNotesSearchQuery } from "@/pages/notes/components/NotesSidebar/processes/useNotesSearchQuery.ts"
import { useHighlightNote } from "@/pages/notes/components/NotesSidebar/processes/useHighlightNote.ts"
import LocateIcon from "@/pages/notes/components/icons/LocateIcon.tsx"
import cls from "@/pages/notes/components/NotesSidebar/NotesSidebar.module.css"

interface VaultOption {
    id: string
    name: string
}

interface NotesSidebarProps {
    vaults: VaultOption[]
    showCreateButton?: boolean
}

export default function NotesSidebar({ vaults, showCreateButton = true }: NotesSidebarProps) {
    const notesStore = useNotes()
    const { OpenDialog } = useDialog()
    const bakeError = useBakeError()
    const [searchQuery, setSearchQuery] = useNotesSearchQuery()
    const { highlightedPath, scrollAreaRef, highlightNote } = useHighlightNote()

    function handleVaultChange(e: React.ChangeEvent<HTMLSelectElement>) {
        void notesStore.selectVault(e.target.value)
    }

    function handleSelectNote(vid: string, path: string) {
        void notesStore.selectNote(vid, path)
    }

    function handleFindCurrentFile() {
        const path = notesStore.selectedPath
        if (!path) return
        if (searchQuery.trim()) setSearchQuery("")
        highlightNote(path)
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
                <div className={cls.SearchInputFlex}>
                    <NotesSearchBar/>
                </div>
                <Button
                    variant="ghost"
                    className={cls.FindCurrentFileBtn}
                    onClick={handleFindCurrentFile}
                    disabled={!notesStore.selectedPath}
                    title="Find currently open file"
                    aria-label="Find currently open file"
                >
                    <LocateIcon/>
                </Button>
            </div>
            <div className={cls.ScrollArea} ref={scrollAreaRef}>
                {!vaultId && (
                    <div className={cls.EmptyVaultHint}>Select a vault to browse notes</div>
                )}
                {vaultId && searchQuery.trim() && (
                    <SearchResultsList
                        notes={notesStore.notes}
                        searchQuery={searchQuery}
                        selectedPath={notesStore.selectedPath}
                        highlightedPath={highlightedPath}
                        onSelectNote={path => handleSelectNote(vaultId, path)}
                    />
                )}
                {vaultId && !searchQuery.trim() && (
                    <FolderSection
                        folders={notesStore.folders}
                        notes={notesStore.notes}
                        selectedPath={notesStore.selectedPath}
                        highlightedPath={highlightedPath}
                        revealPath={highlightedPath}
                        vaultId={vaultId}
                        onSelectNote={handleSelectNote}
                        onCreateNote={folderPath => handleCreateNote(folderPath)}
                        showCreateButton={showCreateButton}
                    />
                )}
            </div>
        </div>
    )
}
