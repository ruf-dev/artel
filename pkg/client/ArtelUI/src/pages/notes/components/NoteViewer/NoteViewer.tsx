import { useMemo } from "react"
import { marked } from "marked"
import DOMPurify from "dompurify"

import cls from "./NoteViewer.module.css"
import WikiChip from "@/pages/notes/components/WikiChip/WikiChip.tsx"

interface NoteViewerProps {
    content: string | null
    fontScale?: number
    onContentClick?: () => void
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

function NoteContent({ rawHtml, onContentClick }: { rawHtml: string; onContentClick?: () => void }) {
    const segments = useMemo(() => parseWikiLinks(rawHtml), [rawHtml])

    return (
        <div className={cls.NoteBody} onClick={onContentClick}>
            {segments.map((seg, i) => {
                if (seg.type === "wiki") {
                    return <WikiChip key={i} name={seg.value} />
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

export default function NoteViewer({ content, fontScale, onContentClick }: NoteViewerProps) {
    const sanitizedHtml = useMemo(() => {
        if (!content) return null
        const rawHtml = marked.parse(content) as string
        return DOMPurify.sanitize(rawHtml)
    }, [content])

    return (
        <div className={cls.NoteViewerContainer} style={{ zoom: fontScale ?? 1 }}>
            {sanitizedHtml === null ? (
                <div className={cls.EmptyState}>Select a note</div>
            ) : (
                <NoteContent rawHtml={sanitizedHtml} onContentClick={onContentClick} />
            )}
        </div>
    )
}
