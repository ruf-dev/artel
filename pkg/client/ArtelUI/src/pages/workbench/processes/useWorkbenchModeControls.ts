import {useState} from "react"

// Simple Chat has no dependency on the Docker workbench at all (no "exists"/
// "status" concept of its own), so it's tracked as an independent top-level mode
// rather than folded into the Docker status state machine in useWorkbenchLifecycle.
// `modeChoice` is the user's explicit choice; `null` means "no explicit choice
// yet — fall back to whatever's implied by whether a Docker workbench already
// exists" (see `effectiveMode` below), which is what keeps an existing Docker
// user's landing experience byte-for-byte unchanged (no picker ever shown once
// `exists` is true). "picking" is a distinct explicit choice (not the same as
// `null`), kept for a future re-entry point back to the picker.
export type WorkbenchMode = "docker" | "simple-chat"
export type WorkbenchModeChoice = WorkbenchMode | "picking" | null

interface Params {
    exists: boolean
    handleCreateDocker: () => void
}

// Bundles the top-level Docker-vs-Simple-Chat mode state and its transition
// handlers — split out purely to keep WorkbenchPage.tsx's render function under
// the max-lines-per-function lint limit, same rationale as
// useWorkbenchPanelControls.ts.
export function useWorkbenchModeControls({exists, handleCreateDocker}: Params) {
    const [modeChoice, setModeChoice] = useState<WorkbenchModeChoice>(null)
    const [simpleChatId, setSimpleChatId] = useState<string | undefined>(undefined)
    const [simpleHistoryOpen, setSimpleHistoryOpen] = useState(false)

    const effectiveMode: WorkbenchMode | "picking" = modeChoice ?? (exists ? "docker" : "picking")

    function handlePickDocker() {
        setModeChoice(null)
        if (!exists) handleCreateDocker()
    }

    function handleSimpleChatCreated(chatId: string) {
        setModeChoice("simple-chat")
        setSimpleChatId(chatId)
    }

    function toggleSimpleHistory() {
        setSimpleHistoryOpen(open => !open)
    }

    function closeSimpleHistory() {
        if (simpleChatId) setSimpleHistoryOpen(false)
    }

    return {
        effectiveMode,
        simpleChatId,
        setSimpleChatId,
        simpleHistoryOpen,
        setSimpleHistoryOpen,
        handlePickDocker,
        handleSimpleChatCreated,
        toggleSimpleHistory,
        closeSimpleHistory,
    }
}
