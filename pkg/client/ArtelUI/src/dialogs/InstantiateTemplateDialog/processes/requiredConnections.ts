import {TractStep} from "@/processes/Tracts.ts"

// Read-only extraction of the distinct non-builtin MoM names a template's step tree
// references — used to render one ConnectionPicker section per MoM in the instantiate
// wizard. Deliberately separate from tractSteps.ts's walkers, which serve step-tree
// editing, not this one-shot read pass over a template's definition.
const BUILTIN_MCP = "artel"

export function requiredConnections(steps: TractStep[]): string[] {
    const moms = new Set<string>()

    function walk(list: TractStep[]) {
        for (const step of list) {
            if (step.type === "action" && step.mcp && step.mcp !== BUILTIN_MCP) {
                moms.add(step.mcp)
            }
            if (step.then) walk(step.then)
            if (step.else) walk(step.else)
            if (step.steps) walk(step.steps)
        }
    }

    walk(steps)
    return Array.from(moms)
}
