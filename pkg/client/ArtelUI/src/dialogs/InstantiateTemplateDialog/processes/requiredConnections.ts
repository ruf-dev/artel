import {TractStep} from "@/processes/Tracts.ts"

// Read-only extraction of the distinct connections a template's step tree requires —
// used to render one ConnectionSection per requirement in the instantiate wizard.
// Deliberately separate from tractSteps.ts's walkers, which serve step-tree editing,
// not this one-shot read pass over a template's definition.
const BUILTIN_MCP = "artel"

// Matches domain.ProviderAnthropic (internal/domain/external_connection.go) — the
// only LLM provider supported today (see LlmConnectionStep.tsx, which hardcodes the
// same filter). Revisit if/when a second LLM provider is added.
const ANTHROPIC_KEY = "anthropic"

export type ConnectionRequirementKind = "mom" | "llm"

export interface ConnectionRequirement {
    key: string
    kind: ConnectionRequirementKind
}

export function requiredConnections(steps: TractStep[]): ConnectionRequirement[] {
    const moms = new Set<string>()
    let needsLlm = false

    function walk(list: TractStep[]) {
        for (const step of list) {
            if (step.type === "action" && step.mcp && step.mcp !== BUILTIN_MCP) {
                moms.add(step.mcp)
            }
            if (step.type === "llm_call") {
                needsLlm = true
            }
            if (step.then) walk(step.then)
            if (step.else) walk(step.else)
            if (step.steps) walk(step.steps)
        }
    }

    walk(steps)
    const requirements: ConnectionRequirement[] = Array.from(moms).map(key => ({key, kind: "mom" as const}))
    if (needsLlm) requirements.push({key: ANTHROPIC_KEY, kind: "llm"})
    return requirements
}
