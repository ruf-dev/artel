import {Button} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import CloseIcon from "@/icons/common/CloseIcon.tsx"
import ChatMessageList from "@/pages/workbench/components/Chat/components/ChatMessageList/ChatMessageList.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import cls from "@/pages/workbench/components/WorkbenchHistoryTranscriptDialog/WorkbenchHistoryTranscriptDialog.module.css"

interface Props {
    title: string
    items: ChatItem[]
    vaultId: string
}

// Read-only transcript view for a past Docker workbench session, opened via
// OpenDialog() from useWorkbenchHistory when a docker-mode history row is
// selected. Docker sessions are closed past transcripts (unlike resumable api
// chats), so there's nothing to interact with — permission decisions are no-ops.
export default function WorkbenchHistoryTranscriptDialog({title, items, vaultId}: Props) {
    const {CloseDialog} = useDialog()

    return (
        <div className={cls.WorkbenchHistoryTranscriptDialogContainer} role="dialog" aria-modal="true">
            <div className={cls.Header}>
                <h2 className={cls.Title}>{title}</h2>
                <Button
                    variant="secondary"
                    className={cls.CloseButton}
                    onClick={CloseDialog}
                    aria-label="Close transcript"
                    title="Close transcript"
                >
                    <CloseIcon className={cls.CloseIcon}/>
                </Button>
            </div>
            <div className={cls.Body}>
                <ChatMessageList items={items} vaultId={vaultId} onPermissionDecision={() => {}}/>
            </div>
        </div>
    )
}
