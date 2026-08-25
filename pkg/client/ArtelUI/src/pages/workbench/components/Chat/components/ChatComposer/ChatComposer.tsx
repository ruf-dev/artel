import {Button, Input} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import PlusIcon from "@/icons/common/PlusIcon.tsx"
import cls from "@/pages/workbench/components/Chat/components/ChatComposer/ChatComposer.module.css"

interface Props {
    value: string
    onChange: (value: string) => void
    onSend: () => void
    onNewChat: () => void
    disabled: boolean
    placeholder: string
    hideNewChatButton?: boolean
}

// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function ChatComposer(props: Props) {
    function handleKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
        if (e.key === "Enter" && props.value.trim() && !props.disabled) {
            props.onSend()
        }
    }

    return (
        <div className={cls.ChatComposerContainer}>
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
            <Input
                value={props.value}
                setValue={props.onChange}
                onKeyDown={handleKeyDown}
                placeholder={props.placeholder}
                disabled={props.disabled}
                className={cls.InputWrapper}
            />
            <Button
                variant="primary"
                className={cls.SendButton}
                onClick={props.onSend}
                disabled={props.disabled || !props.value.trim()}
                aria-label="Send message"
            >
                Send
            </Button>
        </div>
    )
}
