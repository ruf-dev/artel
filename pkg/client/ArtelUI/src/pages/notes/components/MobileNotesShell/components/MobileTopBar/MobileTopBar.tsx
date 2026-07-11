import {Button} from "@vervstack/chures"

import cls from "@/pages/notes/components/MobileNotesShell/components/MobileTopBar/MobileTopBar.module.css"

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

export default function MobileTopBar({ selectedPath, noteContent, sidebarOpen, onHamburgerClick }: MobileTopBarProps) {
    const title = getNoteTitle(selectedPath)
    const meta = getNoteMeta(noteContent)
    const hamburgerClass = `${cls.HamburgerBtn}${sidebarOpen ? ` ${cls.HamburgerBtnActive}` : ""}`

    return (
        <div className={cls.MobileTopBarContainer}>
            <Button variant="ghost" className={hamburgerClass} onClick={onHamburgerClick} aria-label="Open sidebar">
                <svg
                    viewBox="0 0 18 18" width={16} height={16} fill="none" stroke="currentColor"
                    strokeWidth="1.6" strokeLinecap="round"
                >
                    <path d="M3 4.5h12M3 9h8M3 13.5h10" />
                </svg>
            </Button>
            <div className={cls.NoteInfoColumn}>
                <div className={cls.NoteTitle}>{title}</div>
                {meta && <div className={cls.NoteMeta}>{meta}</div>}
            </div>
            <Button variant="ghost" className={cls.ActionBtn} aria-label="Share">
                <svg
                    viewBox="0 0 18 18" width={16} height={16} fill="none" stroke="currentColor"
                    strokeWidth="1.5" strokeLinecap="round"
                >
                    <circle cx="14" cy="3.5" r="1.8" />
                    <circle cx="4" cy="9" r="1.8" />
                    <circle cx="14" cy="14.5" r="1.8" />
                    <path d="M5.7 10l6.5 3.3M12.2 4.7L5.7 8" />
                </svg>
            </Button>
            <Button variant="ghost" className={cls.ActionBtn} aria-label="More options">
                <svg viewBox="0 0 18 18" width={16} height={16} fill="currentColor">
                    <circle cx="9" cy="4" r="1.4" />
                    <circle cx="9" cy="9" r="1.4" />
                    <circle cx="9" cy="14" r="1.4" />
                </svg>
            </Button>
        </div>
    )
}
