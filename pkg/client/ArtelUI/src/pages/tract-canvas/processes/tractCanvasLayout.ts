// Pure layout algorithm turning a TractStep tree into node positions + connector edges for the
// canvas builder. Positions are recomputed on every render from the tree shape — nothing is
// persisted, so there is no backend change needed and no "drag a node" state to reconcile.
//
// Mapping to the mockup's visual language (trigger -> lookup -> [email, telegram] -> update):
//   - sequential steps lay out left-to-right, one column per step
//   - a "parallel" step is NOT drawn as its own node — its lanes fan out from the previous
//     node and fan back into whatever comes next, exactly like the mockup's fork/merge
//   - a "condition" step IS drawn as a node; its then/else branches fork above/below it as
//     dead-end sub-flows, and anything after the condition in the same list continues from the
//     condition node itself (branches don't feed data forward, per checkVisibility semantics)
//   - a "group" step is drawn as a single collapsed node — its nested steps are edited via the
//     inspector (reusing the tree editor), not expanded into the canvas
//
// Vertical placement is a two-pass measure-then-place: measureExtent() walks a branch/lane
// bottom-up to find how much room ABOVE and BELOW its own center row it actually needs
// (recursing into any condition/parallel nested inside it), then layoutFlow() places siblings
// using those real measurements instead of a guessed fixed offset. This is required because a
// fixed offset that merely decays with nesting depth can't know how tall a deeply-nested branch
// really is — it can allot less room than a branch needs and let it spill into a sibling's rows.

import {TractStep} from "@/processes/Tracts.ts"
import {Location, ROOT_LOCATION} from "@/processes/tractSteps.ts"

export const NODE_WIDTH = 208
export const NODE_HEIGHT = 118
export const COLUMN_PITCH = 256
export const MARGIN_X = 48
export const MARGIN_Y = 220

const ROW_GAP = 32

export type CanvasNodeKind = "trigger" | "action" | "condition" | "parallel" | "group" | "script"

export interface CanvasNode {
    id: string
    /** Pixel x/y, already margin-adjusted — ready to use directly as CSS left/top. */
    x: number
    y: number
    kind: CanvasNodeKind
    step?: TractStep
    /** Where this step lives in the tree (root list, or a parent's then/else/steps branch) and
     * its index within that array. Used for deletion and other node-specific lookups. */
    location: Location
    index: number
    /** Where a "+" press on this node should insert the new step. For ordinary nodes this is
     * just {location, index + 1}. For a node that is itself a lane inside a "parallel" step, this
     * instead points right after the whole parallel (its own location/index in ITS containing
     * array) — a lane has no "next" of its own, so "+" always means "join after every lane". */
    nextLocation: Location
    nextIndex: number
}

export interface CanvasEdge {
    id: string
    fromId: string
    toId: string
}

/** Bounding box (final pixel coords) for the soft-bordered container drawn around a non-empty
 * "parallel" step's lanes. Clicking it selects the parallel step itself (kind "parallel" node in
 * CanvasLayout.nodes) so its inspector (lane list, danger zone) opens like any other node. */
export interface ParallelBox {
    id: string
    x: number
    y: number
    width: number
    height: number
}

export interface CanvasLayout {
    nodes: CanvasNode[]
    edges: CanvasEdge[]
    parallelBoxes: ParallelBox[]
    width: number
    height: number
}

interface RawNode {
    id: string
    col: number
    row: number
    kind: CanvasNodeKind
    step?: TractStep
    location: Location
    index: number
    nextLocation: Location
    nextIndex: number
}

const PARALLEL_BOX_PADDING = 16

interface FlowResult {
    tipIds: string[]
    maxCol: number
    minRow: number
    maxRow: number
}

/** How far a branch/lane's rendered content reaches above and below its own center row —
 * everything needed to place it as a sibling without measuring its neighbors. */
interface Extent {
    above: number
    below: number
}

export const TRIGGER_NODE_ID = "trigger"

function addEdges(edges: CanvasEdge[], fromIds: string[], toId: string) {
    for (const fromId of fromIds) {
        edges.push({id: `${fromId}->${toId}`, fromId, toId})
    }
}

/** Bottom-up: how much vertical room does this steps list need above/below its own center row.
 * A plain sequential list needs only its own node height, regardless of length, since sequential
 * steps share one row and differ only in column. Only condition/parallel steps push the
 * requirement out further, by however much their own nested branches/lanes need. */
function measureExtent(steps: TractStep[]): Extent {
    let above = NODE_HEIGHT / 2
    let below = NODE_HEIGHT / 2

    for (const step of steps) {
        if (step.type === "condition") {
            const thenE = measureExtent(step.then ?? [])
            const elseE = measureExtent(step.else ?? [])
            above = Math.max(above, ROW_GAP / 2 + thenE.above + thenE.below)
            below = Math.max(below, ROW_GAP / 2 + elseE.above + elseE.below)
        } else if (step.type === "parallel") {
            const lanes = step.steps ?? []
            if (lanes.length > 0) {
                const total = lanes.reduce((sum, lane) => {
                    const e = measureExtent([lane])
                    return sum + e.above + e.below
                }, 0) + (lanes.length - 1) * ROW_GAP
                above = Math.max(above, total / 2)
                below = Math.max(below, total / 2)
            }
        }
    }

    return {above, below}
}

function layoutFlow(
    steps: TractStep[],
    location: Location,
    col: number,
    row: number,
    prevIds: string[],
    nodes: RawNode[],
    edges: CanvasEdge[],
    /** Overrides where THIS call's own nodes should send a "+" press. Set only when laying out a
     * single lane inside a parallel step (see the "parallel" branch below) — a lane's head node
     * has no "next" of its own within the lane array, so it borrows the enclosing parallel's own
     * (location, index + 1) instead of computing {location, index + 1} from its own array. */
    nextOverride?: {location: Location; index: number},
): FlowResult {
    let curCol = col
    let curPrevIds = prevIds
    let minRow = row
    let maxRow = row

    steps.forEach((step, index) => {
        const next = nextOverride ?? {location, index: index + 1}

        if (step.type === "condition") {
            const id = step.id
            nodes.push({
                id, col: curCol, row, kind: "condition", step, location, index,
                nextLocation: next.location, nextIndex: next.index,
            })
            addEdges(edges, curPrevIds, id)
            curCol++

            const thenSteps = step.then ?? []
            const elseSteps = step.else ?? []
            const thenE = measureExtent(thenSteps)
            const elseE = measureExtent(elseSteps)
            const thenLoc: Location = {parentId: id, branch: "then"}
            const elseLoc: Location = {parentId: id, branch: "else"}
            const thenRow = row - ROW_GAP / 2 - thenE.below
            const elseRow = row + ROW_GAP / 2 + elseE.above
            const thenRes = layoutFlow(thenSteps, thenLoc, curCol, thenRow, [id], nodes, edges)
            const elseRes = layoutFlow(elseSteps, elseLoc, curCol, elseRow, [id], nodes, edges)

            curCol = Math.max(curCol, thenRes.maxCol, elseRes.maxCol)
            minRow = Math.min(minRow, thenRes.minRow, elseRes.minRow)
            maxRow = Math.max(maxRow, thenRes.maxRow, elseRes.maxRow)
            curPrevIds = [id]
            return
        }

        if (step.type === "parallel") {
            const lanes = step.steps ?? []
            if (lanes.length === 0) {
                const id = step.id
                nodes.push({
                    id, col: curCol, row, kind: "parallel", step, location, index,
                    nextLocation: next.location, nextIndex: next.index,
                })
                addEdges(edges, curPrevIds, id)
                curPrevIds = [id]
                curCol++
                return
            }

            const laneLoc: Location = {parentId: step.id, branch: "steps"}
            const laneExtents = lanes.map(lane => measureExtent([lane]))
            const total = laneExtents.reduce((sum, e) => sum + e.above + e.below, 0) + (lanes.length - 1) * ROW_GAP
            let cursor = row - total / 2
            const laneResults = lanes.map((lane, i) => {
                const e = laneExtents[i]
                const laneRow = cursor + e.above
                cursor += e.above + e.below + ROW_GAP
                return layoutFlow([lane], laneLoc, curCol, laneRow, curPrevIds, nodes, edges, next)
            })

            curCol = Math.max(curCol, ...laneResults.map(r => r.maxCol))
            minRow = Math.min(minRow, ...laneResults.map(r => r.minRow))
            maxRow = Math.max(maxRow, ...laneResults.map(r => r.maxRow))
            curPrevIds = laneResults.flatMap(r => r.tipIds)

            // A selectable stand-in for the parallel step itself, so its box (drawn separately,
            // see layoutTract's parallelBoxes) can open the same inspector as any other node.
            // Not part of the visible flow: no column advance, no edges of its own.
            nodes.push({
                id: step.id, col, row, kind: "parallel", step, location, index,
                nextLocation: next.location, nextIndex: next.index,
            })
            return
        }

        // action, script, or group: one node each (a group's own nested steps are edited in
        // the inspector, not expanded onto the canvas)
        const id = step.id
        let kind: CanvasNodeKind = "action"
        if (step.type === "group") kind = "group"
        if (step.type === "script") kind = "script"
        nodes.push({
            id, col: curCol, row, kind, step, location, index,
            nextLocation: next.location, nextIndex: next.index,
        })
        addEdges(edges, curPrevIds, id)
        curPrevIds = [id]
        curCol++
    })

    return {tipIds: curPrevIds, maxCol: curCol, minRow, maxRow}
}

export function layoutTract(rootSteps: TractStep[]): CanvasLayout {
    const rawNodes: RawNode[] = [{
        id: TRIGGER_NODE_ID, col: 0, row: 0, kind: "trigger", location: ROOT_LOCATION, index: -1,
        nextLocation: ROOT_LOCATION, nextIndex: 0,
    }]
    const edges: CanvasEdge[] = []

    const result = layoutFlow(rootSteps, ROOT_LOCATION, 1, 0, [TRIGGER_NODE_ID], rawNodes, edges)

    const rowShift = MARGIN_Y - result.minRow
    const nodes: CanvasNode[] = rawNodes.map(n => ({
        id: n.id,
        x: MARGIN_X + n.col * COLUMN_PITCH,
        y: rowShift + n.row,
        kind: n.kind,
        step: n.step,
        location: n.location,
        index: n.index,
        nextLocation: n.nextLocation,
        nextIndex: n.nextIndex,
    }))

    const parallelBoxes: ParallelBox[] = []
    for (const n of nodes) {
        if (n.kind !== "parallel") continue
        const laneNodes = nodes.filter(m => m.location.parentId === n.id && m.location.branch === "steps")
        if (laneNodes.length === 0) continue // empty parallel — drawn as its own plain node, no box
        const minX = Math.min(...laneNodes.map(m => m.x))
        const minY = Math.min(...laneNodes.map(m => m.y))
        const maxY = Math.max(...laneNodes.map(m => m.y + NODE_HEIGHT))
        parallelBoxes.push({
            id: n.id,
            x: minX - PARALLEL_BOX_PADDING,
            y: minY - PARALLEL_BOX_PADDING,
            width: NODE_WIDTH + PARALLEL_BOX_PADDING * 2,
            height: (maxY - minY) + PARALLEL_BOX_PADDING * 2,
        })
    }

    const width = MARGIN_X * 2 + NODE_WIDTH + result.maxCol * COLUMN_PITCH
    const height = MARGIN_Y * 2 + NODE_HEIGHT + (result.maxRow - result.minRow)

    return {nodes, edges, parallelBoxes, width, height}
}
