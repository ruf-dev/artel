import {useParams} from "react-router-dom"
import {Loader} from "@vervstack/chures"

import cls from "@/pages/workbench/WorkbenchPage.module.css"
import {cn} from "@/app/utils/cn.ts"
import {useVaults} from "@/app/hooks/Vaults.ts"
import {useSimpleChats} from "@/app/hooks/SimpleChat.ts"
import {
    useWorkbench,
    useWorkbenchTerminalTabs,
    useWorkbenchTerminalTabMutations,
} from "@/app/hooks/Workbench.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {useDocumentTitle} from "@/app/hooks/useDocumentTitle.ts"
import PickAuthModeScreen from "@/pages/workbench/components/PickAuthModeScreen/PickAuthModeScreen.tsx"
import PickWorkbenchModeScreen from "@/pages/workbench/components/PickWorkbenchModeScreen/PickWorkbenchModeScreen.tsx"
import WorkbenchPanels from "@/pages/workbench/components/WorkbenchPanels/WorkbenchPanels.tsx"
import WorkbenchSidebar from "@/pages/workbench/components/WorkbenchSidebar/WorkbenchSidebar.tsx"
import WorkbenchTopbar from "@/pages/workbench/components/WorkbenchTopbar/WorkbenchTopbar.tsx"
import WorkbenchTweaksPanel from "@/pages/workbench/components/WorkbenchTweaksPanel/WorkbenchTweaksPanel.tsx"
import {useChatSession} from "@/pages/workbench/processes/useChatSession.ts"
import {
    toSimpleChatSessionBundle,
    useSimpleChatController,
} from "@/pages/workbench/processes/useSimpleChatController.ts"
import {useWorkbenchLifecycle} from "@/pages/workbench/processes/useWorkbenchLifecycle.ts"
import {useWorkbenchPanelControls} from "@/pages/workbench/processes/useWorkbenchPanelControls.ts"
import {useWorkbenchModeControls} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"
import {useWorkbenchContext} from "@/pages/workbench/processes/workbenchContext.ts"
import {useWorkbenchAttachments} from "@/pages/workbench/processes/useWorkbenchAttachments.ts"
import {useWorkbenchVaultFiles} from "@/pages/workbench/processes/useWorkbenchVaultFiles.ts"
import {useWorkbenchSidebar} from "@/pages/workbench/processes/useWorkbenchSidebar.ts"
import {useWorkbenchNavDrawer} from "@/pages/workbench/processes/useWorkbenchNavDrawer.ts"
import {useWorkbenchViewState} from "@/pages/workbench/processes/useWorkbenchViewState.ts"

// eslint-disable-next-line max-lines-per-function -- root page component with complex setup
export default function WorkbenchPage() {
    const {vaultId} = useParams()
    const {vaults} = useVaults()
    const {exists, status, isLoading, pendingTerminalAuthLink} = useWorkbench(vaultId)
    const {chats: simpleChats, isLoading: simpleChatsLoading} = useSimpleChats(vaultId)
    const lifecycle = useWorkbenchLifecycle(vaultId)
    const {tabs} = useWorkbenchTerminalTabs(vaultId, status === "running")
    const {create: createTab, select: selectTab, close: closeTab} = useWorkbenchTerminalTabMutations(vaultId)
    const bakeError = useBakeError()
    const {handleSelectTab, handleCreateTab, handleCloseTab} =
        useWorkbenchPanelControls({selectTab, createTab, closeTab, bakeError})
    const modeControls = useWorkbenchModeControls({
        exists, handleCreateDocker: lifecycle.handleCreate, simpleChats, simpleChatsLoading,
    })
    const {effectiveMode} = modeControls
    // Connect Docker's bridge only while Docker mode is active; the api chat owns its own socket.
    const chatSession = useChatSession(effectiveMode === "docker" && status === "running" ? vaultId : undefined)

    const simpleChatController = useSimpleChatController({
        chatId: modeControls.simpleChatId,
        vaultId: vaultId ?? "",
        active: effectiveMode === "api",
        onChatCreated: modeControls.setSimpleChatId,
        bakeError,
    })

    const vaultName = vaults.find(v => v.id === vaultId)?.name ?? "Vault"

    const ctx = useWorkbenchContext()
    const attachments = useWorkbenchAttachments()
    const vaultFiles = useWorkbenchVaultFiles(vaultId)

    const {view, setView, awaitingAuth, genericCentered, terminalViewActive} = useWorkbenchViewState({
        exists, status, isLoading, effectiveMode, showSetup: lifecycle.showSetup,
        authComplete: chatSession.authComplete, pendingAuthMode: lifecycle.pendingAuthMode,
    })

    // The Vault pane reads vaultFiles + attachments: clicking a file row there calls
    // attachments.toggle, and attached paths render active in the tree.
    const sidebar = useWorkbenchSidebar({
        effectiveMode, isLoading, exists, status, vaultId, vaultName, modeControls, vaultFiles, attachments,
        startNewDockerChat: chatSession.startNewChat, controller: simpleChatController,
    })

    const {navOpen, toggleNav, sidebarHistory} = useWorkbenchNavDrawer(sidebar.history)

    useDocumentTitle(vaultName)

    return (
        <div className={cn(cls.WorkbenchPageContainer, sidebar.show && cls.WithSidebar,
            sidebar.show && !navOpen && cls.NavCollapsed, terminalViewActive && cls.TerminalViewActive)}>
            {/* MainColumn is rendered before the sidebar so the sidebar's mobile fixed
                drawer paints above the topbar's position:relative model-switcher trigger.
                Desktop column order is restored via grid-column placement (not `order`,
                which would also reverse paint order) — see WorkbenchPage.module.css. */}
            <div className={cls.MainColumn}>
                <div className={cls.ContentStack}>
                    <div className={cn(cls.Body, genericCentered && cls.BodyCentered,
                        effectiveMode === "picking" && cls.BodyPicking)}>
                        {isLoading && (
                            <Loader variant="arcs" size="sm" color="var(--coral)"/>
                        )}
                        {!isLoading && effectiveMode === "picking" && vaultId && (
                            <PickWorkbenchModeScreen
                                vaultId={vaultId} onStartDocker={modeControls.handlePickDocker}
                                startingDocker={lifecycle.creating}
                                onSimpleChatCreated={modeControls.handleSimpleChatCreated}
                            />
                        )}
                        {!isLoading && vaultId && (
                            <WorkbenchPanels
                                effectiveMode={effectiveMode} exists={exists} status={status} vaultId={vaultId}
                                view={view} awaitingAuth={awaitingAuth} chatSession={chatSession} tabs={tabs}
                                pendingTerminalAuthLink={pendingTerminalAuthLink} onSelectTab={handleSelectTab}
                                onCreateTab={handleCreateTab} onCloseTab={handleCloseTab}
                                simpleChatId={modeControls.simpleChatId} ctx={ctx}
                                attachments={attachments}
                                simpleChatSession={toSimpleChatSessionBundle(simpleChatController)}
                            />
                        )}
                        {!isLoading && effectiveMode === "docker" && exists && status !== "running"
                            && lifecycle.showSetup && vaultId && (
                            <PickAuthModeScreen
                                onStart={lifecycle.handleStartWorkbench}
                                starting={lifecycle.startingWorkbench}
                            />
                        )}
                    </div>
                    {!isLoading && effectiveMode !== "picking" && (
                        <WorkbenchTopbar
                            effectiveMode={effectiveMode} exists={exists} status={status} vaultId={vaultId}
                            view={view} onViewChange={setView} chatLocked={awaitingAuth}
                            onStart={lifecycle.handleStartClick} onStop={lifecycle.handleStop}
                            stopping={lifecycle.stopping} starting={lifecycle.resuming}
                            onToggleNav={toggleNav} showNavToggle={sidebar.show}
                            onNewChat={sidebarHistory.onNewChat}
                            model={{
                                models: simpleChatController.models, value: simpleChatController.currentModel,
                                isLoading: simpleChatController.modelsLoading, onChange: simpleChatController.setModel,
                            }}
                        />
                    )}
                </div>
                <WorkbenchTweaksPanel ctx={ctx} effectiveMode={effectiveMode} status={status} vaultId={vaultId}/>
            </div>
            {sidebar.show && vaultId && (
                <div className={cls.SidebarSlot}>
                    <WorkbenchSidebar history={sidebarHistory} vault={sidebar.vault}
                        navOpen={navOpen} onToggleNav={toggleNav}/>
                </div>
            )}
        </div>
    )
}
