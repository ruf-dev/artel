import {useWorkbenchHistory, WorkbenchHistory} from "@/pages/workbench/processes/useWorkbenchHistory.ts"
import type {WorkbenchMode} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"
import type {SimpleChatController} from "@/pages/workbench/processes/useSimpleChatController.ts"
import type {WorkbenchVaultFiles} from "@/pages/workbench/processes/useWorkbenchVaultFiles.ts"
import type {WorkbenchAttachments} from "@/pages/workbench/processes/useWorkbenchAttachments.ts"
import type {NoteItem} from "@/app/api/artel/notes.pb.ts"

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
    vaultName: string
    modeControls: ModeControls
    startNewDockerChat: () => void
    controller: SimpleChatController
    vaultFiles: WorkbenchVaultFiles
    attachments: WorkbenchAttachments
}

// Everything the Vault tab pane needs, shaped once here so WorkbenchPage threads a
// single object down (matching WorkbenchHistory). Also serves as VaultPane's props.
export interface VaultPaneBundle {
    vaultName: string
    folders: string[]
    notes: NoteItem[]
    isLoading: boolean
    attachedPaths: string[]
    onToggleAttach: (path: string) => void
}

export interface WorkbenchSidebarState {
    show: boolean
    history: WorkbenchHistory
    vault: VaultPaneBundle
}

// Bundles everything WorkbenchPage needs to render the unified WorkbenchSidebar
// (the flat history list, the Vault pane bundle, and the show/hide predicate) —
// split out purely to keep WorkbenchPage.tsx's render function under the
// max-lines-per-function lint limit, same rationale as
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
        vault: {
            vaultName: p.vaultName,
            folders: p.vaultFiles.folders,
            notes: p.vaultFiles.notes,
            isLoading: p.vaultFiles.isLoading,
            attachedPaths: p.attachments.paths,
            onToggleAttach: p.attachments.toggle,
        },
    }
}
