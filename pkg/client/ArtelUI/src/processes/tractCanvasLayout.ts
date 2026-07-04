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

import {TractStep} from "@/processes/Tracts.ts"
import {Location, ROOT_LOCATION} from "@/processes/tractSteps.ts"

export const NODE_WIDTH = 208
export const NODE_HEIGHT = 118
export const COLUMN_PITCH = 256
export const MARGIN_X = 48
export const MARGIN_Y = 220

const BASE_SPREAD = 96
const MIN_SPREAD = 50

export type CanvasNodeKind = "trigger" | "action" | "condition" | "parallel" | "group"

export interface CanvasNode {
    id: string
    /** Pixel x/y, already margin-adjusted — ready to use directly as CSS left/top. */
    x: number
    y: number
    kind: CanvasNodeKind
    step?: TractStep
    /** Where this step lives in the tree (root list, or a parent's then/else/steps branch) and
     * its index within that array — enough to insert a new step right after it. */
    location: Location
    index: number
}

export interface CanvasEdge {
    id: string
    fromId: string
    toId: string
}

export interface CanvasLayout {
    nodes: CanvasNode[]
    edges: CanvasEdge[]
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
}

interface FlowResult {
    tipIds: string[]
    maxCol: number
    minRow: number
    maxRow: number
}

export const TRIGGER_NODE_ID = "trigger"

function addEdges(edges: CanvasEdge[], fromIds: string[], toId: string) {
    for (const fromId of fromIds) {
        edges.push({id: `${fromId}->${toId}`, fromId, toId})
    }
}

function layoutFlow(
    steps: TractStep[],
    location: Location,
    col: number,
    row: number,
    spread: number,
    prevIds: string[],
    nodes: RawNode[],
    edges: CanvasEdge[],
): FlowResult {
    let curCol = col
    let curPrevIds = prevIds
    let minRow = row
    let maxRow = row

    steps.forEach((step, index) => {
        if (step.type === "condition") {
            const id = step.id
            nodes.push({id, col: curCol, row, kind: "condition", step, location, index})
            addEdges(edges, curPrevIds, id)
            curCol++

            const childSpread = Math.max(spread / 2, MIN_SPREAD)
            const thenLoc: Location = {parentId: id, branch: "then"}
            const elseLoc: Location = {parentId: id, branch: "else"}
            const thenRes = layoutFlow(step.then ?? [], thenLoc, curCol, row - spread, childSpread, [id], nodes, edges)
            const elseRes = layoutFlow(step.else ?? [], elseLoc, curCol, row + spread, childSpread, [id], nodes, edges)

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
                nodes.push({id, col: curCol, row, kind: "parallel", step, location, index})
                addEdges(edges, curPrevIds, id)
                curPrevIds = [id]
                curCol++
                return
            }

            const childSpread = Math.max(spread / 2, MIN_SPREAD)
            const startOffset = -(lanes.length - 1) / 2
            const laneLoc: Location = {parentId: step.id, branch: "steps"}
            const laneResults = lanes.map((lane, i) =>
                layoutFlow([lane], laneLoc, curCol, row + (startOffset + i) * spread, childSpread, curPrevIds, nodes, edges),
            )

            curCol = Math.max(curCol, ...laneResults.map(r => r.maxCol))
            minRow = Math.min(minRow, ...laneResults.map(r => r.minRow))
            maxRow = Math.max(maxRow, ...laneResults.map(r => r.maxRow))
            curPrevIds = laneResults.flatMap(r => r.tipIds)
            return
        }

        // action or group: one node each (a group's own nested steps are edited in the
        // inspector, not expanded onto the canvas)
        const id = step.id
        nodes.push({id, col: curCol, row, kind: step.type === "group" ? "group" : "action", step, location, index})
        addEdges(edges, curPrevIds, id)
        curPrevIds = [id]
        curCol++
    })

    return {tipIds: curPrevIds, maxCol: curCol, minRow, maxRow}
}

export function layoutTract(rootSteps: TractStep[]): CanvasLayout {
    const rawNodes: RawNode[] = [{id: TRIGGER_NODE_ID, col: 0, row: 0, kind: "trigger", location: ROOT_LOCATION, index: -1}]
    const edges: CanvasEdge[] = []

    const result = layoutFlow(rootSteps, ROOT_LOCATION, 1, 0, BASE_SPREAD, [TRIGGER_NODE_ID], rawNodes, edges)

    const rowShift = MARGIN_Y - result.minRow
    const nodes: CanvasNode[] = rawNodes.map(n => ({
        id: n.id,
        x: MARGIN_X + n.col * COLUMN_PITCH,
        y: rowShift + n.row,
        kind: n.kind,
        step: n.step,
        location: n.location,
        index: n.index,
    }))

    const width = MARGIN_X * 2 + NODE_WIDTH + result.maxCol * COLUMN_PITCH
    const height = MARGIN_Y * 2 + NODE_HEIGHT + (result.maxRow - result.minRow)

    return {nodes, edges, width, height}
}
