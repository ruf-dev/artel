import {Button, Input} from "@vervstack/chures"

import cls from "@/pages/workbench/components/Chat/components/ChatComposer/ChatComposer.module.css"

interface Props {
    value: string
    onChange: (value: string) => void
    onSend: () => void
    disabled: boolean
    placeholder: string
}

export default function ChatComposer({value, onChange, onSend, disabled, placeholder}: Props) {
    function handleKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
        if (e.key === "Enter" && value.trim() && !disabled) {
            onSend()
        }
    }

    return (
        <div className={cls.ChatComposerContainer}>
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
