import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import PlusIcon from "@/icons/common/PlusIcon.tsx"
import SendIcon from "@/pages/workbench/components/Chat/components/icons/SendIcon.tsx"
import cls from "@/pages/workbench/components/Chat/components/ComposerRow/ComposerRow.module.css"

interface Props {
    onNewChat: () => void
    onSend: () => void
    disabled: boolean
    sendDisabled: boolean
    hideNewChatButton?: boolean
}

// The bottom control strip of the composer card: new-chat button, spacer, send.
export default function ComposerRow(props: Props) {
    return (
        <div className={cls.ComposerRowContainer}>
            <Button
                variant="secondary"
                className={cn(cls.NewChatButton, props.hideNewChatButton && cls.NewChatButtonHidden)}
                onClick={props.onNewChat}
                disabled={props.disabled}
                aria-label="New chat"
                title="New chat"
            >
                <PlusIcon/>
            </Button>
            <div className={cls.Spacer}/>
            <Button
                variant="primary"
                className={cls.SendButton}
                onClick={props.onSend}
                disabled={props.sendDisabled}
                aria-label="Send message"
            >
                <SendIcon/>
            </Button>
        </div>
    )
}
