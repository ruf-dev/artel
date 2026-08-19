import {useCallback, useEffect, useRef, useState} from "react"

import {ChatEvent, PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"
import {applyEvent, ChatItem} from "@/pages/workbench/processes/chatReducer.ts"

export type ChatConnectionStatus = "connecting" | "open" | "reconnecting" | "closed"

const INITIAL_BACKOFF_MS = 1000
const MAX_BACKOFF_MS = 15000

function buildWsUrl(vaultId: string): string {
    const proto = window.location.protocol === "https:" ? "wss:" : "ws:"
    return `${proto}//${window.location.host}/api/vaults/workbench/${vaultId}/terminal/`
}

// Connects to the workbench chat WebSocket (same route the old ttyd terminal iframe used —
// the backend track keeps the route/port, only the wire protocol changes from raw terminal
// bytes to ChatEvent JSON), reconnecting with exponential backoff on drop. Exposes the
// folded item list (see chatReducer.ts) plus helpers for the three consumer -> bridge event
// types; each helper both sends over the wire and applies the event to local state
// optimistically, since the bridge never echoes these back.
export function useChatSession(vaultId: string | undefined) {
    const [items, setItems] = useState<ChatItem[]>([])
    const [status, setStatus] = useState<ChatConnectionStatus>("connecting")
    const [authComplete, setAuthComplete] = useState(false)
    const wsRef = useRef<WebSocket | null>(null)
    const backoffRef = useRef(INITIAL_BACKOFF_MS)
    const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)
    const closedByUserRef = useRef(false)

    useEffect(() => {
        if (!vaultId) return

        closedByUserRef.current = false
        setItems([])
        setStatus("connecting")
        setAuthComplete(false)

        function connect() {
            if (closedByUserRef.current) return

            const ws = new WebSocket(buildWsUrl(vaultId!))
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
                if (parsed.type === "auth_complete") setAuthComplete(true)
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
    }, [vaultId])

    const dispatch = useCallback((event: ChatEvent) => {
        const ws = wsRef.current
        if (!ws || ws.readyState !== WebSocket.OPEN) return
        ws.send(JSON.stringify(event))
        setItems(prev => applyEvent(prev, event))
    }, [])

    const sendMessage = useCallback((text: string) => {
        dispatch({type: "user_message", text})
    }, [dispatch])

    const sendPermissionDecision = useCallback((id: string, decision: PermissionDecision) => {
        dispatch({type: "permission_decision", id, decision})
    }, [dispatch])

    const sendAuthCode = useCallback((code: string) => {
        dispatch({type: "auth_code_submit", code})
    }, [dispatch])

    return {items, status, authComplete, sendMessage, sendPermissionDecision, sendAuthCode}
}
