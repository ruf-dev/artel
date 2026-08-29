import {Dispatch, MutableRefObject, SetStateAction, useCallback, useEffect, useRef, useState} from "react"

import {ChatEvent, PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"
import {applyEvent, ChatItem, truncateAfterUserMessage} from "@/pages/workbench/processes/chatReducer.ts"

export type ChatConnectionStatus = "connecting" | "open" | "reconnecting" | "closed"

const INITIAL_BACKOFF_MS = 1000
const MAX_BACKOFF_MS = 15000

function buildWsUrl(chatId: string): string {
    const proto = window.location.protocol === "https:" ? "wss:" : "ws:"
    return `${proto}//${window.location.host}/api/simple-chats/${chatId}/ws`
}

// Parses one incoming WS frame, deduplicates it against lastSeqRef, and folds it
// into items - pulled out of the connect effect below (rather than inlined in
// ws.onmessage) purely to keep the hook under the max-lines-per-function budget.
// See useChatSession.ts's identical inline logic for the full seq-dedup rationale.
function handleIncomingFrame(
    raw: string,
    lastSeqRef: MutableRefObject<number>,
    setPendingTurn: Dispatch<SetStateAction<boolean>>,
    setItems: Dispatch<SetStateAction<ChatItem[]>>,
) {
    let parsed: ChatEvent
    try {
        parsed = JSON.parse(raw) as ChatEvent
    } catch {
        return
    }
    if (typeof parsed.seq === "number" && parsed.seq > 0 && parsed.seq <= lastSeqRef.current) {
        return
    }
    if (typeof parsed.seq === "number" && parsed.seq > 0) {
        lastSeqRef.current = parsed.seq
    }
    if (parsed.type === "turn_done" || parsed.type === "error") {
        setPendingTurn(false)
    }
    setItems(prev => applyEvent(prev, parsed))
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
    const [pendingTurn, setPendingTurn] = useState(false)
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
        setPendingTurn(false)
        lastSeqRef.current = 0

        function connect() {
            if (closedByUserRef.current) return

            const ws = new WebSocket(buildWsUrl(chatId!))
            wsRef.current = ws

            ws.onopen = () => {
                backoffRef.current = INITIAL_BACKOFF_MS
                setStatus("open")
            }
            ws.onmessage = event => handleIncomingFrame(event.data as string, lastSeqRef, setPendingTurn, setItems)
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

    // Seeds items from initialItems once the persisted-history fetch resolves,
    // separately from the connect effect above: on a reload/direct-nav to an
    // already-known chatId, the connect effect fires before GetSimpleChat resolves,
    // seeding from an empty initialItems - chatId doesn't change once the real
    // history arrives, so the connect effect never re-runs to pick it up. The
    // prev.length === 0 guard (read live via the functional setItems updater)
    // means this only ever seeds a still-untouched list, so it can't clobber a
    // live WS event that arrived before the query resolved.
    useEffect(() => {
        if (!chatId || initialItems.length === 0) return
        setItems(prev => prev.length === 0 ? initialItems : prev)
    }, [chatId, initialItems])

    const dispatch = useCallback((event: ChatEvent) => {
        const ws = wsRef.current
        if (!ws || ws.readyState !== WebSocket.OPEN) return
        ws.send(JSON.stringify(event))
        setItems(prev => applyEvent(prev, event))
    }, [])

    const sendMessage = useCallback((text: string) => {
        const id = crypto.randomUUID()
        dispatch({type: "user_message", text, id, model: modelRef.current})
        setPendingTurn(true)
    }, [dispatch])

    // Resends a failed/dangling user_message reusing its id - see useChatSession.ts's
    // resendMessage doc comment for the full rationale. Also stamps the current
    // model, mirroring sendMessage.
    const resendMessage = useCallback((id: string, text: string) => {
        setItems(prev => truncateAfterUserMessage(prev, id))
        dispatch({type: "user_message", text, id, model: modelRef.current})
        setPendingTurn(true)
    }, [dispatch])

    const sendPermissionDecision = useCallback((id: string, decision: PermissionDecision) => {
        dispatch({type: "permission_decision", id, decision})
        setPendingTurn(true)
    }, [dispatch])

    return {
        items, status, sendMessage, resendMessage, sendPermissionDecision, currentModel: model, setModel, pendingTurn,
    }
}
