import { useParams } from "react-router-dom"
import { LoadingWrapper } from "@vervstack/chures"

import { usePublicDocs } from "@/pages/docs/hooks/usePublicDocs.ts"
import DocsSidebar from "@/pages/docs/components/DocsSidebar/DocsSidebar.tsx"
import DocsNoteViewer from "@/pages/docs/components/DocsNoteViewer/DocsNoteViewer.tsx"
import cls from "@/pages/docs/DocsPage.module.css"

// Chromeless, fully read-only public reader — no NoteEditor, no ModeBar, no vault switcher,
// no login prompt, no auth. Anonymous data comes from usePublicDocs (PublicDocsAPI).
export default function DocsPage() {
    const { slug } = useParams<{ slug: string }>()
    const docs = usePublicDocs(slug ?? "")

    return (
        <div className={cls.DocsPageContainer}>
            <LoadingWrapper isLoading={docs.isLoading}>
                {docs.error ? (
                    <div className={cls.ErrorState}>This vault could not be found.</div>
                ) : (
                    <div className={cls.Body}>
                        <DocsSidebar
                            vaultName={docs.vaultName}
                            folders={docs.folders}
                            notes={docs.notes}
                            selectedPath={docs.selectedPath}
                            onSelectNote={docs.selectNote}
                        />
                        <DocsNoteViewer content={docs.noteContent} />
                    </div>
                )}
            </LoadingWrapper>
        </div>
    )
}
