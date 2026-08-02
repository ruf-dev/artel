import {useNotes} from "@/app/hooks/Notes.ts"
import {useDialog} from "@/app/hooks/Dialog.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import CreateNoteDialog from "@/pages/notes/components/CreateNoteDialog/CreateNoteDialog.tsx"
import FolderSection from "@/pages/notes/components/NotesSidebar/components/FolderSection/FolderSection.tsx"
import SidebarTopBar from "@/pages/notes/components/NotesSidebar/components/SidebarTopBar/SidebarTopBar.tsx"
import SearchResultsList from "@/pages/notes/components/NotesSidebar/components/SearchResultsList/SearchResultsList.tsx"
import SearchResultsTree from "@/pages/notes/components/NotesSidebar/components/SearchResultsTree/SearchResultsTree.tsx"
import {useNotesSearchQuery} from "@/pages/notes/components/NotesSidebar/processes/useNotesSearchQuery.ts"
import {useNotesTreeView} from "@/pages/notes/components/NotesSidebar/processes/useNotesTreeView.ts"
import {useHighlightNote} from "@/pages/notes/components/NotesSidebar/processes/useHighlightNote.ts"
import {useFolderActions} from "@/pages/notes/components/NotesSidebar/processes/useFolderActions.tsx"
import {useNotesDragAndDrop} from "@/pages/notes/components/NotesSidebar/processes/useNotesDragAndDrop.ts"
import cls from "@/pages/notes/components/NotesSidebar/NotesSidebar.module.css"

interface VaultOption {
    id: string
    name: string
    isPublic?: boolean
}

interface NotesSidebarProps {
    vaults: VaultOption[]
    showCreateButton?: boolean
}

export default function NotesSidebar({vaults, showCreateButton = true}: NotesSidebarProps) {
    const notesStore = useNotes()
    const {OpenDialog} = useDialog()
    const bakeError = useBakeError()
    const [searchQuery] = useNotesSearchQuery()
    const [treeView, toggleTreeView] = useNotesTreeView()
    const {scrollAreaRef} = useHighlightNote(notesStore.highlightedPath)
    const {handleDownloadFolder, handleDeleteFolder, handleUpload} = useFolderActions()
    const {onDragStart: onItemDragStart, onDrop: onFolderDrop} = useNotesDragAndDrop()

    function handleVaultChange(value: string[]) {
        if (!value[0]) return
        void notesStore.selectVault(value[0])
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
            <SidebarTopBar
                vaults={vaults}
                vaultId={vaultId}
                onVaultChange={handleVaultChange}
                treeView={treeView}
                onToggleTreeView={toggleTreeView}
            />
            <div className={cls.ScrollArea} ref={scrollAreaRef}>
                {!vaultId && (
                    <div className={cls.EmptyVaultHint}>Select a vault to browse notes</div>
                )}
                {vaultId && searchQuery.trim() && !treeView && (
                    <SearchResultsList
                        notes={notesStore.notes}
                        searchQuery={searchQuery}
                        selectedPath={notesStore.selectedPath}
                        highlightedPath={notesStore.highlightedPath}
                        onSelectNote={path => handleSelectNote(vaultId, path)}
                    />
                )}
                {vaultId && searchQuery.trim() && treeView && (
                    <SearchResultsTree
                        folders={notesStore.folders}
                        notes={notesStore.notes}
                        searchQuery={searchQuery}
                        selectedPath={notesStore.selectedPath}
                        highlightedPath={notesStore.highlightedPath}
                        onSelectNote={path => handleSelectNote(vaultId, path)}
                        onCreateNote={folderPath => handleCreateNote(folderPath)}
                    />
                )}
                {vaultId && !searchQuery.trim() && (
                    <FolderSection
                        folders={notesStore.folders}
                        notes={notesStore.notes}
                        selectedPath={notesStore.selectedPath}
                        highlightedPath={notesStore.highlightedPath}
                        revealPath={notesStore.highlightedPath}
                        vaultId={vaultId}
                        onSelectNote={handleSelectNote}
                        onCreateNote={folderPath => handleCreateNote(folderPath)}
                        onDownloadFolder={handleDownloadFolder}
                        onDeleteFolder={handleDeleteFolder}
                        onUpload={handleUpload}
                        showCreateButton={showCreateButton}
                        onItemDragStart={onItemDragStart}
                        onFolderDrop={onFolderDrop}
                    />
                )}
            </div>
        </div>
    )
}
