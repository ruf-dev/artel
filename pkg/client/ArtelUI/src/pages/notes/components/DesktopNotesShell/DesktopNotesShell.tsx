import {Button} from "@vervstack/chures"

import { NoteMode } from "@/app/hooks/Notes.ts"
import { SaveStatus } from "@/pages/notes/components/SaveStatusIndicator/SaveStatusIndicator.tsx"
import NotesSidebar from "@/pages/notes/components/NotesSidebar/NotesSidebar.tsx"
import NoteViewer from "@/pages/notes/components/NoteViewer/NoteViewer.tsx"
import NoteEditor from "@/pages/notes/components/NoteEditor/NoteEditor.tsx"
import BreadcrumbBar from "@/pages/notes/components/BreadcrumbBar/BreadcrumbBar.tsx"
import cls from "@/pages/notes/components/DesktopNotesShell/DesktopNotesShell.module.css"

interface VaultOption {
    id: string
    name: string
    isPublic?: boolean
}

export interface DesktopNotesShellProps {
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
    onZoomIn: () => void
    onZoomOut: () => void
    hideEdit?: boolean
}

export default function DesktopNotesShell(props: DesktopNotesShellProps) {
    const { vaultOptions, noteContent, selectedPath, mode, showEditor } = props
    const { saveStatus, saveError, scrollTopRef, fontScale } = props
    const { onModeChange, onChange, onContentClick, onEscape, onRename, hideEdit } = props
    const { onZoomIn, onZoomOut } = props

    return (
        <div className={cls.DesktopNotesShellContainer}>
            <div className={cls.SidebarWrapper}>
                <NotesSidebar vaults={vaultOptions}/>
            </div>
            <div className={cls.ViewerWrapper}>
                {selectedPath && (
                    <BreadcrumbBar
                        path={selectedPath}
                        mode={mode}
                        onModeChange={onModeChange}
                        saveStatus={showEditor ? saveStatus : 'idle'}
                        saveError={saveError}
                        onRename={onRename}
                        hideEdit={hideEdit}
                    />
                )}
                {showEditor ? (
                    <NoteEditor
                        content={noteContent ?? ''}
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
                <div className={cls.ZoomControls}>
                    <Button variant="ghost" className={cls.ZoomBtn} onClick={onZoomOut}>−</Button>
                    <span className={cls.ZoomLabel}>{Math.round(fontScale * 100)}%</span>
                    <Button variant="ghost" className={cls.ZoomBtn} onClick={onZoomIn}>+</Button>
                </div>
            </div>
        </div>
    )
}
