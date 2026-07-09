import {
    ActionStep as PbActionStep,
    ConditionStep as PbConditionStep,
    GroupStep as PbGroupStep,
    ParallelStep as PbParallelStep,
    TractCondition as PbTractCondition,
    TractDefinition as PbTractDefinition,
    TractItem,
    TractRunItem,
    TractRunStepItem,
    TractStep as PbTractStep,
    TractToolItem,
    TriggerItem,
    TriggerSourceItem,
} from "@/app/api/artel/tracts.pb.ts"

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

function safeParseJson<T>(raw: string | undefined, fallback: T): T {
    if (!raw) return fallback
    try {
        return JSON.parse(raw) as T
    } catch {
        return fallback
    }
}

const emptySchema: SchemaNode = {properties: {}}

function conditionToProto(c: TractCondition): PbTractCondition {
    return {left: c.left, op: c.op, right: c.right}
}

function conditionFromProto(c: PbTractCondition): TractCondition {
    return {left: c.left ?? "", op: (c.op ?? "==") as TractCondition["op"], right: c.right ?? ""}
}

function stepToProto(s: TractStep): PbTractStep {
    const base = {id: s.id, name: s.name, description: s.description}
    switch (s.type) {
        case "action":
            return {...base, action: {mcp: s.mcp, tool: s.tool, connectionUuid: s.connection_uuid, params: s.params}}
        case "condition":
            return {
                ...base,
                condition: {
                    conditions: (s.conditions ?? []).map(conditionToProto),
                    then: (s.then ?? []).map(stepToProto),
                    else: (s.else ?? []).map(stepToProto),
                },
            }
        case "parallel":
            return {...base, parallel: {steps: (s.steps ?? []).map(stepToProto)}}
        case "group":
            return {...base, group: {steps: (s.steps ?? []).map(stepToProto)}}
    }
}

function stepFromProto(s: PbTractStep): TractStep {
    const base = {id: s.id ?? "", name: s.name, description: s.description}

    if (s.action) {
        const action: PbActionStep = s.action
        return {...base, type: "action", mcp: action.mcp, tool: action.tool, connection_uuid: action.connectionUuid, params: action.params}
    }
    if (s.condition) {
        const condition: PbConditionStep = s.condition
        return {
            ...base,
            type: "condition",
            conditions: (condition.conditions ?? []).map(conditionFromProto),
            then: (condition.then ?? []).map(stepFromProto),
            else: (condition.else ?? []).map(stepFromProto),
        }
    }
    if (s.parallel) {
        const parallel: PbParallelStep = s.parallel
        return {...base, type: "parallel", steps: (parallel.steps ?? []).map(stepFromProto)}
    }

    const group: PbGroupStep | undefined = s.group
    return {...base, type: "group", steps: (group?.steps ?? []).map(stepFromProto)}
}

export function definitionToProto(def: TractDefinition): PbTractDefinition {
    return {steps: def.steps.map(stepToProto)}
}

export function definitionFromProto(def: PbTractDefinition | undefined): TractDefinition {
    return {steps: (def?.steps ?? []).map(stepFromProto)}
}

function parseSchema(raw: string | undefined): SchemaNode {
    return safeParseJson(raw, emptySchema)
}

export function toTract(item: TractItem): Tract {
    return {
        uuid: item.uuid ?? "",
        name: item.name ?? "",
        description: item.description ?? "",
        enabled: item.enabled ?? false,
        definition: definitionFromProto(item.definition),
        triggers: item.triggers?.map(t => ({
            uuid: t.uuid ?? "",
            name: t.name ?? "",
            kind: t.kind ?? "",
            source: t.source ?? "",
        })),
        lastRun: item.lastRun ? {status: item.lastRun.status ?? "", at: item.lastRun.at ?? ""} : undefined,
        createdAt: item.createdAt ?? "",
        updatedAt: item.updatedAt ?? "",
    }
}

export function toRun(item: TractRunItem): TractRun {
    return {
        uuid: item.uuid ?? "",
        tractUuid: item.tractUuid ?? "",
        triggerUuid: item.triggerUuid ?? "",
        status: item.status ?? "",
        startedBy: item.startedBy ?? "",
        triggerPayload: safeParseJson(item.triggerPayload, null),
        error: item.error ?? "",
        createdAt: item.createdAt ?? "",
        updatedAt: item.updatedAt ?? "",
    }
}

export function toRunStep(item: TractRunStepItem): TractRunStep {
    return {
        stepId: item.stepId ?? "",
        stepName: item.stepName ?? "",
        stepType: item.stepType ?? "",
        status: item.status ?? "",
        input: item.input ? safeParseJson(item.input, undefined) : undefined,
        output: item.output ? safeParseJson(item.output, undefined) : undefined,
        error: item.error ?? "",
        startedAt: item.startedAt ?? "",
        finishedAt: item.finishedAt ?? "",
    }
}

export function toTool(item: TractToolItem): TractTool {
    return {
        mcp: item.mcp ?? "",
        tool: item.tool ?? "",
        description: item.description ?? "",
        inputSchema: parseSchema(item.inputSchema),
        outputSchema: parseSchema(item.outputSchema),
    }
}

export function toTriggerSource(item: TriggerSourceItem): TriggerSource {
    return {
        key: item.key ?? "",
        description: item.description ?? "",
        payloadSchema: parseSchema(item.payloadSchema),
        category: item.category ?? "",
        label: item.label ?? "",
        provider: item.provider ?? "",
    }
}

export function toTrigger(item: TriggerItem): Trigger {
    return {
        uuid: item.uuid ?? "",
        name: item.name ?? "",
        kind: item.kind ?? "",
        source: item.source ?? "",
        config: safeParseJson(item.config, null),
        payloadSchema: parseSchema(item.payloadSchema),
        triggerUuid: item.triggerUuid ?? "",
        enabled: item.enabled ?? false,
        createdAt: item.createdAt ?? "",
        tokenSuffix: item.tokenSuffix ?? "",
    }
}
