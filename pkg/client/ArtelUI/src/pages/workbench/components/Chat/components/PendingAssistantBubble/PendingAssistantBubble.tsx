import {motion} from "framer-motion"
import {Loader, Button} from "@vervstack/chures"

import cls from "@/pages/workbench/components/Chat/components/PendingAssistantBubble/PendingAssistantBubble.module.css"
import AssistantLabel from "@/pages/workbench/components/Chat/components/AssistantLabel/AssistantLabel.tsx"
import RetryIcon from "@/pages/workbench/components/Chat/components/icons/RetryIcon.tsx"

interface Props {
    label?: string
    bucket: "normal" | "slow" | "stuck"
    onRetry: () => void
    retryDisabled?: boolean
}

export default function PendingAssistantBubble(props: Props) {
    const {label, bucket, onRetry, retryDisabled} = props

    const statusText =
        bucket === "normal" ? "Thinking…" :
        bucket === "slow" ? "Still working…" :
        "This is taking longer than usual — the model may be stuck."

    return (
        <motion.div
            className={cls.PendingAssistantBubbleContainer}
            initial={{opacity: 0, y: 14}}
            animate={{opacity: 1, y: 0}}
            exit={{opacity: 0}}
            transition={{duration: 0.22, ease: "easeOut"}}
        >
            {label && <AssistantLabel label={label}/>}
            <div className={cls.StatusRow}>
                <Loader variant="arcs" size="sm" color="var(--coral)"/>
                <span className={cls.StatusText}>{statusText}</span>
            </div>
            {bucket === "stuck" && (
                <Button
                    variant="secondary"
                    className={cls.RetryButton}
                    onClick={onRetry}
                    disabled={retryDisabled}
                >
                    <RetryIcon/>
                    Retry
                </Button>
            )}
        </motion.div>
    )
}
