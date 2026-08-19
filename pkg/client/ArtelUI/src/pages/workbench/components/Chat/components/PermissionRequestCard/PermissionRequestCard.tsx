import {useState} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/pages/workbench/components/Chat/components/PermissionRequestCard/PermissionRequestCard.module.css"
import {cn} from "@/app/utils/cn.ts"
import JsonBlock from "@/components/JsonBlock/JsonBlock.tsx"
import {PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"

interface Props {
    toolName: string
    input?: unknown
    inputSummary?: string
    options: PermissionDecision[]
    decision?: PermissionDecision
    onDecide: (decision: PermissionDecision) => void
}

const CANONICAL_ORDER: PermissionDecision[] = ["allow_once", "allow_always", "deny"]

const LABELS: Record<PermissionDecision, string> = {
    allow_once: "Allow",
    allow_always: "Allow always",
    deny: "Deny",
}

const RESOLVED_LABELS: Record<PermissionDecision, string> = {
    allow_once: "Allowed",
    allow_always: "Allowed always",
    deny: "Denied",
}

const VARIANTS: Record<PermissionDecision, "primary" | "secondary" | "danger"> = {
    allow_once: "primary",
    allow_always: "secondary",
    deny: "danger",
}

export default function PermissionRequestCard({toolName, input, inputSummary, options, decision, onDecide}: Props) {
    const [expanded, setExpanded] = useState(false)
    const hasInput = input !== undefined
    const toggle = hasInput ? () => setExpanded(e => !e) : undefined

    return (
        <div className={cls.PermissionRequestCardContainer}>
            <div
                className={cn(cls.Header, toggle && cls.HeaderClickable)}
                onClick={toggle}
                role={toggle ? "button" : undefined}
                tabIndex={toggle ? 0 : undefined}
            >
                <span className={cls.Icon}>⚠</span>
                <span className={cls.Title}>Permission requested: {toolName}</span>
                {hasInput && <span className={cn(cls.Caret, expanded && cls.CaretOpen)}>▸</span>}
            </div>
            {inputSummary && <p className={cls.Summary}>{inputSummary}</p>}
            {expanded && hasInput && (
                <div className={cls.Detail}>
                    <JsonBlock label="Input" value={input} defaultCollapsed/>
                </div>
            )}
            {decision ? (
                <p className={cls.Resolved}>{RESOLVED_LABELS[decision]}</p>
            ) : (
                <div className={cls.Actions}>
                    {CANONICAL_ORDER.filter(d => options.includes(d)).map(d => (
                        <Button key={d} variant={VARIANTS[d]} onClick={() => onDecide(d)}>{LABELS[d]}</Button>
                    ))}
                </div>
            )}
        </div>
    )
}
