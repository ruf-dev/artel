import { useMemo, useState } from "react"
import { Button, Input } from "@vervstack/chures"

import { useDialog, useDialogKeyboard } from "@/app/hooks/Dialog"
import cls from "@/pages/notes/components/CreateNoteDialog/CreateNoteDialog.module.css"
import SuggestionList from "@/pages/notes/components/CreateNoteDialog/components/SuggestionList/SuggestionList.tsx"

function ensureExtension(path: string): string {
    const filename = path.split("/").pop() ?? path
    return filename.includes(".") ? path : path + ".md"
}

interface Props {
    onConfirm: (path: string) => Promise<void>
    initialPath?: string
    folders?: string[]
}

export default function CreateNoteDialog({ onConfirm, initialPath, folders }: Props) {
    const { CloseDialog } = useDialog()
    const [path, setPath] = useState(initialPath ?? "")
    const [loading, setLoading] = useState(false)
    const [highlightedIdx, setHighlightedIdx] = useState(-1)
    const [inputFocused, setInputFocused] = useState(false)

    const suggestions = useMemo(() => {
        if (!folders?.length || !path) return []
        const lastSlash = path.lastIndexOf("/")
        const prefix = lastSlash === -1 ? "" : path.slice(0, lastSlash + 1)
        const fragment = lastSlash === -1 ? path : path.slice(lastSlash + 1)
        return folders
            .filter(f => {
                if (prefix && !f.startsWith(prefix)) return false
                const remaining = f.slice(prefix.length)
                if (remaining.includes("/")) return false
                return remaining.toLowerCase().includes(fragment.toLowerCase())
            })
            .slice(0, 8)
    }, [folders, path])

    const showSuggestions = inputFocused && suggestions.length > 0

    function handleSuggestionSelect(folder: string) {
        setPath(folder + "/")
        setHighlightedIdx(-1)
    }

    function handleConfirm() {
        if (!path.trim()) return
        setLoading(true)
        onConfirm(ensureExtension(path.trim())).finally(() => {
            setLoading(false)
            CloseDialog()
        })
    }

    const dialogKeyHandler = useDialogKeyboard(handleConfirm)

    function handleKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
        if (showSuggestions) {
            if (e.key === "ArrowDown") {
                e.preventDefault()
                setHighlightedIdx(i => Math.min(i + 1, suggestions.length - 1))
                return
            }
            if (e.key === "ArrowUp") {
                e.preventDefault()
                setHighlightedIdx(i => Math.max(i - 1, -1))
                return
            }
            if (e.key === "Escape") {
                setInputFocused(false)
                return
            }
            if (e.key === "Enter" && highlightedIdx >= 0) {
                e.preventDefault()
                handleSuggestionSelect(suggestions[highlightedIdx])
                return
            }
        }
        dialogKeyHandler(e)
    }

    return (
        <div className={cls.CreateNoteContainer} role="dialog" aria-modal="true">
            <h2 className={cls.CreateNoteTitle}>New Note</h2>
            <p className={cls.CreateNoteLabel}>Path</p>
            <div className={cls.PathWrapper}>
                <Input
                    className={cls.CreateNoteInput}
                    type="text"
                    placeholder="folder/note-name"
                    value={path}
                    setValue={v => { setPath(v); setHighlightedIdx(-1) }}
                    onFocus={() => setInputFocused(true)}
                    onBlur={() => setInputFocused(false)}
                    onKeyDown={handleKeyDown}
                    disabled={loading}
                    autoFocus
                    spellCheck={false}
                />
                {showSuggestions && (
                    <SuggestionList
                        suggestions={suggestions}
                        highlightedIdx={highlightedIdx}
                        onSelect={handleSuggestionSelect}
                    />
                )}
            </div>
            <div className={cls.CreateNoteActions}>
                <Button variant="ghost" onClick={CloseDialog} disabled={loading}>
                    Cancel
                </Button>
                <Button variant="primary" onClick={handleConfirm} disabled={loading || !path.trim()}>
                    {loading ? "…" : "Create"}
                </Button>
            </div>
        </div>
    )
}
