import {useMemo} from "react"
import {marked} from "marked"
import DOMPurify from "dompurify"
import {motion} from "framer-motion"

import cls from "@/pages/workbench/components/Chat/components/AssistantMessageBubble/AssistantMessageBubble.module.css"
import {cn} from "@/app/utils/cn.ts"

interface Props {
    text: string
    done: boolean
}

export default function AssistantMessageBubble({text, done}: Props) {
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
            <div className={cls.Text} dangerouslySetInnerHTML={{__html: html}}/>
        </motion.div>
    )
}
