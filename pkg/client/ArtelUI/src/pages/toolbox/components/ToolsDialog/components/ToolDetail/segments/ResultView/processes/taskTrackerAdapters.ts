export interface TaskTrackerTable {
    columns: string[]
    rows: Record<string, unknown>[]
}

// Single-line preview for a cell's raw value — used as-is for scalars, and as the collapsed
// preview text for object/array values (see TaskTrackerCell, which offers a pretty-printed
// expansion for those instead of falling back to raw JSON for the whole table).
export function stringifyCell(value: unknown): string {
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
// to raw JSON for a recognized task-tracker source. Cell values are kept raw (not stringified)
// so TaskTrackerCell can tell JSON values apart from scalars and offer an expand toggle.
export function toGenericTable(value: unknown): TaskTrackerTable {
    const items = Array.isArray(value) ? value : [value]
    if (items.length === 0) return {columns: [], rows: []}

    if (items.every(isPlainObject)) {
        const objectItems = items as Record<string, unknown>[]
        const columns = Array.from(new Set(objectItems.flatMap(item => Object.keys(item))))
        const rows = objectItems.map(item => Object.fromEntries(columns.map(c => [c, item[c]])))
        return {columns, rows}
    }

    return {
        columns: ["value"],
        rows: items.map(item => ({value: item})),
    }
}

// Pivots the table so fields become rows and records become columns — e.g. a single card
// {id, name, title} renders as one row per field instead of one column per field, which reads
// better when there are few records but many fields.
export function transposeTable(table: TaskTrackerTable): TaskTrackerTable {
    const columns = ["Field", ...table.rows.map((_, i) => `#${i + 1}`)]
    const rows = table.columns.map(field => {
        const row: Record<string, unknown> = {Field: field}
        table.rows.forEach((sourceRow, i) => {
            row[`#${i + 1}`] = sourceRow[field]
        })
        return row
    })
    return {columns, rows}
}
