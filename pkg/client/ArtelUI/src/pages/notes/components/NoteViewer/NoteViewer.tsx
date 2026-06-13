import { useMemo } from "react"
import { marked } from "marked"
import DOMPurify from "dompurify"

import cls from "./NoteViewer.module.css"
import WikiChip from "@/pages/notes/components/WikiChip/WikiChip.tsx"
import { useNotes } from "@/app/hooks/Notes.ts"
import type { NoteItem } from "@/app/hooks/Notes.ts"

interface NoteViewerProps {
    content: string | null
}

interface ContentSegment {
    type: "html" | "wiki"
    value: string
}

function parseWikiLinks(html: string): ContentSegment[] {
    const parts: ContentSegment[] = []
    const regex = /\[\[([^\]]+)\]\]/g
    let lastIndex = 0
    let match: RegExpExecArray | null

    while ((match = regex.exec(html)) !== null) {
        if (match.index > lastIndex) {
            parts.push({ type: "html", value: html.slice(lastIndex, match.index) })
        }
        parts.push({ type: "wiki", value: match[1] })
        lastIndex = match.index + match[0].length
    }

    if (lastIndex < html.length) {
        parts.push({ type: "html", value: html.slice(lastIndex) })
    }

    return parts
}

function resolveNote(notes: NoteItem[], name: string): NoteItem | undefined {
    return notes.find(n =>
        n.path === name ||
        n.path === name + ".md" ||
        n.path?.endsWith("/" + name) ||
        n.path?.endsWith("/" + name + ".md")
    )
}

function NoteContent({ rawHtml }: { rawHtml: string }) {
    const segments = useMemo(() => parseWikiLinks(rawHtml), [rawHtml])
    const { vaultId, notes, selectNote } = useNotes()

    return (
        <div className={cls.NoteBody}>
            {segments.map((seg, i) => {
                if (seg.type === "wiki") {
                    return (
                        <WikiChip
                            key={i}
                            name={seg.value}
                            onClick={() => {
                                const note = resolveNote(notes, seg.value)
                                if (note?.path && vaultId) void selectNote(vaultId, note.path)
                            }}
                        />
                    )
                }
                return (
                    <span
                        key={i}
                        dangerouslySetInnerHTML={{ __html: seg.value }}
                    />
                )
            })}
        </div>
    )
}

export default function NoteViewer({ content }: NoteViewerProps) {
    const sanitizedHtml = useMemo(() => {
        if (!content) return null
        const rawHtml = marked.parse(content) as string
        return DOMPurify.sanitize(rawHtml)
    }, [content])

    return (
        <div className={cls.NoteViewerContainer}>
            {sanitizedHtml === null ? (
                <div className={cls.EmptyState}>Select a note</div>
            ) : (
                <NoteContent rawHtml={sanitizedHtml} />
            )}
        </div>
    )
}
