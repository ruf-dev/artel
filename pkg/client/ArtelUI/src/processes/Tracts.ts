import {
    ActionStep as PbActionStep,
    ConditionStep as PbConditionStep,
    CreateTriggerRequest,
    GroupStep as PbGroupStep,
    ParallelStep as PbParallelStep,
    TractCondition as PbTractCondition,
    TractDefinition as PbTractDefinition,
    TractItem,
    TractRunItem,
    TractRunStepItem,
    TractsAPI,
    TractStep as PbTractStep,
    TractToolItem,
    TriggerItem,
    TriggerSourceItem,
} from "@/app/api/artel/tracts.pb.ts"
import useUser from "@/hooks/user/User.ts"

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

function definitionToProto(def: TractDefinition): PbTractDefinition {
    return {steps: def.steps.map(stepToProto)}
}

function definitionFromProto(def: PbTractDefinition | undefined): TractDefinition {
    return {steps: (def?.steps ?? []).map(stepFromProto)}
}

function parseSchema(raw: string | undefined): SchemaNode {
    return safeParseJson(raw, emptySchema)
}

function toTract(item: TractItem): Tract {
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

function toRun(item: TractRunItem): TractRun {
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

function toRunStep(item: TractRunStepItem): TractRunStep {
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

function toTool(item: TractToolItem): TractTool {
    return {
        mcp: item.mcp ?? "",
        tool: item.tool ?? "",
        description: item.description ?? "",
        inputSchema: parseSchema(item.inputSchema),
        outputSchema: parseSchema(item.outputSchema),
    }
}

function toTriggerSource(item: TriggerSourceItem): TriggerSource {
    return {
        key: item.key ?? "",
        description: item.description ?? "",
        payloadSchema: parseSchema(item.payloadSchema),
        category: item.category ?? "",
        label: item.label ?? "",
        provider: item.provider ?? "",
    }
}

function toTrigger(item: TriggerItem): Trigger {
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
    }
}

export interface ITractsService {
    listTracts: () => Promise<Tract[]>
    getTract: (uuid: string) => Promise<Tract>
    createTract: (name: string, description: string, definition: TractDefinition) => Promise<{ tract: Tract; warnings: string[] }>
    updateTract: (uuid: string, name: string, description: string, definition: TractDefinition) => Promise<{ tract: Tract; warnings: string[] }>
    deleteTract: (uuid: string) => Promise<void>
    setTractEnabled: (uuid: string, enabled: boolean) => Promise<void>
    runTract: (tractUuid: string, params: unknown) => Promise<void>

    listRuns: (tractUuid: string, limit: number) => Promise<TractRun[]>
    getRun: (runUuid: string) => Promise<{ run: TractRun; steps: TractRunStep[] }>

    listTractTools: () => Promise<TractTool[]>
    listTriggerSources: () => Promise<TriggerSource[]>

    createTrigger: (name: string, kind: string, source: string, config: unknown, payloadSchema: unknown) => Promise<CreatedTrigger>
    listTriggers: () => Promise<Trigger[]>
    deleteTrigger: (uuid: string) => Promise<void>
    setTriggerEnabled: (uuid: string, enabled: boolean) => Promise<void>
    rotateTriggerToken: (uuid: string) => Promise<CreatedTrigger>
    linkTrigger: (triggerUuid: string, tractUuid: string, filters: TractCondition[]) => Promise<void>
    unlinkTrigger: (triggerUuid: string, tractUuid: string) => Promise<void>
}

export class TractsService implements ITractsService {
    async listTracts(): Promise<Tract[]> {
        const res = await TractsAPI.ListTracts({}, useUser.getState().auth.getInitReq())
        return (res.tracts ?? []).map(toTract)
    }

    async getTract(uuid: string): Promise<Tract> {
        const res = await TractsAPI.GetTract({uuid}, useUser.getState().auth.getInitReq())
        return toTract(res.tract!)
    }

    async createTract(name: string, description: string, definition: TractDefinition): Promise<{ tract: Tract; warnings: string[] }> {
        const res = await TractsAPI.CreateTract(
            {name, description, definition: definitionToProto(definition)},
            useUser.getState().auth.getInitReq(),
        )
        return {tract: toTract(res.tract!), warnings: res.warnings ?? []}
    }

    async updateTract(uuid: string, name: string, description: string, definition: TractDefinition): Promise<{ tract: Tract; warnings: string[] }> {
        const res = await TractsAPI.UpdateTract(
            {uuid, name, description, definition: definitionToProto(definition)},
            useUser.getState().auth.getInitReq(),
        )
        return {tract: toTract(res.tract!), warnings: res.warnings ?? []}
    }

    async deleteTract(uuid: string): Promise<void> {
        await TractsAPI.DeleteTract({uuid}, useUser.getState().auth.getInitReq())
    }

    async setTractEnabled(uuid: string, enabled: boolean): Promise<void> {
        await TractsAPI.SetTractEnabled({uuid, enabled}, useUser.getState().auth.getInitReq())
    }

    async runTract(tractUuid: string, params: unknown): Promise<void> {
        await TractsAPI.RunTract(
            {tractUuid, params: JSON.stringify(params ?? {})},
            useUser.getState().auth.getInitReq(),
        )
    }

    async listRuns(tractUuid: string, limit: number): Promise<TractRun[]> {
        const res = await TractsAPI.ListRuns({tractUuid, limit}, useUser.getState().auth.getInitReq())
        return (res.runs ?? []).map(toRun)
    }

    async getRun(runUuid: string): Promise<{ run: TractRun; steps: TractRunStep[] }> {
        const res = await TractsAPI.GetRun({runUuid}, useUser.getState().auth.getInitReq())
        return {run: toRun(res.run!), steps: (res.steps ?? []).map(toRunStep)}
    }

    async listTractTools(): Promise<TractTool[]> {
        const res = await TractsAPI.ListTractTools({}, useUser.getState().auth.getInitReq())
        return (res.tools ?? []).map(toTool)
    }

    async listTriggerSources(): Promise<TriggerSource[]> {
        const res = await TractsAPI.ListTriggerSources({}, useUser.getState().auth.getInitReq())
        return (res.sources ?? []).map(toTriggerSource)
    }

    async createTrigger(name: string, kind: string, source: string, config: unknown, payloadSchema: unknown): Promise<CreatedTrigger> {
        const req: CreateTriggerRequest = {
            name,
            kind,
            source,
            config: JSON.stringify(config ?? {}),
            payloadSchema: JSON.stringify(payloadSchema ?? {}),
        }
        const res = await TractsAPI.CreateTrigger(req, useUser.getState().auth.getInitReq())
        return {trigger: toTrigger(res.trigger!), webhookUrl: res.webhookUrl ?? "", webhookToken: res.webhookToken ?? ""}
    }

    async listTriggers(): Promise<Trigger[]> {
        const res = await TractsAPI.ListTriggers({}, useUser.getState().auth.getInitReq())
        return (res.triggers ?? []).map(toTrigger)
    }

    async deleteTrigger(uuid: string): Promise<void> {
        await TractsAPI.DeleteTrigger({uuid}, useUser.getState().auth.getInitReq())
    }

    async setTriggerEnabled(uuid: string, enabled: boolean): Promise<void> {
        await TractsAPI.SetTriggerEnabled({uuid, enabled}, useUser.getState().auth.getInitReq())
    }

    async rotateTriggerToken(uuid: string): Promise<CreatedTrigger> {
        const res = await TractsAPI.RotateTriggerToken({uuid}, useUser.getState().auth.getInitReq())
        return {trigger: toTrigger(res.trigger!), webhookUrl: res.webhookUrl ?? "", webhookToken: res.webhookToken ?? ""}
    }

    async linkTrigger(triggerUuid: string, tractUuid: string, filters: TractCondition[]): Promise<void> {
        await TractsAPI.LinkTrigger(
            {triggerUuid, tractUuid, filters: JSON.stringify(filters ?? [])},
            useUser.getState().auth.getInitReq(),
        )
    }

    async unlinkTrigger(triggerUuid: string, tractUuid: string): Promise<void> {
        await TractsAPI.UnlinkTrigger({triggerUuid, tractUuid}, useUser.getState().auth.getInitReq())
    }
}

export const tractsService = new TractsService()
