import {useSanitizedMarkdownHtml} from "@/app/hooks/useSanitizedMarkdownHtml"
import cls from "@/pages/docs/components/DocsNoteViewer/DocsNoteViewer.module.css"

interface DocsNoteViewerProps {
    content: string | null
}

// Simplified local counterpart to pages/notes/components/NoteViewer/NoteViewer.tsx (which is
// local to NotesPage per the colocation rule). Skips the [[wiki-link]] -> chip parsing that
// NoteContent does — there's no cross-note navigation concept defined for /docs.
export default function DocsNoteViewer({ content }: DocsNoteViewerProps) {
    const sanitizedHtml = useSanitizedMarkdownHtml(content)

    return (
        <div className={cls.DocsNoteViewerContainer}>
            {content === null ? (
                <div className={cls.EmptyState}>Select a note</div>
            ) : (
                <div className={cls.Content} dangerouslySetInnerHTML={{ __html: sanitizedHtml! }} />
            )}
        </div>
    )
}
