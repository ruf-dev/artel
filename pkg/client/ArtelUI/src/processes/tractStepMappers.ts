// TractStep <-> proto conversion — split out of tractsMappers.ts purely to keep that file under
// the project's max-lines-per-file limit. Everything here is specific to the recursive step tree
// (condition/action/parallel/group/script/llm_call); tractsMappers.ts covers the flatter
// Tract/Run/Trigger/Template item conversions and re-exports definitionToProto/definitionFromProto
// from here for its own use and for external callers (e.g. TractsService).

import {
    ActionStep as PbActionStep,
    ConditionStep as PbConditionStep,
    GroupStep as PbGroupStep,
    LlmCallStep as PbLlmCallStep,
    ParallelStep as PbParallelStep,
    ScriptLanguage,
    ScriptStep as PbScriptStep,
    TractCondition as PbTractCondition,
    TractDefinition as PbTractDefinition,
    TractStep as PbTractStep,
} from "@/app/api/artel/tracts.pb.ts"
import {ScriptParam, TractCondition, TractDefinition, TractStep} from "@/processes/tractsTypes.ts"
import {safeParseJson} from "@/app/utils/safeParseJson.ts"

function scriptParamsToJson(params: ScriptParam[] | undefined): string {
    return JSON.stringify(params ?? [])
}

function scriptParamsFromJson(raw: string | undefined): ScriptParam[] {
    return safeParseJson<ScriptParam[]>(raw, [])
}

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
        case "script":
            return {
                ...base,
                script: {
                    language: s.language,
                    code: s.code,
                    inputParams: scriptParamsToJson(s.inputParams),
                    outputParams: scriptParamsToJson(s.outputParams),
                    params: s.params,
                },
            }
        case "llm_call":
            return {
                ...base,
                llmCall: {
                    connectionId: s.llmConnectionId,
                    model: s.llmModel,
                    prompt: s.prompt,
                    systemPrompt: s.systemPrompt,
                    maxTokens: s.maxTokens,
                },
            }
    }
}

function stepFromProto(s: PbTractStep): TractStep {
    const base = {id: s.id ?? "", name: s.name, description: s.description}

    if (s.action) {
        const action: PbActionStep = s.action
        return {
            ...base,
            type: "action",
            mcp: action.mcp,
            tool: action.tool,
            connection_uuid: action.connectionUuid,
            params: action.params,
        }
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
    if (s.script) {
        const script: PbScriptStep = s.script
        return {
            ...base,
            type: "script",
            language: script.language ?? ScriptLanguage.SCRIPT_LANGUAGE_UNSPECIFIED,
            code: script.code,
            inputParams: scriptParamsFromJson(script.inputParams),
            outputParams: scriptParamsFromJson(script.outputParams),
            params: script.params,
        }
    }
    if (s.llmCall) {
        const llmCall: PbLlmCallStep = s.llmCall
        return {
            ...base,
            type: "llm_call",
            llmConnectionId: llmCall.connectionId,
            llmModel: llmCall.model,
            prompt: llmCall.prompt,
            systemPrompt: llmCall.systemPrompt,
            maxTokens: llmCall.maxTokens,
        }
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
