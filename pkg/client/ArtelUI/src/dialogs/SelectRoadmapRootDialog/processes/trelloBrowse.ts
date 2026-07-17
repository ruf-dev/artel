// Thin wrappers around the generic MoM ExecuteMomTool RPC for the trello-specific "browse" tools
// this dialog needs (list_boards/list_lists/list_cards). Colocated here rather than as a global
// src/processes/ file since only this dialog's screens use them.

import {useMcpKeys} from "@/app/hooks/McpKeys.ts"

export interface TrelloBoardLite {
    id: string
    name: string
    desc: string
    closed: boolean
    url: string
}

export interface TrelloListLite {
    id: string
    name: string
    closed: boolean
    idBoard: string
}

export interface TrelloCardLite {
    id: string
    name: string
    desc: string
    idList: string
    idBoard: string
    closed: boolean
    url: string
    labels?: unknown[]
}

function execute(trackerId: string, toolName: string, params: Record<string, unknown>): Promise<string> {
    return useMcpKeys.getState().executeMomTool("trello", toolName, trackerId, params)
}

export function fetchBoards(trackerId: string): Promise<TrelloBoardLite[]> {
    return execute(trackerId, "list_boards", {}).then(raw => JSON.parse(raw) as TrelloBoardLite[])
}

export function fetchLists(trackerId: string, boardId: string): Promise<TrelloListLite[]> {
    return execute(trackerId, "list_lists", {board_id: boardId}).then(raw => JSON.parse(raw) as TrelloListLite[])
}

export function fetchCards(trackerId: string, boardId: string): Promise<TrelloCardLite[]> {
    return execute(trackerId, "list_cards", {board_id: boardId}).then(raw => JSON.parse(raw) as TrelloCardLite[])
}
