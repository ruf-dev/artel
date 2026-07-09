import {TemplateSource} from "@/components/TemplateInput/TemplateInput.tsx"
import {SchemaNode, TractStep, TractTool} from "@/processes/Tracts.ts"
import {computeVisibleStepIds, findStepById} from "@/processes/tractTemplate.ts"

export function buildSourcesFor(rootSteps: TractStep[], targetId: string, tools: TractTool[], triggerSchema?: SchemaNode): TemplateSource[] {
    const visibleIds = computeVisibleStepIds(rootSteps, targetId)
    const sources: TemplateSource[] = [{id: "trigger", label: "trigger", schema: triggerSchema}]

    for (const id of visibleIds) {
        const step = findStepById(rootSteps, id)
        if (!step) continue
        if (step.type === "action") {
            const tool = tools.find(t => t.mcp === step.mcp && t.tool === step.tool)
            sources.push({id, label: step.name || id, schema: tool?.outputSchema})
        } else if (step.type === "condition") {
            sources.push({id, label: step.name || id, schema: {properties: {result: {type: "boolean", description: "condition result"}}}})
        }
    }

    return sources
}

export function collectIdsFromRoot(rootSteps: TractStep[]): Set<string> {
    const ids = new Set<string>()
    function walk(list: TractStep[]) {
        for (const s of list) {
            ids.add(s.id)
            if (s.then) walk(s.then)
            if (s.else) walk(s.else)
            if (s.steps) walk(s.steps)
        }
    }
    walk(rootSteps)
    return ids
}
