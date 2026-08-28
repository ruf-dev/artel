import {motion} from "framer-motion"
import {Button} from "@vervstack/chures"

import cls from "@/pages/workbench/components/Chat/components/ErrorCard/ErrorCard.module.css"
import RetryIcon from "@/pages/workbench/components/Chat/components/icons/RetryIcon.tsx"

interface Props {
    text: string
    onRetry?: () => void
    retryDisabled?: boolean
}

export default function ErrorCard({text, onRetry, retryDisabled}: Props) {
    return (
        <motion.div
            className={cls.ErrorCardContainer}
            initial={{opacity: 0, y: 14}}
            animate={{opacity: 1, y: 0}}
            exit={{opacity: 0}}
            transition={{duration: 0.22, ease: "easeOut"}}
        >
            <span className={cls.Icon}>⨯</span>
            <div className={cls.Content}>
                <p className={cls.Text}>{text}</p>
                {onRetry && (
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
                )}
            </div>
        </motion.div>
    )
}
