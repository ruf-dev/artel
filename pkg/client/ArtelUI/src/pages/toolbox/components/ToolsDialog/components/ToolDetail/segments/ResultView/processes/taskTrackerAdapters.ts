export interface TaskTrackerTable {
    columns: string[]
    rows: Record<string, string>[]
}

type TaskTrackerAdapter = (value: unknown) => TaskTrackerTable

// One adapter per MoM source name (candidate.name). Only trello exists today; a future
// tracker (Jira, Linear, ...) registers here under its own MoM name.
const TASK_TRACKER_ADAPTERS: Record<string, TaskTrackerAdapter> = {
    trello: toGenericTable,
}

export function buildTaskTrackerTable(source: string, value: unknown): TaskTrackerTable | null {
    const adapter = TASK_TRACKER_ADAPTERS[source]
    return adapter ? adapter(value) : null
}

function stringifyCell(value: unknown): string {
    if (value === null || value === undefined) return "—"
    if (Array.isArray(value)) return value.map(stringifyCell).join(", ")
    if (typeof value === "object") return JSON.stringify(value)
    return String(value)
}

function isPlainObject(value: unknown): value is Record<string, unknown> {
    return typeof value === "object" && value !== null && !Array.isArray(value)
}

// Renders any JSON value as a table: an array of objects becomes one row per item with the
// union of their keys as columns (e.g. trello's list_cards/list_boards responses). A bare
// array or single scalar/object degrades to a single "value" column so nothing falls back
// to raw JSON for a recognized task-tracker source.
function toGenericTable(value: unknown): TaskTrackerTable {
    const items = Array.isArray(value) ? value : [value]
    if (items.length === 0) return {columns: [], rows: []}

    if (items.every(isPlainObject)) {
        const objectItems = items as Record<string, unknown>[]
        const columns = Array.from(new Set(objectItems.flatMap(item => Object.keys(item))))
        const rows = objectItems.map(item => Object.fromEntries(columns.map(c => [c, stringifyCell(item[c])])))
        return {columns, rows}
    }

    return {
        columns: ["value"],
        rows: items.map(item => ({value: stringifyCell(item)})),
    }
}
