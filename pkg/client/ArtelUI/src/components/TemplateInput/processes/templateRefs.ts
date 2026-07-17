import {SchemaNode, SchemaProperty} from "@/processes/Tracts.ts"

export interface TemplateSource {
    id: string
    label: string
    schema?: SchemaNode
}

export interface FlatEntry {
    key: string
    ref: string
    label: string
    type?: string
    description?: string
    depth: number
}

export interface FilteredGroup {
    source: TemplateSource
    entries: FlatEntry[]
    rootMatches: boolean
}

export interface FlatListItem {
    ref: string
    label: string
}

function flattenArrayItems(ref: string, items: SchemaProperty, depth: number, out: FlatEntry[]) {
    const arrRef = `${ref}[0]`
    out.push({
        key: arrRef, ref: arrRef, label: "[0]", type: items.type,
        description: items.description, depth,
    })
    if (items.type === "object" && items.properties) {
        for (const [childName, childProp] of Object.entries(items.properties)) {
            flattenProperty(arrRef, childName, childProp, depth + 1, out)
        }
    }
}

function flattenProperty(baseRef: string, name: string, prop: SchemaProperty, depth: number, out: FlatEntry[]) {
    const ref = `${baseRef}.${name}`
    out.push({key: ref, ref, label: name, type: prop.type, description: prop.description, depth})

    if (prop.type === "object" && prop.properties) {
        for (const [childName, childProp] of Object.entries(prop.properties)) {
            flattenProperty(ref, childName, childProp, depth + 1, out)
        }
    } else if (prop.type === "array" && prop.items) {
        flattenArrayItems(ref, prop.items, depth + 1, out)
    }
}

function flattenSource(source: TemplateSource): FlatEntry[] {
    const out: FlatEntry[] = []
    if (source.schema?.isArray && source.schema.items) {
        flattenArrayItems(source.id, source.schema.items, 1, out)
    } else if (source.schema) {
        for (const [name, prop] of Object.entries(source.schema.properties)) {
            flattenProperty(source.id, name, prop, 1, out)
        }
    }
    return out
}

export function filterSources(sources: TemplateSource[], query: string): FilteredGroup[] {
    const q = query.toLowerCase()
    return sources.map(source => ({
        source,
        entries: flattenSource(source).filter(e => !q || e.ref.toLowerCase().includes(q)),
        rootMatches: !q || source.id.toLowerCase().includes(q),
    })).filter(g => g.rootMatches || g.entries.length > 0)
}

export function buildFlatList(filtered: FilteredGroup[]): FlatListItem[] {
    const list: FlatListItem[] = []
    for (const g of filtered) {
        if (g.rootMatches) list.push({ref: g.source.id, label: g.source.label})
        for (const e of g.entries) list.push({ref: e.ref, label: e.label})
    }
    return list
}
