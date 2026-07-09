import { useMemo } from "react"

import WikiChip from "@/pages/notes/components/WikiChip/WikiChip.tsx"
import cls from "@/pages/notes/components/NoteViewer/components/NoteContent/NoteContent.module.css"

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

interface NoteContentProps {
    rawHtml: string
    onContentClick?: () => void
}

export default function NoteContent({ rawHtml, onContentClick }: NoteContentProps) {
    const segments = useMemo(() => parseWikiLinks(rawHtml), [rawHtml])

    return (
        <div className={cls.NoteContentContainer} onClick={onContentClick}>
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
