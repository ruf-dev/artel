export interface TractDefinition {
    steps: TractStep[]
}

export interface TractStep {
    id: string
    name?: string
    description?: string
    type: "action" | "condition" | "parallel" | "group"
    mcp?: string
    tool?: string
    connection_uuid?: string
    params?: Record<string, string>
    conditions?: TractCondition[]
    then?: TractStep[]
    else?: TractStep[]
    steps?: TractStep[]
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
