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
import WorkbenchHeader from "@/pages/workbench/components/WorkbenchHeader/WorkbenchHeader.tsx"
import WorkbenchPanels from "@/pages/workbench/components/WorkbenchPanels/WorkbenchPanels.tsx"
import {useChatSession} from "@/pages/workbench/processes/useChatSession.ts"
import {
    toSimpleChatSessionBundle,
    useSimpleChatController,
} from "@/pages/workbench/processes/useSimpleChatController.ts"
import {useWorkbenchLifecycle} from "@/pages/workbench/processes/useWorkbenchLifecycle.ts"
import {useWorkbenchPanelControls} from "@/pages/workbench/processes/useWorkbenchPanelControls.ts"
import {useWorkbenchModeControls} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"
import {useWorkbenchViewState} from "@/pages/workbench/processes/useWorkbenchViewState.ts"

export default function WorkbenchPage() {
    const {vaultId} = useParams()
    const {vaults} = useVaults()
    const {exists, status, isLoading, pendingTerminalAuthLink} = useWorkbench(vaultId)
    const {chats: simpleChats, isLoading: simpleChatsLoading} = useSimpleChats(vaultId)
    const lifecycle = useWorkbenchLifecycle(vaultId)
    const {tabs} = useWorkbenchTerminalTabs(vaultId, status === "running")
    const {create: createTab, select: selectTab, close: closeTab} = useWorkbenchTerminalTabMutations(vaultId)
    const bakeError = useBakeError()
    const {historyOpen, handleSelectTab, handleCreateTab, handleCloseTab, toggleHistory, closeHistory} =
        useWorkbenchPanelControls({selectTab, createTab, closeTab, bakeError})
    const modeControls = useWorkbenchModeControls({
        exists,
        handleCreateDocker: lifecycle.handleCreate,
        simpleChats,
        simpleChatsLoading,
    })
    const {effectiveMode} = modeControls
    // Connect Docker's bridge only while Docker mode is active; Simple Chat owns its own socket.
    const chatSession = useChatSession(effectiveMode === "docker" && status === "running" ? vaultId : undefined)

    const simpleChatController = useSimpleChatController({
        chatId: modeControls.simpleChatId,
        vaultId: vaultId ?? "",
        active: effectiveMode === "simple-chat",
        onChatCreated: modeControls.setSimpleChatId,
        bakeError,
    })

    const vault = vaults.find(v => v.id === vaultId)
    const vaultName = vault?.name ?? "Vault"

    const {view, setView, awaitingAuth, genericCentered, terminalViewActive} = useWorkbenchViewState({
        exists,
        status,
        isLoading,
        effectiveMode,
        showSetup: lifecycle.showSetup,
        authComplete: chatSession.authComplete,
        pendingAuthMode: lifecycle.pendingAuthMode,
    })

    useDocumentTitle(vaultName)

    return (
        <div className={cn(cls.WorkbenchPageContainer, terminalViewActive && cls.TerminalViewActive)}>
            <div className={cn(cls.Body, genericCentered && cls.BodyCentered)}>
                {isLoading && (
                    <Loader variant="arcs" size="sm" color="var(--coral)"/>
                )}
                {!isLoading && effectiveMode === "picking" && vaultId && (
                    <PickWorkbenchModeScreen
                        vaultId={vaultId}
                        onStartDocker={modeControls.handlePickDocker}
                        startingDocker={lifecycle.creating}
                        onSimpleChatCreated={modeControls.handleSimpleChatCreated}
                    />
                )}
                {!isLoading && vaultId && (
                    <WorkbenchPanels
                        effectiveMode={effectiveMode}
                        exists={exists}
                        status={status}
                        vaultId={vaultId}
                        view={view}
                        awaitingAuth={awaitingAuth}
                        chatSession={chatSession}
                        historyOpen={historyOpen}
                        onCloseHistory={closeHistory}
                        tabs={tabs}
                        pendingTerminalAuthLink={pendingTerminalAuthLink}
                        onSelectTab={handleSelectTab}
                        onCreateTab={handleCreateTab}
                        onCloseTab={handleCloseTab}
                        simpleChatId={modeControls.simpleChatId}
                        onSelectSimpleChat={modeControls.setSimpleChatId}
                        simpleChatSession={toSimpleChatSessionBundle(simpleChatController)}
                        onCloseSimpleChat={modeControls.handleCloseSimpleChat}
                    />
                )}
                {!isLoading && effectiveMode === "docker" && exists && status !== "running" && lifecycle.showSetup
                    && vaultId && (
                    <PickAuthModeScreen
                        onStart={lifecycle.handleStartWorkbench}
                        starting={lifecycle.startingWorkbench}
                    />
                )}
            </div>
            <WorkbenchHeader
                isLoading={isLoading}
                effectiveMode={effectiveMode}
                exists={exists}
                vaultId={vaultId}
                status={status}
                onStart={lifecycle.handleStartClick}
                onStop={lifecycle.handleStop}
                stopping={lifecycle.stopping}
                starting={lifecycle.resuming}
                view={view}
                onViewChange={setView}
                chatLocked={awaitingAuth}
                onToggleHistory={toggleHistory}
            />
        </div>
    )
}
