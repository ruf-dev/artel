import { useEffect, useRef, useState } from "react"
import cls from "./MobileNotesShell.module.css"
import { NoteMode, useNotes } from "@/app/hooks/Notes.ts"
import { useDialog } from "@/app/hooks/Dialog.ts"
import { useBakeError } from "@/app/hooks/useErrorToast.ts"
import { SaveStatus } from "@/pages/notes/components/SaveStatusIndicator/SaveStatusIndicator.tsx"
import ModeBar from "@/pages/notes/components/ModeBar/ModeBar.tsx"
import NoteEditor from "@/pages/notes/components/NoteEditor/NoteEditor.tsx"
import NoteViewer from "@/pages/notes/components/NoteViewer/NoteViewer.tsx"
import NotesSidebar from "@/pages/notes/components/NotesSidebar/NotesSidebar.tsx"
import CreateNoteDialog from "@/pages/notes/components/CreateNoteDialog/CreateNoteDialog.tsx"
import {Button} from "@vervstack/chures"
import ArtelLogoIcon from "@/pages/notes/components/icons/ArtelLogoIcon.tsx"
import CloseIcon from "@/pages/notes/components/icons/CloseIcon.tsx"

interface VaultOption {
    id: string
    name: string
}

export interface MobileNotesShellProps {
    vaultOptions: VaultOption[]
    noteContent: string | null
    selectedPath: string | null
    mode: NoteMode
    showEditor: boolean
    saveStatus: SaveStatus
    saveError?: string
    scrollTopRef: React.MutableRefObject<number>
    fontScale: number
    onModeChange: (mode: NoteMode) => void
    onChange: (content: string) => void
    onContentClick: (() => void) | undefined
    onEscape: (() => void) | undefined
    onRename: () => void
}

function getNoteTitle(selectedPath: string | null): string {
    if (!selectedPath) return "No note selected"
    const parts = selectedPath.split("/")
    const filename = parts[parts.length - 1] || selectedPath
    return filename.endsWith(".md") ? filename.slice(0, -3) : filename
}

function getNoteMeta(content: string | null): string {
    if (!content) return ""
    const words = content.trim() ? content.trim().split(/\s+/).length : 0
    const links = (content.match(/\[\[.+?\]\]/g) || []).length
    if (!words && !links) return ""
    const parts: string[] = []
    if (words) parts.push(`${words} words`)
    if (links) parts.push(`${links} links`)
    return parts.join(" · ")
}

interface MobileTopBarProps {
    selectedPath: string | null
    noteContent: string | null
    sidebarOpen: boolean
    onHamburgerClick: () => void
}

function MobileTopBar({ selectedPath, noteContent, sidebarOpen, onHamburgerClick }: MobileTopBarProps) {
    const title = getNoteTitle(selectedPath)
    const meta = getNoteMeta(noteContent)
    const hamburgerClass = `${cls.HamburgerBtn}${sidebarOpen ? ` ${cls.HamburgerBtnActive}` : ""}`

    return (
        <div className={cls.TopBarContainer}>
            <button type="button" className={hamburgerClass} onClick={onHamburgerClick} aria-label="Open sidebar">
                <svg viewBox="0 0 18 18" width={16} height={16} fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round">
                    <path d="M3 4.5h12M3 9h8M3 13.5h10" />
                </svg>
            </button>
            <div className={cls.NoteInfoColumn}>
                <div className={cls.NoteTitle}>{title}</div>
                {meta && <div className={cls.NoteMeta}>{meta}</div>}
            </div>
            <button type="button" className={cls.ActionBtn} aria-label="Share">
                <svg viewBox="0 0 18 18" width={16} height={16} fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round">
                    <circle cx="14" cy="3.5" r="1.8" />
                    <circle cx="4" cy="9" r="1.8" />
                    <circle cx="14" cy="14.5" r="1.8" />
                    <path d="M5.7 10l6.5 3.3M12.2 4.7L5.7 8" />
                </svg>
            </button>
            <button type="button" className={cls.ActionBtn} aria-label="More options">
                <svg viewBox="0 0 18 18" width={16} height={16} fill="currentColor">
                    <circle cx="9" cy="4" r="1.4" />
                    <circle cx="9" cy="9" r="1.4" />
                    <circle cx="9" cy="14" r="1.4" />
                </svg>
            </button>
        </div>
    )
}

function DrawerCloseButton({onClose}: {onClose: () => void}) {
    return (
        <Button variant="ghost" className={cls.DrawerCloseBtn} onClick={onClose} aria-label="Close sidebar">
            <CloseIcon/>
        </Button>
    )
}

interface MobileDrawerProps {
    open: boolean
    onClose: () => void
    vaultOptions: VaultOption[]
}

function MobileDrawer({ open, onClose, vaultOptions }: MobileDrawerProps) {
    const { vaultId, folders, createNote } = useNotes()
    const { OpenDialog } = useDialog()
    const bakeError = useBakeError()

    function handleNewNote() {
        OpenDialog(
            <CreateNoteDialog
                initialPath=""
                folders={folders}
                onConfirm={(path: string) =>
                    createNote(path).catch(err => bakeError("Failed to create note", err))
                }
            />
        )
        onClose()
    }

    const backdropClass = `${cls.Backdrop}${open ? ` ${cls.BackdropOpen}` : ""}`
    const panelClass = `${cls.DrawerPanel}${open ? ` ${cls.DrawerPanelOpen}` : ""}`

    return (
        <>
            <div className={backdropClass} onClick={onClose} />
            <div className={panelClass}>
                <div className={cls.DrawerStatusSpacer} />
                <div className={cls.DrawerHeader}>
                    <ArtelLogoIcon/>
                    <span className={cls.DrawerWordmark}>artel</span>
                    <span className={cls.DrawerNotesLabel}>notes</span>
                    <span className={cls.DrawerSpacer} />
                    <DrawerCloseButton onClose={onClose}/>
                </div>
                <div className={cls.DrawerBody}>
                    <NotesSidebar vaults={vaultOptions} />
                </div>
                {vaultId && (
                    <div className={cls.DrawerFooter}>
                        <button type="button" className={cls.NewNoteBtn} onClick={handleNewNote}>
                            <svg viewBox="0 0 12 12" width={11} height={11} fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round">
                                <path d="M6 1v10M1 6h10" />
                            </svg>
                            New Note
                        </button>
                    </div>
                )}
            </div>
        </>
    )
}

export default function MobileNotesShell({
    vaultOptions,
    noteContent,
    selectedPath,
    mode,
    showEditor,
    scrollTopRef,
    fontScale,
    onModeChange,
    onChange,
    onContentClick,
    onEscape,
}: MobileNotesShellProps) {
    const [sidebarOpen, setSidebarOpen] = useState(false)
    const { selectedPath: storeSelectedPath } = useNotes()
    const prevSelectedPath = useRef(storeSelectedPath)

    useEffect(() => {
        if (storeSelectedPath && storeSelectedPath !== prevSelectedPath.current) {
            setSidebarOpen(false)
        }
        prevSelectedPath.current = storeSelectedPath
    }, [storeSelectedPath])

    const contentAreaClass = `${cls.ContentArea}${mode === "read" ? ` ${cls.ContentAreaRead}` : ""}`

    return (
        <div className={cls.MobileShellContainer}>
            <div className={cls.StatusBarSpacer} />
            <MobileTopBar
                selectedPath={selectedPath}
                noteContent={noteContent}
                sidebarOpen={sidebarOpen}
                onHamburgerClick={() => setSidebarOpen(v => !v)}
            />
            <ModeBar active={mode} onModeChange={onModeChange} />
            <div className={contentAreaClass}>
                {showEditor ? (
                    <NoteEditor
                        content={noteContent ?? ""}
                        onChange={onChange}
                        scrollTopRef={scrollTopRef}
                        fontScale={fontScale}
                        onEscape={onEscape}
                    />
                ) : (
                    <NoteViewer
                        content={noteContent}
                        fontScale={fontScale}
                        onContentClick={onContentClick}
                    />
                )}
            </div>
            <MobileDrawer
                open={sidebarOpen}
                onClose={() => setSidebarOpen(false)}
                vaultOptions={vaultOptions}
            />
        </div>
    )
}
