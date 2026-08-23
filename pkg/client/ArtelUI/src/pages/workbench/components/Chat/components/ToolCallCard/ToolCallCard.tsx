import {useState} from "react"
import {AnimatePresence, motion} from "framer-motion"

import cls from "@/pages/workbench/components/Chat/components/ToolCallCard/ToolCallCard.module.css"
import {cn} from "@/app/utils/cn.ts"
import JsonBlock from "@/components/JsonBlock/JsonBlock.tsx"

interface Props {
    toolName: string
    input?: unknown
    inputSummary?: string
    output?: string
    isError?: boolean
    done: boolean
}

export default function ToolCallCard({toolName, input, inputSummary, output, isError, done}: Props) {
    const [expanded, setExpanded] = useState(false)
    const hasDetail = input !== undefined || Boolean(output)
    const toggle = hasDetail ? () => setExpanded(e => !e) : undefined
    const statusClass = !done ? cls.DotPending : isError ? cls.DotErr : cls.DotOk

    return (
        <motion.div
            className={cls.ToolCallCardContainer}
            initial={{opacity: 0, y: 14}}
            animate={{opacity: 1, y: 0}}
            exit={{opacity: 0}}
            transition={{duration: 0.22, ease: "easeOut"}}
        >
            <div
                className={cn(cls.Line, toggle && cls.LineClickable)}
                onClick={toggle}
                role={toggle ? "button" : undefined}
                tabIndex={toggle ? 0 : undefined}
            >
                <span className={cn(cls.Dot, statusClass)}/>
                <span className={cls.ToolName}>{toolName}</span>
                {inputSummary && <span className={cls.Summary}>{inputSummary}</span>}
                {!done && <span className={cls.Pending}>running…</span>}
                {hasDetail && <span className={cn(cls.Caret, expanded && cls.CaretOpen)}>▸</span>}
            </div>
            <AnimatePresence initial={false}>
                {expanded && (
                    <motion.div
                        key="detail"
                        className={cls.Detail}
                        initial={{height: 0, opacity: 0}}
                        animate={{height: "auto", opacity: 1}}
                        exit={{height: 0, opacity: 0}}
                        transition={{duration: 0.18, ease: "easeInOut"}}
                    >
                        {input !== undefined && <JsonBlock label="Input" value={input} defaultCollapsed/>}
                        {output && <p className={cls.OutputLabel}>Output</p>}
                        {output && <pre className={cn(cls.Output, isError && cls.OutputError)}>{output}</pre>}
                    </motion.div>
                )}
            </AnimatePresence>
        </motion.div>
    )
}
