import {AnimatePresence} from "framer-motion"

import cls from "@/pages/workbench/components/Chat/components/ChatMessageColumn/ChatMessageColumn.module.css"
import ChatMessageItem from "@/pages/workbench/components/Chat/components/ChatMessageItem/ChatMessageItem.tsx"
import PendingAssistantBubble
    from "@/pages/workbench/components/Chat/components/PendingAssistantBubble/PendingAssistantBubble.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"

interface Props {
    items: ChatItem[]
    assistantLabel?: string
    retryDisabled?: boolean
    onRetryMessage?: (text: string) => void
    onPermissionDecision: (id: string, decision: PermissionDecision) => void
    pending?: {bucket: "normal" | "slow" | "stuck"; label?: string}
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

// Last user_message text in the items list, or undefined.
function lastUserText(items: ChatItem[]): string | undefined {
    for (let i = items.length - 1; i >= 0; i--) {
        const it = items[i]
        if (it.kind === "user_message") return it.text
    }
    return undefined
}

// Centered reading column inside the full-width scroll container.
export default function ChatMessageColumn(props: Props) {
    const trailingRetryText = !props.pending && (() => {
        const last = props.items[props.items.length - 1]
        if (!last) return undefined
        if (last.kind === "user_message" || last.kind === "error") {
            return lastUserText(props.items)
        }
        return undefined
    })()

    return (
        <div className={cls.ChatMessageColumnContainer}>
            <AnimatePresence initial={false}>
                {props.items.map((item, index) => {
                    const retryText = item.kind === "assistant_message"
                        ? precedingUserText(props.items, index)
                        : undefined
                    const isLastItem = index === props.items.length - 1
                    const showTrailingRetry = isLastItem && trailingRetryText
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
                            trailingRetry={showTrailingRetry ? {
                                onRetry: () => trailingRetryText && props.onRetryMessage?.(trailingRetryText),
                                disabled: props.retryDisabled || !trailingRetryText,
                            } : undefined}
                        />
                    )
                })}
                {props.pending && (
                    <PendingAssistantBubble
                        key="pending"
                        label={props.pending.label}
                        bucket={props.pending.bucket}
                        onRetry={() => {
                            const t = lastUserText(props.items)
                            if (t) props.onRetryMessage?.(t)
                        }}
                        retryDisabled={props.retryDisabled || !lastUserText(props.items)}
                    />
                )}
            </AnimatePresence>
        </div>
    )
}
