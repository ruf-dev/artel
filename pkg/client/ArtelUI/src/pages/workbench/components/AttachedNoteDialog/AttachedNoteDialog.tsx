import {Button, Loader} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import {useNoteContent} from "@/app/hooks/Notes"
import {useSanitizedMarkdownHtml} from "@/app/hooks/useSanitizedMarkdownHtml"
import {getItemName} from "@/components/FileTree/fileTree.ts"
import CloseIcon from "@/icons/common/CloseIcon.tsx"
import cls from "@/pages/workbench/components/AttachedNoteDialog/AttachedNoteDialog.module.css"

interface Props {
    vaultId: string
    path: string
}

// Dialog for viewing an attached note file. Opened via OpenDialog() from
// UserMessageBubble when a user clicks an attachment chip. Fetches and displays
// the note content as sanitized markdown HTML.
export default function AttachedNoteDialog({vaultId, path}: Props) {
    const {CloseDialog} = useDialog()
    const {content, isLoading} = useNoteContent(vaultId, path)
    const html = useSanitizedMarkdownHtml(content)

    return (
        <div className={cls.AttachedNoteDialogContainer} role="dialog" aria-modal="true">
            <div className={cls.Header}>
                <h2 className={cls.Title}>{getItemName({path})}</h2>
                <Button
                    variant="secondary"
                    className={cls.CloseButton}
                    onClick={CloseDialog}
                    aria-label="Close note"
                    title="Close note"
                >
                    <CloseIcon className={cls.CloseIcon}/>
                </Button>
            </div>
            <div className={cls.Body}>
                {isLoading ? (
                    <div className={cls.LoaderWrapper}><Loader/></div>
                ) : (
                    <div className={cls.Content} dangerouslySetInnerHTML={{__html: html ?? ""}}/>
                )}
            </div>
        </div>
    )
}
