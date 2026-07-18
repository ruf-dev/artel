// Pure tree-editing utilities over the nested TractStep tree (then/else/steps), used by the
// tract editor's local (unsaved) definition state. All functions are non-mutating.

import {ScriptLanguage, TractStep} from "@/processes/Tracts.ts"
import {StepDraft} from "@/components/StepPickerDialog/StepPickerDialog.tsx"

/** Location identifies one ordered array of steps within the tree: the root list, or a named
 * branch (then/else/steps) of a specific parent step id. */
export type Location =
    | { parentId: null; branch: "root" }
    | { parentId: string; branch: "then" | "else" | "steps" }

export const ROOT_LOCATION: Location = {parentId: null, branch: "root"}

function branchArray(step: TractStep, branch: "then" | "else" | "steps"): TractStep[] {
    if (branch === "then") return step.then ?? []
    if (branch === "else") return step.else ?? []
    return step.steps ?? []
}

function withBranch(step: TractStep, branch: "then" | "else" | "steps", arr: TractStep[]): TractStep {
    if (branch === "then") return {...step, then: arr}
    if (branch === "else") return {...step, else: arr}
    return {...step, steps: arr}
}

function walkAndReplaceBranch(
    steps: TractStep[],
    parentId: string,
    branch: "then" | "else" | "steps",
    updater: (arr: TractStep[]) => TractStep[]
): TractStep[] {
    return steps.map(step => {
        if (step.id === parentId) {
            return withBranch(step, branch, updater(branchArray(step, branch)))
        }
        return {
            ...step,
            then: step.then ? walkAndReplaceBranch(step.then, parentId, branch, updater) : step.then,
            else: step.else ? walkAndReplaceBranch(step.else, parentId, branch, updater) : step.else,
            steps: step.steps ? walkAndReplaceBranch(step.steps, parentId, branch, updater) : step.steps,
        }
    })
}

/** mutateAt applies updater to the ordered array identified by location, returning a new root tree. */
export function mutateAt(
    rootSteps: TractStep[], location: Location, updater: (arr: TractStep[]) => TractStep[]
): TractStep[] {
    if (location.parentId === null) return updater(rootSteps)
    return walkAndReplaceBranch(rootSteps, location.parentId, location.branch, updater)
}

export function insertStepAt(
    rootSteps: TractStep[], location: Location, index: number, newStep: TractStep
): TractStep[] {
    return mutateAt(rootSteps, location, arr => {
        const copy = [...arr]
        copy.splice(index, 0, newStep)
        return copy
    })
}

export function appendStep(rootSteps: TractStep[], location: Location, newStep: TractStep): TractStep[] {
    return mutateAt(rootSteps, location, arr => [...arr, newStep])
}

/** insertBlockAfter inserts newStep at index within location's array, for the canvas's "add
 * step after this node" action. If a step already occupies that position (the clicked node
 * already has a "next"), newStep does NOT get spliced into the sequence — instead it becomes
 * a parallel sibling of the existing step, both forking from the same predecessor and
 * rejoining afterward, matching the canvas's fork/join visual model (see
 * pages/tract-canvas/processes/tractCanvasLayout.ts). If the existing step at that position is
 * already a "parallel" step, newStep is added as another lane instead of nesting parallels. */
export function insertBlockAfter(
    rootSteps: TractStep[],
    location: Location,
    index: number,
    newStep: TractStep,
    existingIds: Set<string>
): TractStep[] {
    return mutateAt(rootSteps, location, arr => {
        const existing = arr[index]
        if (!existing) {
            const copy = [...arr]
            copy.splice(index, 0, newStep)
            return copy
        }

        const copy = [...arr]
        if (existing.type === "parallel") {
            copy[index] = {...existing, steps: [...(existing.steps ?? []), newStep]}
            return copy
        }

        const parallelId = generateStepId("parallel", existingIds)
        copy[index] = {id: parallelId, name: parallelId, type: "parallel", steps: [existing, newStep]}
        return copy
    })
}

export function removeStepAt(rootSteps: TractStep[], location: Location, stepId: string): TractStep[] {
    const removed = mutateAt(rootSteps, location, arr => arr.filter(s => s.id !== stepId))
    return collapseThinParallels(removed)
}

/** collapseThinParallels walks the whole tree and unwraps any "parallel" step left with 0 or 1
 * lanes (e.g. after deleting a lane down to the last one or two) — a parallel step only makes
 * sense with 2+ concurrent lanes, so once it's thinner than that it should disappear (0 lanes) or
 * be replaced by its sole remaining lane in place (1 lane) rather than linger as a dead wrapper. */
export function collapseThinParallels(steps: TractStep[]): TractStep[] {
    const result: TractStep[] = []
    for (const step of steps) {
        const withCollapsedChildren: TractStep = {
            ...step,
            then: step.then ? collapseThinParallels(step.then) : step.then,
            else: step.else ? collapseThinParallels(step.else) : step.else,
            steps: step.steps ? collapseThinParallels(step.steps) : step.steps,
        }

        if (withCollapsedChildren.type === "parallel") {
            const lanes = withCollapsedChildren.steps ?? []
            if (lanes.length === 0) continue
            if (lanes.length === 1) {
                result.push(lanes[0])
                continue
            }
        }

        result.push(withCollapsedChildren)
    }
    return result
}

/** replaceStep finds stepId anywhere in the tree and replaces it via updater, preserving its
 * children unless updater itself changes them. */
export function replaceStep(
    rootSteps: TractStep[], stepId: string, updater: (step: TractStep) => TractStep
): TractStep[] {
    return rootSteps.map(step => {
        if (step.id === stepId) return updater(step)
        return {
            ...step,
            then: step.then ? replaceStep(step.then, stepId, updater) : step.then,
            else: step.else ? replaceStep(step.else, stepId, updater) : step.else,
            steps: step.steps ? replaceStep(step.steps, stepId, updater) : step.steps,
        }
    })
}

export function hasChildren(step: TractStep): boolean {
    return (step.then?.length ?? 0) > 0 || (step.else?.length ?? 0) > 0 || (step.steps?.length ?? 0) > 0
}

export function collectAllStepIds(steps: TractStep[]): Set<string> {
    const ids = new Set<string>()
    function walk(list: TractStep[]) {
        for (const s of list) {
            ids.add(s.id)
            if (s.then) walk(s.then)
            if (s.else) walk(s.else)
            if (s.steps) walk(s.steps)
        }
    }
    walk(steps)
    return ids
}

function slugify(name: string): string {
    const cleaned = name.toLowerCase().replace(/[^a-z0-9_]+/g, "_").replace(/^_+|_+$/g, "").replace(/_+/g, "_")
    if (!cleaned) return "step"
    return /^[a-z]/.test(cleaned) ? cleaned : `s_${cleaned}`
}

/** generateStepId derives a valid, unique step id (^[a-z][a-z0-9_]*$, not "trigger") from a
 * human name (typically the tool name, or "condition"/"parallel"/"group"), suffixing _1/_2... on
 * collision with an existing id. Mirrors the server's naming convention note in 00-overview.md
 * ("Default naming ... is done by UI/agent — server only validates"). */
export function generateStepId(base: string, existingIds: Set<string>): string {
    let candidate = slugify(base)
    if (candidate === "trigger") candidate = "trigger_step"
    if (!existingIds.has(candidate)) return candidate

    let i = 1
    while (existingIds.has(`${candidate}_${i}`)) i++
    return `${candidate}_${i}`
}

/** buildStepFromDraft turns a StepPickerDialog draft into a fresh TractStep with a generated,
 * unique id/name and the shape validateStepShape requires for its type (e.g. a condition
 * always starts with one empty condition row; parallel/group start with no children — the
 * server rejects an empty parallel/group on Save, prompting the user to add a lane/step). */
export function buildStepFromDraft(draft: StepDraft, existingIds: Set<string>): TractStep {
    const baseName = draft.type === "action" ? (draft.tool ?? "action") : draft.type
    const id = generateStepId(baseName, existingIds)

    if (draft.type === "action") {
        return {
            id, name: id, type: "action", mcp: draft.mcp, tool: draft.tool,
            connection_uuid: draft.connectionUuid, params: {},
        }
    }
    if (draft.type === "condition") {
        return {id, name: id, type: "condition", conditions: [{left: "", op: "==", right: ""}], then: [], else: []}
    }
    if (draft.type === "script") {
        return {
            id, name: id, type: "script", language: ScriptLanguage.SCRIPT_LANGUAGE_JAVASCRIPT,
            code: "", params: {}, inputParams: [], outputParams: [],
        }
    }
    return {id, name: id, type: draft.type, steps: []}
}
