import {useMemo} from "react"
import {marked} from "marked"
import DOMPurify from "dompurify"
import {motion} from "framer-motion"

import cls from "@/pages/workbench/components/Chat/components/AssistantMessageBubble/AssistantMessageBubble.module.css"
import AssistantLabel from "@/pages/workbench/components/Chat/components/AssistantLabel/AssistantLabel.tsx"
import MessageActions from "@/pages/workbench/components/Chat/components/MessageActions/MessageActions.tsx"
import {cn} from "@/app/utils/cn.ts"

interface Props {
    text: string
    done: boolean
    label?: string
    onCopy: () => void
    onRetry: () => void
    retryDisabled?: boolean
}

export default function AssistantMessageBubble(props: Props) {
    const {text, done, label} = props

    const html = useMemo(() => {
        if (!text) return ""
        return DOMPurify.sanitize(marked.parse(text) as string)
    }, [text])

    return (
        <motion.div
            className={cn(cls.AssistantMessageBubbleContainer, !done && cls.Streaming)}
            initial={{opacity: 0, y: 14}}
            animate={{opacity: 1, y: 0}}
            exit={{opacity: 0}}
            transition={{duration: 0.22, ease: "easeOut"}}
        >
            {label && <AssistantLabel label={label}/>}
            <div className={cls.Text} dangerouslySetInnerHTML={{__html: html}}/>
            {done && (
                <MessageActions
                    className={cls.Actions}
                    onCopy={props.onCopy}
                    onRetry={props.onRetry}
                    retryDisabled={props.retryDisabled}
                />
            )}
        </motion.div>
    )
}
