/** safeParseJson parses raw JSON, falling back to `fallback` on empty/invalid input — shared by
 * every mapper that unmarshals a JSON-string field off a proto item (schemas, script params,
 * trigger payloads, ...). */
export function safeParseJson<T>(raw: string | undefined, fallback: T): T {
    if (!raw) return fallback
    try {
        return JSON.parse(raw) as T
    } catch {
        return fallback
    }
}
