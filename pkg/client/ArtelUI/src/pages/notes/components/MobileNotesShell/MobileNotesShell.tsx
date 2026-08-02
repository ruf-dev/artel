import { useEffect, useRef, useState } from "react"

import { NoteMode, useNotes } from "@/app/hooks/Notes.ts"
import { SaveStatus } from "@/pages/notes/components/SaveStatusIndicator/SaveStatusIndicator.tsx"
import ModeBar from "@/pages/notes/components/ModeBar/ModeBar.tsx"
import NoteEditor from "@/pages/notes/components/NoteEditor/NoteEditor.tsx"
import NoteViewer from "@/pages/notes/components/NoteViewer/NoteViewer.tsx"
import MobileTopBar from "@/pages/notes/components/MobileNotesShell/components/MobileTopBar/MobileTopBar.tsx"
import MobileDrawer from "@/pages/notes/components/MobileNotesShell/components/MobileDrawer/MobileDrawer.tsx"
import cls from "@/pages/notes/components/MobileNotesShell/MobileNotesShell.module.css"

interface VaultOption {
    id: string
    name: string
    isPublic?: boolean
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
    hideEdit?: boolean
}

export default function MobileNotesShell(props: MobileNotesShellProps) {
    const { vaultOptions, noteContent, selectedPath, mode, showEditor } = props
    const { scrollTopRef, fontScale, onModeChange, onChange, onContentClick } = props
    const { onEscape, hideEdit } = props

    const [sidebarOpen, setSidebarOpen] = useState(false)
    const { selectedPath: storeSelectedPath, highlightNote } = useNotes()
    const prevSelectedPath = useRef(storeSelectedPath)

    useEffect(() => {
        if (storeSelectedPath && storeSelectedPath !== prevSelectedPath.current) {
            setSidebarOpen(false)
        }
        prevSelectedPath.current = storeSelectedPath
    }, [storeSelectedPath])

    function handleFindCurrentFile() {
        if (!storeSelectedPath) return
        setSidebarOpen(true)
        highlightNote(storeSelectedPath)
    }

    const contentAreaClass = `${cls.ContentArea}${mode === "read" ? ` ${cls.ContentAreaRead}` : ""}`

    return (
        <div className={cls.MobileShellContainer}>
            <div className={cls.StatusBarSpacer} />
            <MobileTopBar
                selectedPath={selectedPath}
                noteContent={noteContent}
                sidebarOpen={sidebarOpen}
                onHamburgerClick={() => setSidebarOpen(v => !v)}
                onFindCurrentFile={handleFindCurrentFile}
            />
            <ModeBar active={mode} onModeChange={onModeChange} hideEdit={hideEdit} />
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
