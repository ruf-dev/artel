// Pure tree logic for the canvas's drag-and-drop step reordering, split out of tractSteps.ts
// (which owns the general tree-editing primitives this builds on) to keep both files under the
// line-count limit.

import {TractStep} from "@/processes/Tracts.ts"
import {
    branchArray, collapseThinParallels, collectAllStepIds, insertBlockAfter, insertStepAt, Location, mutateAt,
    ROOT_LOCATION,
} from "@/processes/tractSteps.ts"

export type DropSide = "before" | "after"
export type MoveConflictChoice = "parallel" | "between"

// Must match TRIGGER_NODE_ID in pages/tract-canvas/processes/tractCanvasLayout.ts — the canvas's
// synthetic trigger node lives outside the TractStep tree, so a drop targeting it is special-cased
// below instead of resolved via tree lookup.
const TRIGGER_STEP_ID = "trigger"

export function findStepById(rootSteps: TractStep[], id: string): TractStep | null {
    for (const step of rootSteps) {
        if (step.id === id) return step
        for (const branch of [step.then, step.else, step.steps]) {
            if (!branch) continue
            const found = findStepById(branch, id)
            if (found) return found
        }
    }
    return null
}

/** findLocation returns where stepId currently lives: the array (Location) that holds it and its
 * index within that array. */
export function findLocation(rootSteps: TractStep[], stepId: string): {location: Location; index: number} | null {
    function walk(steps: TractStep[], location: Location): {location: Location; index: number} | null {
        for (let i = 0; i < steps.length; i++) {
            const step = steps[i]
            if (step.id === stepId) return {location, index: i}
            if (step.then) {
                const found = walk(step.then, {parentId: step.id, branch: "then"})
                if (found) return found
            }
            if (step.else) {
                const found = walk(step.else, {parentId: step.id, branch: "else"})
                if (found) return found
            }
            if (step.steps) {
                const found = walk(step.steps, {parentId: step.id, branch: "steps"})
                if (found) return found
            }
        }
        return null
    }
    return walk(rootSteps, ROOT_LOCATION)
}

function arrayAt(rootSteps: TractStep[], location: Location): TractStep[] {
    if (location.parentId === null) return rootSteps
    const parent = findStepById(rootSteps, location.parentId)
    if (!parent) return []
    return branchArray(parent, location.branch)
}

/** Where a "drop after target" lands. A plain node's "next" is the following index in its own
 * array. A parallel lane has no "next" of its own — dropping after a lane means "after the whole
 * parallel finishes", mirroring the nextOverride redirect in tractCanvasLayout.ts's layoutFlow. */
function locateInsertAfter(rootSteps: TractStep[], targetId: string): {location: Location; index: number} | null {
    if (targetId === TRIGGER_STEP_ID) return {location: ROOT_LOCATION, index: 0}

    const found = findLocation(rootSteps, targetId)
    if (!found) return null
    if (found.location.branch === "steps") {
        const parent = findStepById(rootSteps, found.location.parentId)
        if (parent?.type === "parallel") return locateInsertAfter(rootSteps, parent.id)
    }
    return {location: found.location, index: found.index + 1}
}

function isSelfOrDescendant(step: TractStep, id: string): boolean {
    if (step.id === id) return true
    const children = [...(step.then ?? []), ...(step.else ?? []), ...(step.steps ?? [])]
    return children.some(child => isSelfOrDescendant(child, id))
}

interface ResolvedDrop {
    steps: TractStep[]
    location: Location
    index: number
    conflict: boolean
}

/** Shared resolution for a drag-and-drop step move: removes source from wherever it currently
 * lives (collapsing any parallel left with under 2 lanes, same as delete), then locates the drop
 * point fresh against that already-collapsed tree. This makes a degenerate case like dragging one
 * lane out of a 2-lane parallel onto its own former sibling "just work" — the sibling's location
 * is resolved AFTER the collapse already inlined it, never against stale pre-removal indices. */
function resolveDrop(
    rootSteps: TractStep[], sourceId: string, targetId: string, side: DropSide
): ResolvedDrop | null {
    if (sourceId === targetId) return null
    const source = findStepById(rootSteps, sourceId)
    if (!source) return null
    if (isSelfOrDescendant(source, targetId)) return null

    const sourceLoc = findLocation(rootSteps, sourceId)
    if (!sourceLoc) return null
    const withoutSource = collapseThinParallels(
        mutateAt(rootSteps, sourceLoc.location, arr => arr.filter(s => s.id !== sourceId))
    )

    if (side === "before") {
        if (targetId === TRIGGER_STEP_ID) return null
        const target = findLocation(withoutSource, targetId)
        if (!target) return null
        return {steps: withoutSource, location: target.location, index: target.index, conflict: false}
    }

    const insertAt = locateInsertAfter(withoutSource, targetId)
    if (!insertAt) return null
    const conflict = arrayAt(withoutSource, insertAt.location)[insertAt.index] !== undefined
    return {steps: withoutSource, location: insertAt.location, index: insertAt.index, conflict}
}

/** Whether dropping sourceId onto targetId's given side would land on an already-occupied slot —
 * i.e. the caller should prompt for a MoveConflictChoice before calling moveStepRelativeTo. */
export function wouldConflict(rootSteps: TractStep[], sourceId: string, targetId: string, side: DropSide): boolean {
    return resolveDrop(rootSteps, sourceId, targetId, side)?.conflict ?? false
}

/** Moves sourceId to before/after targetId, per the canvas's drag-and-drop reordering. A "before"
 * drop always splices in directly (there's nothing to conflict with — see the DropSide comment at
 * the call site). An "after" drop into an occupied slot needs conflictChoice: "parallel" reuses
 * insertBlockAfter's fork/lane logic, "between" splices source directly ahead of target's existing
 * next step. */
export function moveStepRelativeTo(
    rootSteps: TractStep[], sourceId: string, targetId: string, side: DropSide, conflictChoice?: MoveConflictChoice,
): TractStep[] {
    const resolved = resolveDrop(rootSteps, sourceId, targetId, side)
    if (!resolved) return rootSteps
    const source = findStepById(rootSteps, sourceId)
    if (!source) return rootSteps

    if (resolved.conflict && conflictChoice === "parallel") {
        return insertBlockAfter(
            resolved.steps, resolved.location, resolved.index, source, collectAllStepIds(resolved.steps))
    }
    return insertStepAt(resolved.steps, resolved.location, resolved.index, source)
}
