import Textarea from "@/components/atoms/Textarea/Textarea.tsx"
import ComposerRow from "@/pages/workbench/components/Chat/components/ComposerRow/ComposerRow.tsx"
import ComposerChipRow from "@/pages/workbench/components/Chat/components/ComposerChipRow/ComposerChipRow.tsx"
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
    function handleKeyDown(e: React.KeyboardEvent<HTMLTextAreaElement>) {
        if (e.key !== "Enter" || e.shiftKey) return
        e.preventDefault()
        if (props.value.trim() && !props.disabled) props.onSend()
    }

    return (
        <div className={cls.ChatComposerContainer}>
            <div className={cls.Composer}>
                <div className={cls.CtxRow}/>
                <Textarea
                    value={props.value}
                    setValue={props.onChange}
                    onKeyDown={handleKeyDown}
                    placeholder={props.placeholder}
                    disabled={props.disabled}
                    autoGrow
                    className={cls.Input}
                />
                <ComposerRow
                    onNewChat={props.onNewChat}
                    onSend={props.onSend}
                    disabled={props.disabled}
                    sendDisabled={props.disabled || !props.value.trim()}
                    hideNewChatButton={props.hideNewChatButton}
                />
            </div>
            <ComposerChipRow/>
        </div>
    )
}
