import cls from "@/components/JsonView/JsonView.module.css"

type JsonTokenKind = "key" | "string" | "number" | "boolean" | "null" | "plain"

interface JsonToken {
    text: string
    kind: JsonTokenKind
}

// Matches JSON string literals (optionally followed by a colon, marking them as keys),
// booleans, null, and numbers — everything else (braces, commas, whitespace) is left as plain text.
const JSON_TOKEN_RE = /"(?:\\.|[^"\\])*"(?:\s*:)?|\btrue\b|\bfalse\b|\bnull\b|-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?/g

function tokenizeJson(json: string): JsonToken[] {
    const tokens: JsonToken[] = []
    let lastIndex = 0
    for (const match of json.matchAll(JSON_TOKEN_RE)) {
        const idx = match.index
        if (idx > lastIndex) tokens.push({text: json.slice(lastIndex, idx), kind: "plain"})
        const text = match[0]
        let kind: JsonTokenKind
        if (text.startsWith('"')) {
            kind = text.trimEnd().endsWith(":") ? "key" : "string"
        } else if (text === "true" || text === "false") {
            kind = "boolean"
        } else if (text === "null") {
            kind = "null"
        } else {
            kind = "number"
        }
        tokens.push({text, kind})
        lastIndex = idx + text.length
    }
    if (lastIndex < json.length) tokens.push({text: json.slice(lastIndex), kind: "plain"})
    return tokens
}

const JSON_KIND_CLASS: Record<JsonTokenKind, string | undefined> = {
    key: cls.JsonKey,
    string: cls.JsonString,
    number: cls.JsonNumber,
    boolean: cls.JsonBoolean,
    null: cls.JsonNull,
    plain: undefined,
}

export default function JsonView({value}: { value: unknown }) {
    const json = JSON.stringify(value, null, 2).replace(/[ \t]+,/g, ",")
    return (
        <pre className={cls.JsonViewContainer}>
            {tokenizeJson(json).map((t, i) => t.kind === "plain" ? t.text : (
                <span key={i} className={JSON_KIND_CLASS[t.kind]}>{t.text}</span>
            ))}
        </pre>
    )
}
