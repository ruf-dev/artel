// TS mirror of internal/service/v1/tract/template.go's scanning logic. Only the
// extraction/validation half is ported here (no render/eval) — this task has no live template
// preview; task 08's editor is the first consumer that might need render/eval, and should not
// assume they exist in this file.

import {TractStep} from "@/processes/Tracts.ts"

interface PathSegment {
    field?: string
    index?: number
    isIndex: boolean
}

function isOpen(s: string, i: number): boolean {
    return s[i] === "{" && s[i + 1] === "{"
}

function isEscapedOpen(s: string, i: number): boolean {
    return s[i] === "$" && s[i + 1] === "{" && s[i + 2] === "{"
}

/**
 * extractRefs scans s for `{{ ... }}` tract tokens, skipping literal `${{ ... }}` MoM escapes,
 * and returns each token's raw reference expression. `length(...)` wrappers are unwrapped to
 * their inner reference; `$`-prefixed vars (e.g. `$now`) are omitted. Mirrors
 * internal/service/v1/tract/template.go's extractRefs exactly.
 */
export function extractRefs(s: string): string[] {
    const refs: string[] = []
    let i = 0
    while (i < s.length) {
        if (isEscapedOpen(s, i)) {
            const end = s.indexOf("}}", i)
            if (end === -1) throw new Error("unterminated ${{ ... }} sequence")
            i = end + 2
            continue
        }

        if (isOpen(s, i)) {
            const end = s.indexOf("}}", i)
            if (end === -1) throw new Error("unterminated {{ ... }} token")
            const expr = s.slice(i + 2, end).trim()
            if (!expr.startsWith("$")) {
                let ref = expr
                if (expr.startsWith("length(") && expr.endsWith(")")) {
                    ref = expr.slice("length(".length, expr.length - 1).trim()
                }
                refs.push(ref)
            }
            i = end + 2
            continue
        }

        i++
    }
    return refs
}

/** splitRef splits a reference expression into its base identifier and trailing path segments. */
function splitRef(ref: string): { base: string; segments: PathSegment[] } {
    if (!ref) throw new Error("empty reference")

    let i = 0
    while (i < ref.length && ref[i] !== "." && ref[i] !== "[") i++
    const base = ref.slice(0, i)
    if (!base) throw new Error(`reference missing base identifier: ${ref}`)

    return {base, segments: parsePathSegments(ref.slice(i))}
}

function parsePathSegments(rest: string): PathSegment[] {
    const segments: PathSegment[] = []
    let i = 0
    while (i < rest.length) {
        if (rest[i] === ".") {
            i++
            const start = i
            while (i < rest.length && rest[i] !== "." && rest[i] !== "[") i++
            const name = rest.slice(start, i)
            if (!name) throw new Error(`empty path segment in: ${rest}`)
            segments.push({field: name, isIndex: false})
        } else if (rest[i] === "[") {
            const end = rest.indexOf("]", i)
            if (end === -1) throw new Error(`unterminated [ index in: ${rest}`)
            const numStr = rest.slice(i + 1, end)
            const idx = Number(numStr)
            if (!Number.isInteger(idx)) throw new Error(`non-numeric index in: ${rest}`)
            segments.push({index: idx, isIndex: true})
            i = end + 1
        } else {
            throw new Error(`unexpected character in path: ${rest}`)
        }
    }
    return segments
}

/**
 * validateRef reports whether ref's base identifier is either the literal "trigger" or a
 * member of availableStepIds. This is a primitive only — task 08's editor computes which step
 * ids are actually visible (ancestors/earlier-siblings) at a given tree position and passes
 * that set in; this function does not replicate the backend's full checkVisibility tree walk.
 */
export function validateRef(ref: string, availableStepIds: string[]): boolean {
    try {
        const {base} = splitRef(ref)
        return base === "trigger" || availableStepIds.includes(base)
    } catch {
        return false
    }
}

// -- Visibility (mirrors internal/service/v1/tract/tract.go's checkVisibility) --
//
// A step may only reference "trigger" or a step id that is DEFINITELY executed before it: a
// preceding sibling at any ancestor level. Parallel lanes never see each other while running,
// but all of a parallel/group step's descendant ids fan in for steps that follow it. A
// condition step's own id becomes visible after it; ids declared inside its then/else branches
// never escape the branch.

function walkProduced(steps: TractStep[], visible: Set<string>): Set<string> {
    const local = new Set(visible)
    const produced = new Set<string>()

    for (const step of steps) {
        if (step.type === "action" || step.type === "script") {
            local.add(step.id)
            produced.add(step.id)
        } else if (step.type === "condition") {
            walkProduced(step.then ?? [], local)
            walkProduced(step.else ?? [], local)
            local.add(step.id)
            produced.add(step.id)
        } else if (step.type === "parallel") {
            const fannedIn = new Set<string>()
            for (const lane of step.steps ?? []) {
                const laneProduced = walkProduced([lane], local)
                laneProduced.forEach(id => fannedIn.add(id))
            }
            fannedIn.forEach(id => {
                local.add(id)
                produced.add(id)
            })
        } else if (step.type === "group") {
            const groupProduced = walkProduced(step.steps ?? [], local)
            groupProduced.forEach(id => {
                local.add(id)
                produced.add(id)
            })
        }
    }

    return produced
}

function findVisible(steps: TractStep[], targetId: string, visible: Set<string>): Set<string> | null {
    const local = new Set(visible)

    for (const step of steps) {
        if (step.id === targetId) return new Set(local)

        if (step.type === "action" || step.type === "script") {
            local.add(step.id)
        } else if (step.type === "condition") {
            const thenHit = findVisible(step.then ?? [], targetId, local)
            if (thenHit) return thenHit
            const elseHit = findVisible(step.else ?? [], targetId, local)
            if (elseHit) return elseHit
            local.add(step.id)
        } else if (step.type === "parallel") {
            const fannedIn = new Set<string>()
            let hit: Set<string> | null = null
            for (const lane of step.steps ?? []) {
                const laneHit = findVisible([lane], targetId, local)
                if (laneHit) {
                    hit = laneHit
                    break
                }
                const laneProduced = walkProduced([lane], local)
                laneProduced.forEach(id => fannedIn.add(id))
            }
            if (hit) return hit
            fannedIn.forEach(id => local.add(id))
        } else if (step.type === "group") {
            const groupHit = findVisible(step.steps ?? [], targetId, local)
            if (groupHit) return groupHit
            const groupProduced = walkProduced(step.steps ?? [], local)
            groupProduced.forEach(id => local.add(id))
        }
    }

    return null
}

/**
 * computeVisibleStepIds returns every step id visible to targetId (a preceding sibling at any
 * ancestor level, per the backend's checkVisibility rule) — the set TemplateInput's hint
 * dropdown and validateRef should treat as valid non-"trigger" references. Returns [] if
 * targetId isn't found in the tree (e.g. a brand-new unsaved step).
 */
export function computeVisibleStepIds(rootSteps: TractStep[], targetId: string): string[] {
    const result = findVisible(rootSteps, targetId, new Set())
    return result ? Array.from(result) : []
}

/** findStepById searches the whole step tree for a step by id, recursing into then/else/steps. */
export function findStepById(steps: TractStep[], id: string): TractStep | undefined {
    for (const step of steps) {
        if (step.id === id) return step
        const nested = findStepById(step.then ?? [], id)
            ?? findStepById(step.else ?? [], id)
            ?? findStepById(step.steps ?? [], id)
        if (nested) return nested
    }
    return undefined
}
