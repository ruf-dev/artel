import {useState} from "react"

import cls
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/tract-canvas/components/TractCanvasLogPanel/components/RunLog/components/LogEntryRow/LogEntryRow.module.css"
import {cn} from "@/app/utils/cn.ts"
import JsonBlock from "@/components/JsonBlock/JsonBlock.tsx"
import {LogLine} from "@/pages/tract-canvas/components/TractCanvasLogPanel/components/RunLog/processes/logLines.ts"

function cap(s: string): string {
    return s.charAt(0).toUpperCase() + s.slice(1)
}

interface Props {
    line: LogLine
}

export default function LogEntryRow({line}: Props) {
    const [expanded, setExpanded] = useState(false)
    const hasDetail = line.input !== undefined && line.input !== null || Boolean(line.error)
    const toggle = hasDetail ? () => setExpanded(e => !e) : undefined

    return (
        <div className={cn(
            cls.LogEntryRowContainer,
            line.terminal && cls.Terminal,
            line.terminal && cls[`Terminal${cap(line.kind)}`],
        )}>
            <div
                className={cn(cls.Line, toggle && cls.LineClickable)}
                onClick={toggle}
                role={toggle ? "button" : undefined}
                tabIndex={toggle ? 0 : undefined}
            >
                <span className={cn(cls.Dot, cls[`Dot${cap(line.kind)}`])}/>
                <span className={cls.LineT}>{line.t}</span>
                <span className={cls.LineTitle}>{line.title}</span>
                {line.meta && <span className={cls.LineMeta}>{line.meta}</span>}
                {line.error && !expanded && <span className={cls.LineError}>{line.error}</span>}
                {hasDetail && <span className={cn(cls.Caret, expanded && cls.CaretOpen)}>▸</span>}
            </div>
            {expanded && (
                <div className={cls.Detail}>
                    {line.error && <p className={cls.FullError}>{line.error}</p>}
                    {line.input !== undefined && line.input !== null && (
                        <JsonBlock label="Input" value={line.input} defaultCollapsed/>
                    )}
                </div>
            )}
        </div>
    )
}
