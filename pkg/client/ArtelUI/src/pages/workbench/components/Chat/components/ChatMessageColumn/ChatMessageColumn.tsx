import {AnimatePresence} from "framer-motion"

import cls from "@/pages/workbench/components/Chat/components/ChatMessageColumn/ChatMessageColumn.module.css"
import ChatMessageItem from "@/pages/workbench/components/Chat/components/ChatMessageItem/ChatMessageItem.tsx"
import PendingAssistantBubble
    from "@/pages/workbench/components/Chat/components/PendingAssistantBubble/PendingAssistantBubble.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"

interface Props {
    items: ChatItem[]
    vaultId: string
    assistantLabel?: string
    retryDisabled?: boolean
    onRetryMessage?: (text: string) => void
    onResendMessage?: (id: string, text: string) => void
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

// Last user_message in the items list (id + text), or undefined. The id is
// carried alongside the text so a trailing-failure retry can resend in place
// (see resendMessage in useChatSession.ts/useSimpleChatSession.ts) instead of
// minting a new message — id is undefined for a legacy persisted row from
// before ids were stored, which the caller uses to disable that retry.
function lastUserMessage(items: ChatItem[]): {id?: string; text: string} | undefined {
    for (let i = items.length - 1; i >= 0; i--) {
        const it = items[i]
        if (it.kind === "user_message") return {id: it.id, text: it.text}
    }
    return undefined
}

// Centered reading column inside the full-width scroll container.
export default function ChatMessageColumn(props: Props) {
    const trailingFailure = !props.pending && (() => {
        const last = props.items[props.items.length - 1]
        if (!last) return undefined
        if (last.kind === "user_message" || last.kind === "error") {
            return lastUserMessage(props.items)
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
                    const showTrailingRetry = isLastItem && trailingFailure
                    return (
                        <ChatMessageItem
                            key={item.key}
                            item={item}
                            vaultId={props.vaultId}
                            assistantLabel={props.assistantLabel}
                            onCopy={() => navigator.clipboard?.writeText(
                                item.kind === "assistant_message" ? item.text : "",
                            )}
                            onRetry={() => retryText && props.onRetryMessage?.(retryText)}
                            retryDisabled={props.retryDisabled || !retryText}
                            onPermissionDecision={props.onPermissionDecision}
                            trailingRetry={showTrailingRetry ? {
                                onRetry: () => trailingFailure?.id
                                    && props.onResendMessage?.(trailingFailure.id, trailingFailure.text),
                                disabled: props.retryDisabled || !trailingFailure?.id,
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
                            const t = lastUserMessage(props.items)?.text
                            if (t) props.onRetryMessage?.(t)
                        }}
                        retryDisabled={props.retryDisabled || !lastUserMessage(props.items)?.text}
                    />
                )}
            </AnimatePresence>
        </div>
    )
}
