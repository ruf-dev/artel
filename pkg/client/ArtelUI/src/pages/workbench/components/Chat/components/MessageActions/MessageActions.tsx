import {Button} from "@vervstack/chures"

import CopyIcon from "@/pages/workbench/components/Chat/components/icons/CopyIcon.tsx"
import RetryIcon from "@/pages/workbench/components/Chat/components/icons/RetryIcon.tsx"
import cls from "@/pages/workbench/components/Chat/components/MessageActions/MessageActions.module.css"
import {cn} from "@/app/utils/cn.ts"

interface Props {
    onCopy: () => void
    onRetry: () => void
    retryDisabled?: boolean
    className?: string
}

// Hover actions under a finished assistant message. Visibility is CSS-only — the
// parent AssistantMessageBubble reveals this row on :hover.
export default function MessageActions({onCopy, onRetry, retryDisabled, className}: Props) {
    return (
        <div className={cn(cls.MessageActionsContainer, className)}>
            <Button
                variant="unstyled"
                className={cls.ActBtn}
                onClick={onCopy}
                aria-label="Copy message"
                data-tooltip-id="root-tooltip"
                data-tooltip-content="Copy"
            >
                <CopyIcon/>
            </Button>
            <Button
                variant="unstyled"
                className={cls.ActBtn}
                aria-label="Retry message"
                aria-disabled={retryDisabled || undefined}
                data-tooltip-id="root-tooltip"
                data-tooltip-content={retryDisabled ? "Nothing to retry" : "Retry"}
                onClick={() => {
                    if (!retryDisabled) onRetry()
                }}
            >
                <RetryIcon/>
            </Button>
        </div>
    )
}
