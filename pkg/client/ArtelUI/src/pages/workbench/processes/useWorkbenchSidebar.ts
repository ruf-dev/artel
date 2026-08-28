import {useWorkbenchHistory, WorkbenchHistory} from "@/pages/workbench/processes/useWorkbenchHistory.ts"
import type {WorkbenchMode} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"
import type {SimpleChatController} from "@/pages/workbench/processes/useSimpleChatController.ts"

interface ModeControls {
    simpleChatId?: string
    setSimpleChatId: (id: string | undefined) => void
}

interface Params {
    effectiveMode: WorkbenchMode | "picking"
    isLoading: boolean
    exists: boolean
    status: string
    vaultId?: string
    modeControls: ModeControls
    startNewDockerChat: () => void
    controller: SimpleChatController
}

export interface WorkbenchSidebarState {
    show: boolean
    history: WorkbenchHistory
    showCloseButton: boolean
}

// Bundles everything WorkbenchPage needs to render the unified WorkbenchSidebar
// (the flat history list, the exit-to-home control, and the show/hide
// predicate) — split out purely to keep WorkbenchPage.tsx's render function under
// the max-lines-per-function lint limit, same rationale as
// useWorkbenchPanelControls.ts / useWorkbenchViewState.ts.
export function useWorkbenchSidebar(p: Params): WorkbenchSidebarState {
    const isApi = p.effectiveMode === "api"
    const mode: WorkbenchMode = p.effectiveMode === "docker" ? "docker" : "api"

    const history = useWorkbenchHistory({
        mode,
        vaultId: p.vaultId ?? "",
        activeApiChatId: p.modeControls.simpleChatId,
        onSelectApiChat: p.modeControls.setSimpleChatId,
        onNewChat: mode === "docker" ? p.startNewDockerChat : p.controller.handleNewChat,
    })

    const show = !p.isLoading && !!p.vaultId && (
        isApi || (p.effectiveMode === "docker" && p.exists && p.status === "running")
    )

    return {
        show,
        history,
        showCloseButton: isApi,
    }
}
