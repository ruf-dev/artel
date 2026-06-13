import { useState } from "react"
import cls from "./NotesSidebar.module.css"
import { NoteItem, useNotes } from "@/app/hooks/Notes.ts"
import { useDialog } from "@/app/hooks/Dialog.ts"
import { useBakeError } from "@/app/hooks/useErrorToast.ts"
import CreateNoteDialog from "@/pages/notes/components/CreateNoteDialog/CreateNoteDialog.tsx"

interface VaultOption {
    id: string
    name: string
}

interface NotesSidebarProps {
    vaults: VaultOption[]
}

function ArrowIcon({ open }: { open: boolean }) {
    return (
        <svg
            viewBox="0 0 8 8"
            width={8}
            height={8}
            style={{ flexShrink: 0, opacity: 0.35, transform: open ? "rotate(90deg)" : "none", transition: "transform 0.15s" }}
        >
            <path d="M2 1.5l4 2.5-4 2.5z" fill="currentColor" />
        </svg>
    )
}

function FileIcon() {
    return (
        <svg viewBox="0 0 14 14" width={10} height={10} fill="none" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" style={{ flexShrink: 0, opacity: 0.35 }}>
            <rect x="2" y="1" width="10" height="12" rx="1.5" />
            <path d="M5 4.5h4M5 7h4M5 9.5h2" />
        </svg>
    )
}

function FolderIcon() {
    return (
        <svg viewBox="0 0 14 14" width={10} height={10} fill="none" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" style={{ flexShrink: 0, opacity: 0.45 }}>
            <path d="M1 3.5h4l1.5 1.5H13v7H1z" />
        </svg>
    )
}

interface TreeItemProps {
    name: string
    active?: boolean
    depth?: number
    isFolder?: boolean
    isOpen?: boolean
    onClick?: () => void
}

function TreeItem({ name, active, depth = 0, isFolder, isOpen, onClick }: TreeItemProps) {
    const paddingLeft = 16 + depth * 12
    const rowClass = `${cls.TreeItemRow}${active ? ` ${cls.TreeItemRowActive}` : ""}`

    return (
        <div className={rowClass} style={{ padding: `4px 16px 4px ${paddingLeft}px` }} onClick={onClick}>
            {isFolder ? <ArrowIcon open={!!isOpen} /> : <span className={cls.ArrowSpacer} />}
            {isFolder ? <FolderIcon /> : <FileIcon />}
            <span className={cls.TreeItemLabel}>{name}</span>
        </div>
    )
}

interface FolderSectionProps {
    folders: string[]
    notes: NoteItem[]
    selectedPath: string | null
    vaultId: string
    onSelectNote: (vaultId: string, path: string) => void
    onCreateNote: () => void
}

function FolderSection({ folders, notes, selectedPath, vaultId, onSelectNote, onCreateNote }: FolderSectionProps) {
    const [openFolders, setOpenFolders] = useState<Set<string>>(new Set())

    function toggleFolder(folder: string) {
        setOpenFolders(prev => {
            const next = new Set(prev)
            if (next.has(folder)) {
                next.delete(folder)
            } else {
                next.add(folder)
            }
            return next
        })
    }

    const rootNotes = notes.filter(n => n.path && !n.path.includes("/"))
    const folderNotes = (folder: string) => notes.filter(n => n.path && n.path.startsWith(folder + "/"))

    function getNoteName(note: NoteItem): string {
        if (!note.path) return "Untitled"
        const parts = note.path.split("/")
        return parts[parts.length - 1] || note.path
    }

    return (
        <>
            <div className={cls.SectionHeader}>
                <span className={cls.SectionLabel}>All Notes</span>
                <button className={cls.CreateNoteBtn} onClick={onCreateNote} title="New note">
                    <svg viewBox="0 0 12 12" width={10} height={10} fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round">
                        <path d="M6 1v10M1 6h10" />
                    </svg>
                </button>
            </div>
            {folders.map(folder => (
                <div key={folder}>
                    <TreeItem
                        name={folder}
                        isFolder
                        isOpen={openFolders.has(folder)}
                        onClick={() => toggleFolder(folder)}
                    />
                    {openFolders.has(folder) && folderNotes(folder).map(note => (
                        <TreeItem
                            key={note.path}
                            name={getNoteName(note)}
                            depth={1}
                            active={selectedPath === note.path}
                            onClick={() => note.path && onSelectNote(vaultId, note.path)}
                        />
                    ))}
                </div>
            ))}
            {rootNotes.map(note => (
                <TreeItem
                    key={note.path}
                    name={getNoteName(note)}
                    active={selectedPath === note.path}
                    onClick={() => note.path && onSelectNote(vaultId, note.path)}
                />
            ))}
        </>
    )
}

export default function NotesSidebar({ vaults }: NotesSidebarProps) {
    const { vaultId, folders, notes, selectedPath, selectVault, selectNote, createNote } = useNotes()
    const { OpenDialog } = useDialog()
    const bakeError = useBakeError()

    function handleVaultChange(e: React.ChangeEvent<HTMLSelectElement>) {
        void selectVault(e.target.value)
    }

    function handleSelectNote(vid: string, path: string) {
        void selectNote(vid, path)
    }

    function handleCreateNote() {
        OpenDialog(
            <CreateNoteDialog
                onConfirm={async (path: string) => {
                    try {
                        await createNote(path)
                    } catch (err) {
                        bakeError("Failed to create note", err)
                    }
                }}
            />
        )
    }

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
                    <svg viewBox="0 0 16 16" width={12} height={12} fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" className={cls.SearchIcon}>
                        <circle cx="6.5" cy="6.5" r="4.5" />
                        <path d="M10.5 10.5l3 3" />
                    </svg>
                    <span className={cls.SearchPlaceholder}>Search notes…</span>
                </div>
            </div>
            <div className={cls.ScrollArea}>
                {!vaultId && (
                    <div className={cls.EmptyVaultHint}>Select a vault to browse notes</div>
                )}
                {vaultId && (
                    <FolderSection
                        folders={folders}
                        notes={notes}
                        selectedPath={selectedPath}
                        vaultId={vaultId}
                        onSelectNote={handleSelectNote}
                        onCreateNote={handleCreateNote}
                    />
                )}
            </div>
        </div>
    )
}
