import UserMessageBubble from "@/pages/workbench/components/Chat/components/UserMessageBubble/UserMessageBubble.tsx"
import AssistantMessageBubble
    from "@/pages/workbench/components/Chat/components/AssistantMessageBubble/AssistantMessageBubble.tsx"
import ToolCallCard from "@/pages/workbench/components/Chat/components/ToolCallCard/ToolCallCard.tsx"
import PermissionRequestCard
    from "@/pages/workbench/components/Chat/components/PermissionRequestCard/PermissionRequestCard.tsx"
import ErrorCard from "@/pages/workbench/components/Chat/components/ErrorCard/ErrorCard.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"

interface Props {
    item: ChatItem
    assistantLabel?: string
    onCopy: () => void
    onRetry: () => void
    retryDisabled?: boolean
    onPermissionDecision: (id: string, decision: PermissionDecision) => void
}

export default function ChatMessageItem(props: Props) {
    const {item, onPermissionDecision} = props
    switch (item.kind) {
        case "user_message":
            return <UserMessageBubble text={item.text}/>
        case "assistant_message":
            return (
                <AssistantMessageBubble
                    text={item.text}
                    done={item.done}
                    label={props.assistantLabel}
                    onCopy={props.onCopy}
                    onRetry={props.onRetry}
                    retryDisabled={props.retryDisabled}
                />
            )
        case "tool_call":
            return (
                <ToolCallCard
                    toolName={item.toolName}
                    input={item.input}
                    inputSummary={item.inputSummary}
                    output={item.output}
                    isError={item.isError}
                    done={item.done}
                />
            )
        case "permission_request":
            return (
                <PermissionRequestCard
                    toolName={item.toolName}
                    input={item.input}
                    inputSummary={item.inputSummary}
                    options={item.options}
                    decision={item.decision}
                    onDecide={decision => onPermissionDecision(item.id, decision)}
                />
            )
        case "error":
            return <ErrorCard text={item.text}/>
        default:
            return null
    }
}
