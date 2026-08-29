import {Button} from "@vervstack/chures"
import {MorphIcon} from "morphicons/react"
import {Circle, Send} from "lucide"

import {cn} from "@/app/utils/cn.ts"
import PlusIcon from "@/icons/common/PlusIcon.tsx"
import cls from "@/pages/workbench/components/Chat/components/ComposerRow/ComposerRow.module.css"

interface Props {
    onNewChat: () => void
    onSend: () => void
    disabled: boolean
    sendDisabled: boolean
    hasText: boolean
    hideNewChatButton?: boolean
}

// Send starts life as a bare circle glyph — nothing to send yet — and morphs into
// the paper-plane send glyph the moment the composer holds text.
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
            <Button
                variant="primary"
                className={cls.SendButton}
                onClick={props.onSend}
                disabled={props.sendDisabled}
                aria-label="Send message"
            >
                <MorphIcon
                    icon={props.hasText ? Send : Circle}
                    size={15}
                    strokeWidth={2}
                    className={cls.SendIcon}
                />
            </Button>
        </div>
    )
}
