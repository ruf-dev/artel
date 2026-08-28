import {motion} from "framer-motion"
import {Button} from "@vervstack/chures"

import cls from "@/pages/workbench/components/Chat/components/UserMessageBubble/UserMessageBubble.module.css"
import RetryIcon from "@/pages/workbench/components/Chat/components/icons/RetryIcon.tsx"

interface Props {
    text: string
    showRetry?: boolean
    onRetry?: () => void
    retryDisabled?: boolean
}

export default function UserMessageBubble({text, showRetry, onRetry, retryDisabled}: Props) {
    return (
        <motion.div
            className={cls.UserMessageBubbleContainer}
            initial={{opacity: 0, y: 28, scale: 0.9}}
            animate={{opacity: 1, y: 0, scale: 1}}
            exit={{opacity: 0, scale: 0.9}}
            transition={{type: "spring", stiffness: 420, damping: 32}}
        >
            <p className={cls.Text}>{text}</p>
            {showRetry && (
                <div className={cls.RetryRow}>
                    <Button
                        variant="unstyled"
                        className={cls.ActBtn}
                        aria-label="Retry message"
                        aria-disabled={retryDisabled || undefined}
                        data-tooltip-id="root-tooltip"
                        data-tooltip-content={retryDisabled ? "Nothing to retry" : "Retry"}
                        onClick={() => {
                            if (!retryDisabled) onRetry?.()
                        }}
                    >
                        <RetryIcon/>
                    </Button>
                </div>
            )}
        </motion.div>
    )
}
