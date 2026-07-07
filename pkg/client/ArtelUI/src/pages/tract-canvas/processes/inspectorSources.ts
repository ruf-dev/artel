import {SchemaNode, TractStep, TractTool} from "@/processes/Tracts.ts"
import {computeVisibleStepIds, findStepById} from "@/processes/tractTemplate.ts"
import {TemplateSource} from "@/components/TemplateInput/TemplateInput.tsx"

export function buildSources(rootSteps: TractStep[], targetId: string, tools: TractTool[], triggerSchema?: SchemaNode): TemplateSource[] {
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
