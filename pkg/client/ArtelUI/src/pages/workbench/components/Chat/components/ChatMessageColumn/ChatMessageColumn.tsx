import {AnimatePresence} from "framer-motion"

import cls from "@/pages/workbench/components/Chat/components/ChatMessageColumn/ChatMessageColumn.module.css"
import ChatMessageItem from "@/pages/workbench/components/Chat/components/ChatMessageItem/ChatMessageItem.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"

interface Props {
    items: ChatItem[]
    assistantLabel?: string
    retryDisabled?: boolean
    onRetryMessage?: (text: string) => void
    onPermissionDecision: (id: string, decision: PermissionDecision) => void
}

// Nearest preceding user_message text for the assistant item at `index` — the
// turn that a Retry action re-sends. Undefined when the assistant message has no
// user turn before it (shouldn't happen in practice, but keeps Retry disabled).
function precedingUserText(items: ChatItem[], index: number): string | undefined {
    for (let i = index - 1; i >= 0; i--) {
        const it = items[i]
        if (it.kind === "user_message") return it.text
    }
    return undefined
}

// Centered reading column inside the full-width scroll container.
export default function ChatMessageColumn(props: Props) {
    return (
        <div className={cls.ChatMessageColumnContainer}>
            <AnimatePresence initial={false}>
                {props.items.map((item, index) => {
                    const retryText = item.kind === "assistant_message"
                        ? precedingUserText(props.items, index)
                        : undefined
                    return (
                        <ChatMessageItem
                            key={item.key}
                            item={item}
                            assistantLabel={props.assistantLabel}
                            onCopy={() => navigator.clipboard?.writeText(
                                item.kind === "assistant_message" ? item.text : "",
                            )}
                            onRetry={() => retryText && props.onRetryMessage?.(retryText)}
                            retryDisabled={props.retryDisabled || !retryText}
                            onPermissionDecision={props.onPermissionDecision}
                        />
                    )
                })}
            </AnimatePresence>
        </div>
    )
}
