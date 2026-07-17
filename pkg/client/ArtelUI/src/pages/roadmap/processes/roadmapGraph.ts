// Graph state shape + the async "expand a card" step for the roadmap canvas (pages/roadmap).
// Kept as plain data + pure/async functions rather than a Zustand store: the graph is entirely
// page-local (nothing outside pages/roadmap needs it), and a global app/hooks store would have to
// reach into this colocated file from outside its owning page, which the colocation rule in
// pkg/client/ArtelUI/CLAUDE.md forbids ("don't reach into another page's local folder from
// outside it"). RoadmapPage.tsx owns this as useState and calls expandNode() directly.

import {ParsedTaskLink, parseTaskLinks, RoadmapRelation} from "@/pages/roadmap/processes/taskLinkParser.ts"

export interface RoadmapNode {
    id: string
    name: string
    desc: string
    url: string
    /** Trello shortLink for this card — always present because every node is loaded via get_card,
     * which returns it; used to resolve other cards' parsed links against this node. */
    shortLink: string
    idBoard: string
    idList: string
    closed: boolean
    expanded: boolean
    /** Every relation line found in this card's comments, resolved or not — RoadmapInspector
     * renders each as either a "jump to node" (if `resolvedNodeId` matches a loaded node) or an
     * external-link chip otherwise. */
    links: ParsedTaskLink[]
}

export interface RoadmapEdge {
    id: string
    fromId: string
    toId: string
    relation: RoadmapRelation
}

export interface RoadmapGraph {
    rootId: string
    nodes: Record<string, RoadmapNode>
    edges: RoadmapEdge[]
}

interface TrelloCardDetail {
    id: string
    name: string
    desc: string
    idList: string
    idBoard: string
    closed: boolean
    url: string
    shortLink: string
}

interface TrelloCommentAction {
    id: string
    data?: {text?: string}
}

/** Calls one MoM tool (already bound to mcpName "trello" + a tracker's externalConnectionId) and
 * returns its raw JSON result string, exactly what McpKeysAPI.ExecuteMomTool hands back. */
export type MomToolExecutor = (toolName: string, params: Record<string, unknown>) => Promise<string>

export function createRoadmapGraph(rootId: string): RoadmapGraph {
    return {rootId, nodes: {}, edges: []}
}

function cardToNode(card: TrelloCardDetail, links: ParsedTaskLink[], expanded: boolean): RoadmapNode {
    return {
        id: card.id,
        name: card.name,
        desc: card.desc,
        url: card.url,
        shortLink: card.shortLink,
        idBoard: card.idBoard,
        idList: card.idList,
        closed: card.closed,
        expanded,
        links,
    }
}

function findNodeByShortLink(nodes: Record<string, RoadmapNode>, shortLink: string): RoadmapNode | undefined {
    return Object.values(nodes).find(n => n.shortLink === shortLink || n.id === shortLink)
}

/** Loads (if needed) + expands a single card: fetches its detail (skipped if already loaded),
 * always re-fetches + re-parses its comments so re-running this after AddTaskLinkDialog adds a
 * new comment picks up the new link. Any parsed link whose target is already loaded becomes an
 * edge immediately; everything else stays an unresolved link on the node for the inspector to
 * show as an external chip until/unless the user expands that target too. */
export async function expandNode(graph: RoadmapGraph, cardId: string, execute: MomToolExecutor): Promise<RoadmapGraph> {
    let node = graph.nodes[cardId]

    if (!node) {
        const rawCard = await execute("get_card", {id: cardId})
        const card = JSON.parse(rawCard) as TrelloCardDetail
        node = cardToNode(card, [], false)
    }

    const rawComments = await execute("list_card_comments", {id: cardId})
    const comments = JSON.parse(rawComments) as TrelloCommentAction[]
    const links = comments.flatMap(c => parseTaskLinks(c.data?.text ?? ""))

    const nodes: Record<string, RoadmapNode> = {...graph.nodes, [node.id]: {...node, expanded: true, links}}
    const edges = [...graph.edges]

    for (const link of links) {
        if (!link.shortLink) continue
        const target = findNodeByShortLink(nodes, link.shortLink)
        if (!target || target.id === node.id) continue

        const edgeId = `${node.id}->${target.id}->${link.relation}`
        if (!edges.some(e => e.id === edgeId)) {
            edges.push({id: edgeId, fromId: node.id, toId: target.id, relation: link.relation})
        }
    }

    return {...graph, nodes, edges}
}
