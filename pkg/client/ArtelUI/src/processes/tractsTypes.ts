import {ScriptLanguage} from "@/app/api/artel/tracts.pb.ts"

export {ScriptLanguage}

export interface TractDefinition {
    steps: TractStep[]
}

// ScriptParam is one named, typed, ordered entry in a script step's declared input or
// output list — order drives the function signature ScriptEditor generates around the
// step's code. `type` reuses the same {type, description, enum, ...} shape as a MoM tool's
// input/output schema (SchemaProperty below), wrapped with a name because — unlike
// SchemaNode.properties, a map — a script step's params must preserve declaration order.
export interface ScriptParam {
    name: string
    type: SchemaProperty
}

export interface TractStep {
    id: string
    name?: string
    description?: string
    type: "action" | "condition" | "parallel" | "group" | "script"
    mcp?: string
    tool?: string
    connection_uuid?: string
    params?: Record<string, string>
    conditions?: TractCondition[]
    then?: TractStep[]
    else?: TractStep[]
    steps?: TractStep[]
    language?: ScriptLanguage
    code?: string
    inputParams?: ScriptParam[]
    outputParams?: ScriptParam[]
}

export interface TractCondition {
    left: string
    op: "==" | "!=" | ">" | "<" | ">=" | "<=" | "contains" | "glob" | "regex"
    right: string
}

// Shared JSON-schema wire shape for input_schema/output_schema/payload_schema fields.
export interface SchemaProperty {
    type: "string" | "integer" | "number" | "boolean" | "array" | "object"
    description?: string
    enum?: string[]
    properties?: Record<string, SchemaProperty>
    items?: SchemaProperty
    required?: string[]
}

export interface SchemaNode {
    properties: Record<string, SchemaProperty>
    required?: string[]
    isArray?: boolean
    items?: SchemaProperty
}

export interface TractTriggerSummary {
    uuid: string
    name: string
    kind: string
    source: string
}

export interface TractLastRun {
    status: string
    at: string
}

export interface Tract {
    uuid: string
    name: string
    description: string
    enabled: boolean
    definition: TractDefinition
    triggers?: TractTriggerSummary[]
    lastRun?: TractLastRun
    createdAt: string
    updatedAt: string
}

export interface TractRun {
    uuid: string
    tractUuid: string
    triggerUuid: string
    status: "running" | "done" | "failed" | string
    startedBy: "webhook" | "manual" | "mcp" | string
    triggerPayload: unknown
    error: string
    createdAt: string
    updatedAt: string
}

export interface TractRunStep {
    stepId: string
    stepName: string
    stepType: string
    status: "running" | "done" | "failed" | string
    input?: unknown
    output?: unknown
    error: string
    startedAt: string
    finishedAt: string
}

export interface TractTool {
    mcp: string
    tool: string
    description: string
    inputSchema: SchemaNode
    outputSchema: SchemaNode
}

export interface TriggerSource {
    key: string
    description: string
    payloadSchema: SchemaNode
    category: string
    label: string
    provider: string
}

export interface Trigger {
    uuid: string
    name: string
    kind: string
    source: string
    config: unknown
    payloadSchema: SchemaNode
    triggerUuid: string
    enabled: boolean
    createdAt: string
    tokenSuffix: string
}

export interface CreatedTrigger {
    trigger: Trigger
    webhookUrl: string
    webhookToken: string
}

export interface TractTemplateSummary {
    uuid: string
    ownerUuid: string
    name: string
    description: string
    category: string
    installCount: number
    publishedAt: string
}

export interface TractTemplate {
    uuid: string
    sourceTractUuid: string
    ownerUuid: string
    name: string
    description: string
    definition: TractDefinition
    category: string
    installCount: number
    publishedAt: string
    updatedAt: string
}
