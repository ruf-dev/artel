import {createContext, useContext, type Dispatch, type SetStateAction} from "react"

import {SchemaNode, SchemaProperty, TractCondition} from "@/processes/Tracts.ts"

export type AddTriggerStep = "mode" | "create" | "link"

export const FIELD_TYPES: SchemaProperty["type"][] = ["string", "integer", "number", "boolean", "array", "object"]

export type SchemaFieldRow = {
    name: string
    type: SchemaProperty["type"]
    description: string
    required: boolean
    properties?: SchemaFieldRow[] // nested fields, present when type === "object"
    items?: SchemaFieldRow // item schema, present when type === "array"
}

export function emptySchemaField(): SchemaFieldRow {
    return {name: "", type: "string", description: "", required: false}
}

function fieldToSchemaProperty(f: SchemaFieldRow): SchemaProperty {
    return {
        type: f.type,
        description: f.description || undefined,
        properties: f.type === "object" ? fieldsToSchemaProperties(f.properties ?? []) : undefined,
        items: f.type === "array" && f.items ? fieldToSchemaProperty(f.items) : undefined,
    }
}

function fieldsToSchemaProperties(rows: SchemaFieldRow[]): Record<string, SchemaProperty> {
    const properties: Record<string, SchemaProperty> = {}
    for (const f of rows) {
        if (!f.name) continue
        properties[f.name] = fieldToSchemaProperty(f)
    }
    return properties
}

export function fieldsToSchemaNode(rows: SchemaFieldRow[]): SchemaNode {
    const required = rows.filter(f => f.name && f.required).map(f => f.name)
    return {properties: fieldsToSchemaProperties(rows), required}
}

export interface AddTriggerDialogState {
    tractUuid: string
    linkedUuids: Set<string>

    step: AddTriggerStep
    setStep: Dispatch<SetStateAction<AddTriggerStep>>
    saving: boolean
    setSaving: Dispatch<SetStateAction<boolean>>

    // create-trigger fields
    name: string
    setName: Dispatch<SetStateAction<string>>
    kind: "webhook" | "manual"
    setKind: Dispatch<SetStateAction<"webhook" | "manual">>
    source: string
    setSource: Dispatch<SetStateAction<string>>
    fields: SchemaFieldRow[]
    setFields: Dispatch<SetStateAction<SchemaFieldRow[]>>

    // link-existing fields
    selectedTriggerId: string
    setSelectedTriggerId: Dispatch<SetStateAction<string>>
    filters: TractCondition[]
    setFilters: Dispatch<SetStateAction<TractCondition[]>>
}

export const AddTriggerDialogContext = createContext<AddTriggerDialogState | null>(null)

export function useAddTriggerDialog(): AddTriggerDialogState {
    const ctx = useContext(AddTriggerDialogContext)
    if (!ctx) throw new Error("useAddTriggerDialog must be used within AddTriggerDialogContext.Provider")
    return ctx
}
