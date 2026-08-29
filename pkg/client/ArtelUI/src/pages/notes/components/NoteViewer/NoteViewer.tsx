import {useSanitizedMarkdownHtml} from "@/app/hooks/useSanitizedMarkdownHtml"
import cls from "@/pages/notes/components/NoteViewer/NoteViewer.module.css"
import NoteContent from "@/pages/notes/components/NoteViewer/components/NoteContent/NoteContent.tsx"

interface NoteViewerProps {
    content: string | null
    fontScale?: number
    onContentClick?: () => void
}

export default function NoteViewer({ content, fontScale, onContentClick }: NoteViewerProps) {
    const sanitizedHtml = useSanitizedMarkdownHtml(content)

    return (
        <div className={cls.NoteViewerContainer} style={{ zoom: fontScale ?? 1 }}>
            {content === null ? (
                <div className={cls.EmptyState}>Select a note</div>
            ) : (
                <NoteContent rawHtml={sanitizedHtml!} onContentClick={onContentClick} />
            )}
        </div>
    )
}
