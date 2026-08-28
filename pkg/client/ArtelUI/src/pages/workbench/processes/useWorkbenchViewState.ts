import {useEffect, useState} from "react"

import {type WorkbenchView} from "@/pages/workbench/processes/workbenchView.ts"

interface Params {
    exists: boolean
    status: string
    isLoading: boolean
    effectiveMode: string
    showSetup: boolean
    authComplete: boolean
    pendingAuthMode: string | undefined
}

// Owns WorkbenchPage's Chat/Terminal view toggle plus the layout flags derived from
// workbench/lifecycle state — split out purely to keep WorkbenchPage.tsx's render
// function under the max-lines-per-function lint limit.
export function useWorkbenchViewState(p: Params) {
    const [view, setView] = useState<WorkbenchView>("chat")

    const awaitingAuth = !p.authComplete && p.pendingAuthMode === "subscription_login"

    const dockerCentered = !p.exists || (p.status !== "running" && !p.showSetup)
    const genericCentered = p.isLoading || (p.effectiveMode === "docker" && dockerCentered)
    const terminalViewActive = p.effectiveMode === "docker" && p.status === "running" && view === "terminal"

    // Terminal login and the chat's own sign-in flow share the same in-container
    // credentials file, so there's no chat-side auth screen anymore — lock the Chat
    // toggle to Terminal while unauthenticated instead. Unlock-only: once auth
    // completes this stops firing, but view is never forced back to "chat" on its own.
    useEffect(() => {
        if (awaitingAuth && view === "chat") {
            setView("terminal")
        }
    }, [awaitingAuth, view])

    return {view, setView, awaitingAuth, genericCentered, terminalViewActive}
}
