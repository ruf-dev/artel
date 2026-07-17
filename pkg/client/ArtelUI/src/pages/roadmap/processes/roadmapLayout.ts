// Layered DAG layout turning a RoadmapGraph into node positions + connector edges for the canvas.
// Shape mirrors pages/tract-canvas/processes/tractCanvasLayout.ts's CanvasLayout on purpose (same
// {nodes, edges, width, height} contract, same constant values) but is its own implementation —
// colocation forbids importing across page-local subtrees, so the constants below are a
// deliberate redeclaration, not a drift risk (tract-canvas's are for a step tree, not a graph).
//
// Layering is longest-path-from-root via DFS: a node's layer is the length of the longest
// directed path root -> ... -> node found by the walk. Depends-on/Blocks pairs are naturally
// cyclic (A depends_on B and B blocks A both point away from A), so the walk tracks nodes
// currently on the recursion stack and refuses to recurse into one of them again ("back edge") —
// that node still gets a layer from whichever path reached it, it just doesn't get walked twice.

import {RoadmapEdge, RoadmapGraph, RoadmapNode} from "@/pages/roadmap/processes/roadmapGraph.ts"

export const NODE_WIDTH = 208
export const NODE_HEIGHT = 118
export const COLUMN_PITCH = 256
export const MARGIN_X = 48
export const MARGIN_Y = 220
const ROW_GAP = 32

export interface RoadmapLayoutNode {
    id: string
    x: number
    y: number
    node: RoadmapNode
}

export interface RoadmapLayoutEdge {
    id: string
    fromId: string
    toId: string
    relation: RoadmapEdge["relation"]
}

export interface RoadmapLayout {
    nodes: RoadmapLayoutNode[]
    edges: RoadmapLayoutEdge[]
    width: number
    height: number
}

function computeLayers(rootId: string, nodeIds: string[], edges: RoadmapEdge[]): Map<string, number> {
    const adjacency = new Map<string, string[]>()
    for (const id of nodeIds) adjacency.set(id, [])
    for (const e of edges) adjacency.get(e.fromId)?.push(e.toId)

    const layer = new Map<string, number>()
    const onStack = new Set<string>()

    function visit(id: string, depth: number) {
        const existing = layer.get(id)
        if (existing !== undefined && existing >= depth) return
        layer.set(id, depth)

        if (onStack.has(id)) return // back edge into a node already on this path — stop here
        onStack.add(id)
        for (const next of adjacency.get(id) ?? []) {
            visit(next, depth + 1)
        }
        onStack.delete(id)
    }

    visit(rootId, 0)
    for (const id of nodeIds) {
        if (!layer.has(id)) layer.set(id, 0) // unreached node (shouldn't normally happen) — park at root layer
    }

    return layer
}

export function layoutRoadmap(graph: RoadmapGraph): RoadmapLayout {
    const nodeIds = Object.keys(graph.nodes)
    if (nodeIds.length === 0) {
        return {nodes: [], edges: [], width: MARGIN_X * 2 + NODE_WIDTH, height: MARGIN_Y * 2 + NODE_HEIGHT}
    }

    const layerOf = computeLayers(graph.rootId, nodeIds, graph.edges)
    const byLayer = new Map<number, string[]>()
    for (const id of nodeIds) {
        const l = layerOf.get(id) ?? 0
        const arr = byLayer.get(l) ?? []
        arr.push(id)
        byLayer.set(l, arr)
    }

    const layers = Array.from(byLayer.keys())
    const maxLayer = Math.max(...layers)
    const maxRows = Math.max(...Array.from(byLayer.values()).map(a => a.length))

    const layoutNodes: RoadmapLayoutNode[] = []
    for (const [layer, ids] of byLayer) {
        const orderedIds = [...ids].sort()
        orderedIds.forEach((id, row) => {
            layoutNodes.push({
                id,
                x: MARGIN_X + layer * COLUMN_PITCH,
                y: MARGIN_Y + row * (NODE_HEIGHT + ROW_GAP),
                node: graph.nodes[id],
            })
        })
    }

    const layoutEdges: RoadmapLayoutEdge[] = graph.edges
        .filter(e => graph.nodes[e.fromId] && graph.nodes[e.toId])
        .map(e => ({id: e.id, fromId: e.fromId, toId: e.toId, relation: e.relation}))

    return {
        nodes: layoutNodes,
        edges: layoutEdges,
        width: MARGIN_X * 2 + NODE_WIDTH + maxLayer * COLUMN_PITCH,
        height: MARGIN_Y * 2 + NODE_HEIGHT + (maxRows - 1) * (NODE_HEIGHT + ROW_GAP),
    }
}
