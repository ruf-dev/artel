import {
    TractItem,
    TractRunItem,
    TractRunStepItem,
    TractTemplateItem as PbTractTemplateItem,
    TractTemplateSummary as PbTractTemplateSummary,
    TractToolItem,
    TriggerItem,
    TriggerSourceItem,
} from "@/app/api/artel/tracts.pb.ts"
import {
    SchemaNode,
    Tract,
    TractRun,
    TractRunStep,
    TractTemplate,
    TractTemplateSummary,
    TractTool,
    Trigger,
    TriggerSource,
} from "@/processes/tractsTypes.ts"
import {definitionFromProto, definitionToProto} from "@/processes/tractStepMappers.ts"
import {safeParseJson} from "@/app/utils/safeParseJson.ts"

export * from "@/processes/tractsTypes.ts"
export {definitionFromProto, definitionToProto}

const emptySchema: SchemaNode = {properties: {}}

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

export function toTractTemplateSummary(item: PbTractTemplateSummary): TractTemplateSummary {
    return {
        uuid: item.uuid ?? "",
        ownerUuid: item.ownerUuid ?? "",
        name: item.name ?? "",
        description: item.description ?? "",
        category: item.category ?? "",
        installCount: item.installCount ?? 0,
        publishedAt: item.publishedAt ?? "",
    }
}

export function toTractTemplate(item: PbTractTemplateItem): TractTemplate {
    return {
        uuid: item.uuid ?? "",
        sourceTractUuid: item.sourceTractUuid ?? "",
        ownerUuid: item.ownerUuid ?? "",
        name: item.name ?? "",
        description: item.description ?? "",
        definition: definitionFromProto(item.definition),
        category: item.category ?? "",
        installCount: item.installCount ?? 0,
        publishedAt: item.publishedAt ?? "",
        updatedAt: item.updatedAt ?? "",
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
