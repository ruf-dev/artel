import {cn} from "@/app/utils/cn.ts"
import Textarea from "@/components/atoms/Textarea/Textarea.tsx"
import ComposerRow from "@/pages/workbench/components/Chat/components/ComposerRow/ComposerRow.tsx"
import ComposerChipRow from "@/pages/workbench/components/Chat/components/ComposerChipRow/ComposerChipRow.tsx"
import ComposerCtxRow from "@/pages/workbench/components/Chat/components/ComposerCtxRow/ComposerCtxRow.tsx"
import type {TweaksSection} from "@/pages/workbench/processes/workbenchContext.ts"
import cls from "@/pages/workbench/components/Chat/components/ChatComposer/ChatComposer.module.css"

// Must match SendButton/NewChatButton's width/height in ComposerRow.module.css —
// the empty textarea is floored to this height so it centers exactly against them.
const BUTTON_SIZE_REM = 2.125

interface Props {
    value: string
    onChange: (value: string) => void
    onSend: () => void
    onNewChat: () => void
    disabled: boolean
    placeholder: string
    hideNewChatButton?: boolean
    onOpenTweaks: (section?: TweaksSection) => void
    attachedPaths: string[]
    onRemoveAttachment: (path: string) => void
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
                <ComposerCtxRow attachedPaths={props.attachedPaths} onRemoveAttachment={props.onRemoveAttachment}/>
                <div className={cls.InputRow}>
                    <Textarea
                        value={props.value}
                        setValue={props.onChange}
                        onKeyDown={handleKeyDown}
                        placeholder={props.placeholder}
                        disabled={props.disabled}
                        autoGrow
                        minHeightRem={BUTTON_SIZE_REM}
                        className={cn(cls.Input, !props.hideNewChatButton && cls.InputWide)}
                    />
                    <ComposerRow
                        onNewChat={props.onNewChat}
                        onSend={props.onSend}
                        disabled={props.disabled}
                        sendDisabled={props.disabled || !props.value.trim()}
                        hasText={!!props.value.trim()}
                        hideNewChatButton={props.hideNewChatButton}
                    />
                </div>
            </div>
            <ComposerChipRow onOpenTweaks={props.onOpenTweaks}/>
        </div>
    )
}
