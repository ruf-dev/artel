import {Button} from "@vervstack/chures"

import cls from "@/pages/workbench/components/Chat/components/ChatHeader/ChatHeader.module.css"

interface Props {
    onNewChat: () => void
    disabled?: boolean
}

export default function ChatHeader({onNewChat, disabled}: Props) {
    return (
        <div className={cls.ChatHeaderContainer}>
            <Button variant="secondary" className={cls.NewChatButton} onClick={onNewChat} disabled={disabled}>
                <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
                     strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                New Chat
            </Button>
        </div>
    )
}
