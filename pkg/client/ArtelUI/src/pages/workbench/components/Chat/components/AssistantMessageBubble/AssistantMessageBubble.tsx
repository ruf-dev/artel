import {useMemo} from "react"
import {marked} from "marked"
import DOMPurify from "dompurify"

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
        <div className={cn(cls.AssistantMessageBubbleContainer, !done && cls.Streaming)}>
            <div className={cls.Text} dangerouslySetInnerHTML={{__html: html}}/>
        </div>
    )
}
