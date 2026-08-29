import {motion} from "framer-motion"
import {Button} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import {stripAttachmentsPreamble} from "@/pages/workbench/processes/workbenchAttachments"
import AttachmentChip from "@/pages/workbench/components/Chat/components/AttachmentChip/AttachmentChip.tsx"
import AttachedNoteDialog from "@/pages/workbench/components/AttachedNoteDialog/AttachedNoteDialog.tsx"
import cls from "@/pages/workbench/components/Chat/components/UserMessageBubble/UserMessageBubble.module.css"
import RetryIcon from "@/pages/workbench/components/Chat/components/icons/RetryIcon.tsx"

interface Props {
    text: string
    attachments?: string[]
    vaultId?: string
    showRetry?: boolean
    onRetry?: () => void
    retryDisabled?: boolean
}

export default function UserMessageBubble({text, attachments, vaultId, showRetry, onRetry, retryDisabled}: Props) {
    const {OpenDialog} = useDialog()

    const caption = stripAttachmentsPreamble(text, attachments ?? [])

    function handleChipClick(path: string): void {
        if (vaultId) {
            OpenDialog(<AttachedNoteDialog vaultId={vaultId} path={path}/>)
        }
    }

    return (
        <motion.div
            className={cls.UserMessageBubbleContainer}
            initial={{opacity: 0, y: 28, scale: 0.9}}
            animate={{opacity: 1, y: 0, scale: 1}}
            exit={{opacity: 0, scale: 0.9}}
            transition={{type: "spring", stiffness: 420, damping: 32}}
        >
            {attachments && attachments.length > 0 && (
                <div className={cls.AttachmentRow}>
                    {attachments.map(path => (
                        <AttachmentChip
                            key={path}
                            path={path}
                            onClick={handleChipClick}
                        />
                    ))}
                </div>
            )}
            <p className={cls.Text}>{caption}</p>
            {showRetry && (
                <div className={cls.RetryRow}>
                    <Button
                        variant="unstyled"
                        className={cls.ActBtn}
                        aria-label="Retry message"
                        aria-disabled={retryDisabled || undefined}
                        data-tooltip-id="root-tooltip"
                        data-tooltip-content={retryDisabled ? "Nothing to retry" : "Retry"}
                        onClick={() => {
                            if (!retryDisabled) onRetry?.()
                        }}
                    >
                        <RetryIcon/>
                    </Button>
                </div>
            )}
        </motion.div>
    )
}
