import cls from "@/components/JsonBlock/JsonBlock.module.css"
import {cn} from "@/app/utils/cn.ts"

type TokenKind = "string" | "number" | "boolean" | "null"

interface Props {
    keyLabel?: string
    value: unknown
    path: string
    depth: number
    collapsedPaths: Set<string>
    onToggle: (path: string) => void
    trailingComma: boolean
}

function primitiveKind(value: unknown): TokenKind {
    if (value === null) return "null"
    if (typeof value === "boolean") return "boolean"
    if (typeof value === "number") return "number"
    return "string"
}

function formatPrimitive(value: unknown): string {
    return typeof value === "string" ? JSON.stringify(value) : String(value)
}

function tokenClass(kind: TokenKind): string {
    switch (kind) {
        case "string": return cls.JsonString
        case "number": return cls.JsonNumber
        case "boolean": return cls.JsonBoolean
        case "null": return cls.JsonNull
    }
}

export default function JsonNode(props: Props) {
    const {keyLabel, value, path, depth, collapsedPaths, onToggle} = props
    const {trailingComma} = props
    const indent = {paddingLeft: `${depth * 0.9}rem`}
    const isArray = Array.isArray(value)
    const isBranch = isArray || (value !== null && typeof value === "object")

    if (!isBranch) {
        return (
            <div className={cls.Row} style={indent}>
                {keyLabel !== undefined && <span className={cls.JsonKey}>&quot;{keyLabel}&quot;: </span>}
                <span className={tokenClass(primitiveKind(value))}>{formatPrimitive(value)}</span>
                {trailingComma && ","}
            </div>
        )
    }

    const entries = isArray
        ? (value as unknown[]).map((v, i) => [String(i), v] as const)
        : Object.entries(value as Record<string, unknown>)
    const open = isArray ? "[" : "{"
    const close = isArray ? "]" : "}"

    if (entries.length === 0) {
        return (
            <div className={cls.Row} style={indent}>
                {keyLabel !== undefined && <span className={cls.JsonKey}>&quot;{keyLabel}&quot;: </span>}
                <span className={cls.Punct}>{open}{close}</span>
                {trailingComma && ","}
            </div>
        )
    }

    const collapsed = collapsedPaths.has(path)

    return (
        <>
            <div className={cn(cls.Row, cls.RowToggle)} style={indent} onClick={() => onToggle(path)}>
                <span className={cn(cls.Caret, !collapsed && cls.CaretOpen)}>▸</span>
                {keyLabel !== undefined && <span className={cls.JsonKey}>&quot;{keyLabel}&quot;: </span>}
                <span className={cls.Punct}>{open}</span>
                {collapsed && (
                    <span className={cls.Ellipsis}>
                        {entries.length} {isArray ? "items" : "keys"}{close}{trailingComma && ","}
                    </span>
                )}
            </div>
            {!collapsed && entries.map(([k, v], i) => (
                <JsonNode
                    key={isArray ? i : k}
                    keyLabel={isArray ? undefined : k}
                    value={v}
                    path={isArray ? `${path}[${k}]` : `${path}.${k}`}
                    depth={depth + 1}
                    collapsedPaths={collapsedPaths}
                    onToggle={onToggle}
                    trailingComma={i < entries.length - 1}
                />
            ))}
            {!collapsed && (
                <div className={cls.Row} style={indent}>
                    <span className={cls.Punct}>{close}</span>
                    {trailingComma && ","}
                </div>
            )}
        </>
    )
}
