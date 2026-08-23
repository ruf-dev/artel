import {Button, Input} from "@vervstack/chures"

import cls from "@/pages/workbench/components/Chat/components/ChatComposer/ChatComposer.module.css"

interface Props {
    value: string
    onChange: (value: string) => void
    onSend: () => void
    onNewChat: () => void
    disabled: boolean
    placeholder: string
}

export default function ChatComposer({value, onChange, onSend, onNewChat, disabled, placeholder}: Props) {
    function handleKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
        if (e.key === "Enter" && value.trim() && !disabled) {
            onSend()
        }
    }

    return (
        <div className={cls.ChatComposerContainer}>
            <Button
                variant="secondary"
                className={cls.NewChatButton}
                onClick={onNewChat}
                disabled={disabled}
                aria-label="New chat"
                title="New chat"
            >
                <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
                     strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
            </Button>
            <Input
                value={value}
                setValue={onChange}
                onKeyDown={handleKeyDown}
                placeholder={placeholder}
                disabled={disabled}
                className={cls.InputWrapper}
            />
            <Button variant="primary" onClick={onSend} disabled={disabled || !value.trim()} aria-label="Send message">
                Send
            </Button>
        </div>
    )
}
