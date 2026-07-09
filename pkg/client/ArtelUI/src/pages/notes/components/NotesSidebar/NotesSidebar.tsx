import { useState } from "react"
import {Button} from "@vervstack/chures"

import { NoteItem, useNotes } from "@/app/hooks/Notes.ts"
import { useDialog } from "@/app/hooks/Dialog.ts"
import { useBakeError } from "@/app/hooks/useErrorToast.ts"
import CreateNoteDialog from "@/pages/notes/components/CreateNoteDialog/CreateNoteDialog.tsx"
import PlusIcon from "@/pages/notes/components/icons/PlusIcon.tsx"
import SearchIcon from "@/pages/notes/components/icons/SearchIcon.tsx"

import cls from "./NotesSidebar.module.css"


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
            width={11}
            height={11}
            style={{ flexShrink: 0, opacity: 0.35, transform: open ? "rotate(90deg)" : "none", transition: "transform 0.15s" }}
        >
            <path d="M2 1.5l4 2.5-4 2.5z" fill="currentColor" />
        </svg>
    )
}

function FileIcon() {
    return (
        <svg viewBox="0 0 14 14" width={14} height={14} fill="none" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" style={{ flexShrink: 0, opacity: 0.35 }}>
            <rect x="2" y="1" width="10" height="12" rx="1.5" />
            <path d="M5 4.5h4M5 7h4M5 9.5h2" />
        </svg>
    )
}

function FolderIcon() {
    return (
        <svg viewBox="0 0 14 14" width={14} height={14} fill="none" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" style={{ flexShrink: 0, opacity: 0.45 }}>
            <path d="M1 3.5h4l1.5 1.5H13v7H1z" />
        </svg>
    )
}

interface TreeItemProps {
    name: string
    subtitle?: string
    active?: boolean
    depth?: number
    isFolder?: boolean
    isOpen?: boolean
    onClick?: () => void
    onAddInFolder?: () => void
}

function TreeItem({ name, subtitle, active, depth = 0, isFolder, isOpen, onClick, onAddInFolder }: TreeItemProps) {
    const paddingLeft = 1.12 + depth * 0.84
    const rowClass = `${cls.TreeItemRow}${active ? ` ${cls.TreeItemRowActive}` : ""}`

    return (
        <div className={rowClass} style={{ padding: `0.28rem 1.12rem 0.28rem ${paddingLeft}rem` }} onClick={onClick}>
            {isFolder ? <ArrowIcon open={!!isOpen} /> : <span className={cls.ArrowSpacer} />}
            {isFolder ? <FolderIcon /> : <FileIcon />}
            <div className={cls.TreeItemMain}>
                <span className={cls.TreeItemLabel}>{name}</span>
                {subtitle && <span className={cls.TreeItemSubtitle}>{subtitle}</span>}
            </div>
            {isFolder && (
                <Button
                    variant="ghost"
                    className={cls.FolderAddBtn}
                    onClick={e => { e.stopPropagation(); onAddInFolder?.() }}
                    title="New note here"
                >
                    <svg viewBox="0 0 12 12" width={11} height={11} fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round">
                        <path d="M6 1v10M1 6h10" />
                    </svg>
                </Button>
            )}
        </div>
    )
}

interface FolderNode {
    name: string
    path: string
    children: FolderNode[]
}

function buildFolderTree(folders: string[]): FolderNode[] {
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

function getNoteName(note: NoteItem): string {
    if (!note.path) return "Untitled"
    const parts = note.path.split("/")
    const filename = parts[parts.length - 1] || note.path
    return filename.endsWith(".md") ? filename.slice(0, -3) : filename
}

function getDirectNotes(folderPath: string, notes: NoteItem[]): NoteItem[] {
    const prefix = folderPath + "/"
    return notes.filter(n => {
        if (!n.path || !n.path.startsWith(prefix)) return false
        return !n.path.slice(prefix.length).includes("/")
    })
}

interface FolderNodeItemProps {
    node: FolderNode
    notes: NoteItem[]
    openFolders: Set<string>
    selectedPath: string | null
    depth: number
    onToggle: (path: string) => void
    onSelectNote: (path: string) => void
    onCreateNoteInFolder: (path: string) => void
}

function FolderNodeItem({ node, notes, openFolders, selectedPath, depth, onToggle, onSelectNote, onCreateNoteInFolder }: FolderNodeItemProps) {
    const isOpen = openFolders.has(node.path)
    const directNotes = getDirectNotes(node.path, notes)

    return (
        <>
            <TreeItem
                name={node.name}
                isFolder
                isOpen={isOpen}
                depth={depth}
                onClick={() => onToggle(node.path)}
                onAddInFolder={() => onCreateNoteInFolder(node.path)}
            />
            {isOpen && (
                <>
                    {node.children.map(child => (
                        <FolderNodeItem
                            key={child.path}
                            node={child}
                            notes={notes}
                            openFolders={openFolders}
                            selectedPath={selectedPath}
                            depth={depth + 1}
                            onToggle={onToggle}
                            onSelectNote={onSelectNote}
                            onCreateNoteInFolder={onCreateNoteInFolder}
                        />
                    ))}
                    {directNotes.map(note => (
                        <TreeItem
                            key={note.path}
                            name={getNoteName(note)}
                            depth={depth + 1}
                            active={selectedPath === note.path}
                            onClick={() => note.path && onSelectNote(note.path)}
                        />
                    ))}
                </>
            )}
        </>
    )
}

interface FolderSectionProps {
    folders: string[]
    notes: NoteItem[]
    selectedPath: string | null
    vaultId: string
    onSelectNote: (vaultId: string, path: string) => void
    onCreateNote: (folderPath?: string) => void
}

function FolderSection({ folders, notes, selectedPath, vaultId, onSelectNote, onCreateNote }: FolderSectionProps) {
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

    const tree = buildFolderTree(folders)
    const rootNotes = notes.filter(n => n.path && !n.path.includes("/"))

    return (
        <>
            <div className={cls.SectionHeader}>
                <span className={cls.SectionLabel}>All Notes</span>
                <Button variant="ghost" className={cls.CreateNoteBtn} onClick={() => onCreateNote()} title="New note">
                    <PlusIcon/>
                </Button>
            </div>
            {tree.map(node => (
                <FolderNodeItem
                    key={node.path}
                    node={node}
                    notes={notes}
                    openFolders={openFolders}
                    selectedPath={selectedPath}
                    depth={0}
                    onToggle={toggleFolder}
                    onSelectNote={path => onSelectNote(vaultId, path)}
                    onCreateNoteInFolder={onCreateNote}
                />
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
    const [searchQuery, setSearchQuery] = useState("")

    function handleVaultChange(e: React.ChangeEvent<HTMLSelectElement>) {
        void selectVault(e.target.value)
    }

    function handleSelectNote(vid: string, path: string) {
        void selectNote(vid, path)
    }

    function handleCreateNote(folderPath?: string) {
        OpenDialog(
            <CreateNoteDialog
                initialPath={folderPath ? folderPath + "/" : ""}
                folders={folders}
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
                    <SearchIcon className={cls.SearchIcon}/>
                    <input
                        className={cls.SearchInput}
                        placeholder="Search notes…"
                        value={searchQuery}
                        onChange={e => setSearchQuery(e.target.value)}
                    />
                </div>
            </div>
            <div className={cls.ScrollArea}>
                {!vaultId && (
                    <div className={cls.EmptyVaultHint}>Select a vault to browse notes</div>
                )}
                {vaultId && searchQuery.trim() && notes
                    .filter(n => getNoteName(n).toLowerCase().includes(searchQuery.trim().toLowerCase()))
                    .map(note => {
                        const dir = note.path ? note.path.split("/").slice(0, -1).join("/") : undefined
                        return (
                            <TreeItem
                                key={note.path}
                                name={getNoteName(note)}
                                subtitle={dir || undefined}
                                active={selectedPath === note.path}
                                onClick={() => note.path && handleSelectNote(vaultId, note.path)}
                            />
                        )
                    })
                }
                {vaultId && !searchQuery.trim() && (
                    <FolderSection
                        folders={folders}
                        notes={notes}
                        selectedPath={selectedPath}
                        vaultId={vaultId}
                        onSelectNote={handleSelectNote}
                        onCreateNote={folderPath => handleCreateNote(folderPath)}
                    />
                )}
            </div>
        </div>
    )
}
