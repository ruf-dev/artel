import {useCallback, useEffect, useRef, useState} from "react"

import {ChatEvent, PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"
import {applyEvent, ChatItem} from "@/pages/workbench/processes/chatReducer.ts"

export type ChatConnectionStatus = "connecting" | "open" | "reconnecting" | "closed"

const INITIAL_BACKOFF_MS = 1000
const MAX_BACKOFF_MS = 15000

function buildWsUrl(chatId: string): string {
    const proto = window.location.protocol === "https:" ? "wss:" : "ws:"
    return `${proto}//${window.location.host}/api/simple-chats/${chatId}/ws`
}

// Connects to a Simple Chat thread's WebSocket. Mirrors useChatSession.ts's shape
// (state, exponential-backoff reconnect starting at 1s doubling to a 15s cap,
// seq-based dedup via lastSeqRef, dispatch folding both directions into local
// items via chatReducer's applyEvent) but differs from it in three ways:
//  - keyed by `chatId` (a saved thread) rather than `vaultId` (a running bridge) —
//    switching chats means switching the socket, not just clearing local state.
//  - seeds `items` from `initialItems` (the thread's persisted history, fetched
//    separately via useSimpleChat/GetSimpleChat — see simpleChatMessages.ts)
//    instead of always starting empty, since a resumed thread has prior turns.
//  - stamps `model` onto every outgoing user_message with whatever the caller's
//    currently-selected model is (`currentModel`/`setModel`) — this is the
//    mechanism that makes mid-conversation model switching work end-to-end, since
//    the backend reads `event.model` per turn rather than fixing it at creation.
export function useSimpleChatSession(chatId: string | undefined, initialModel: string, initialItems: ChatItem[] = []) {
    const [items, setItems] = useState<ChatItem[]>(initialItems)
    const [status, setStatus] = useState<ChatConnectionStatus>("closed")
    const [model, setModel] = useState(initialModel)
    const wsRef = useRef<WebSocket | null>(null)
    const backoffRef = useRef(INITIAL_BACKOFF_MS)
    const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)
    const closedByUserRef = useRef(false)
    const lastSeqRef = useRef(0)
    const modelRef = useRef(initialModel)
    const initialItemsRef = useRef(initialItems)
    const initialModelRef = useRef(initialModel)

    // Kept fresh every render (not effect deps) so the connect-effect below can
    // read the latest seed values when it fires for a chatId change, without
    // re-running (and reconnecting) every time initialItems/initialModel's
    // identity changes from an unrelated refetch.
    initialItemsRef.current = initialItems
    initialModelRef.current = initialModel

    useEffect(() => {
        modelRef.current = model
    }, [model])

    useEffect(() => {
        if (!chatId) {
            setStatus("closed")
            return
        }

        closedByUserRef.current = false
        setItems(initialItemsRef.current)
        setStatus("connecting")
        setModel(initialModelRef.current)
        modelRef.current = initialModelRef.current
        lastSeqRef.current = 0

        function connect() {
            if (closedByUserRef.current) return

            const ws = new WebSocket(buildWsUrl(chatId!))
            wsRef.current = ws

            ws.onopen = () => {
                backoffRef.current = INITIAL_BACKOFF_MS
                setStatus("open")
            }
            ws.onmessage = event => {
                let parsed: ChatEvent
                try {
                    parsed = JSON.parse(event.data as string) as ChatEvent
                } catch {
                    return
                }
                // Deduplicate already-seen events based on seq — see
                // useChatSession.ts's identical logic for the full rationale.
                if (typeof parsed.seq === "number" && parsed.seq > 0 && parsed.seq <= lastSeqRef.current) {
                    return
                }
                if (typeof parsed.seq === "number" && parsed.seq > 0) {
                    lastSeqRef.current = parsed.seq
                }
                setItems(prev => applyEvent(prev, parsed))
            }
            ws.onclose = () => {
                if (closedByUserRef.current) return
                setStatus("reconnecting")
                const delay = backoffRef.current
                backoffRef.current = Math.min(backoffRef.current * 2, MAX_BACKOFF_MS)
                reconnectTimerRef.current = setTimeout(connect, delay)
            }
            ws.onerror = () => {
                ws.close()
            }
        }

        connect()

        return () => {
            closedByUserRef.current = true
            clearTimeout(reconnectTimerRef.current)
            wsRef.current?.close()
            wsRef.current = null
            setStatus("closed")
        }
    }, [chatId])

    const dispatch = useCallback((event: ChatEvent) => {
        const ws = wsRef.current
        if (!ws || ws.readyState !== WebSocket.OPEN) return
        ws.send(JSON.stringify(event))
        setItems(prev => applyEvent(prev, event))
    }, [])

    const sendMessage = useCallback((text: string) => {
        const id = crypto.randomUUID()
        dispatch({type: "user_message", text, id, model: modelRef.current})
    }, [dispatch])

    const sendPermissionDecision = useCallback((id: string, decision: PermissionDecision) => {
        dispatch({type: "permission_decision", id, decision})
    }, [dispatch])

    return {items, status, sendMessage, sendPermissionDecision, currentModel: model, setModel}
}
