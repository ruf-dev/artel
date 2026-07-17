import {useState} from "react"

import cls from "@/components/JsonBlock/JsonBlock.module.css"
import {cn} from "@/app/utils/cn.ts"
import JsonNode from "@/components/JsonBlock/components/JsonNode.tsx"

interface Props {
    label: string
    value: unknown
    defaultCollapsed?: boolean
}

const ROOT_PATH = "$"

function isNonEmptyBranch(value: unknown): boolean {
    if (Array.isArray(value)) return value.length > 0
    if (value !== null && typeof value === "object") return Object.keys(value).length > 0
    return false
}

export default function JsonBlock({label, value, defaultCollapsed}: Props) {
    const [collapsedPaths, setCollapsedPaths] = useState<Set<string>>(
        () => defaultCollapsed ? new Set([ROOT_PATH]) : new Set(),
    )

    function toggle(path: string) {
        setCollapsedPaths(prev => {
            const next = new Set(prev)
            if (next.has(path)) next.delete(path)
            else next.add(path)
            return next
        })
    }

    const toggleable = isNonEmptyBranch(value)
    const rootOpen = !collapsedPaths.has(ROOT_PATH)

    return (
        <div className={cls.JsonBlockContainer}>
            <div
                className={cn(cls.LabelRow, toggleable && cls.LabelRowToggle)}
                onClick={toggleable ? () => toggle(ROOT_PATH) : undefined}
            >
                {toggleable && <span className={cn(cls.Caret, rootOpen && cls.CaretOpen)}>▸</span>}
                <span className={cls.Label}>{label}</span>
            </div>
            <div className={cls.Json}>
                <JsonNode
                    value={value}
                    path={ROOT_PATH}
                    depth={0}
                    collapsedPaths={collapsedPaths}
                    onToggle={toggle}
                    trailingComma={false}
                />
            </div>
        </div>
    )
}
