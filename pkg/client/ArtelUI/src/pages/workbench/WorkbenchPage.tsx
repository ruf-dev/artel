import {useEffect, useState} from "react"
import {useParams} from "react-router-dom"
import {Loader} from "@vervstack/chures"

import cls from "@/pages/workbench/WorkbenchPage.module.css"
import {cn} from "@/app/utils/cn.ts"
import {useVaults} from "@/app/hooks/Vaults.ts"
import {
    useWorkbench,
    useWorkbenchTerminalTabs,
    useWorkbenchTerminalTabMutations,
} from "@/app/hooks/Workbench.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import PickAuthModeScreen from "@/pages/workbench/components/PickAuthModeScreen/PickAuthModeScreen.tsx"
import PickWorkbenchModeScreen from "@/pages/workbench/components/PickWorkbenchModeScreen/PickWorkbenchModeScreen.tsx"
import WorkbenchHeader from "@/pages/workbench/components/WorkbenchHeader/WorkbenchHeader.tsx"
import WorkbenchPanels from "@/pages/workbench/components/WorkbenchPanels/WorkbenchPanels.tsx"
import {type WorkbenchView} from "@/pages/workbench/components/WorkbenchToolbar/WorkbenchToolbar.tsx"
import {useChatSession} from "@/pages/workbench/processes/useChatSession.ts"
import {
    toSimpleChatSessionBundle,
    toSimpleChatTopBarProps,
    useSimpleChatController,
} from "@/pages/workbench/processes/useSimpleChatController.ts"
import {useWorkbenchLifecycle} from "@/pages/workbench/processes/useWorkbenchLifecycle.ts"
import {useWorkbenchPanelControls} from "@/pages/workbench/processes/useWorkbenchPanelControls.ts"
import {useWorkbenchModeControls} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"

export default function WorkbenchPage() {
    const {vaultId} = useParams()
    const {vaults} = useVaults()
    const {exists, status, isLoading, pendingTerminalAuthLink} = useWorkbench(vaultId)
    const lifecycle = useWorkbenchLifecycle(vaultId)
    const {tabs} = useWorkbenchTerminalTabs(vaultId, status === "running")
    const {create: createTab, select: selectTab, close: closeTab} = useWorkbenchTerminalTabMutations(vaultId)
    const bakeError = useBakeError()
    const [view, setView] = useState<WorkbenchView>("chat")
    const {historyOpen, handleSelectTab, handleCreateTab, handleCloseTab, toggleHistory, closeHistory} =
        useWorkbenchPanelControls({selectTab, createTab, closeTab, bakeError})

    const modeControls = useWorkbenchModeControls({exists, handleCreateDocker: lifecycle.handleCreate})
    const {effectiveMode} = modeControls

    // Only connect the Docker bridge WebSocket while Docker mode is actually the
    // active mode — otherwise a vault with both a running Docker workbench and an
    // active Simple Chat thread would hold two live sockets open at once for no
    // reason, since Simple Chat's own WS connection lives in SimpleChat.tsx.
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

    const awaitingAuth = !chatSession.authComplete && lifecycle.pendingAuthMode === "subscription_login"

    const dockerCentered = !exists || (status !== "running" && !lifecycle.showSetup)
    const modeCentered = effectiveMode === "docker" ? dockerCentered : false
    const genericCentered = isLoading || modeCentered
    const terminalViewActive = effectiveMode === "docker" && status === "running" && view === "terminal"

    // Terminal login and the chat's own sign-in flow share the same in-container
    // credentials file, so there's no chat-side auth screen anymore — lock the Chat
    // toggle to Terminal while unauthenticated instead. Unlock-only: once auth
    // completes this stops firing, but view is never forced back to "chat" on its own.
    useEffect(() => {
        if (awaitingAuth && view === "chat") {
            setView("terminal")
        }
    }, [awaitingAuth, view])

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
                        simpleHistoryOpen={modeControls.simpleHistoryOpen}
                        onCloseSimpleHistory={modeControls.closeSimpleHistory}
                        simpleChatSession={toSimpleChatSessionBundle(simpleChatController)}
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
                vaultName={vaultName}
                status={status}
                onStart={lifecycle.handleStartClick}
                onStop={lifecycle.handleStop}
                stopping={lifecycle.stopping}
                starting={lifecycle.resuming}
                view={view}
                onViewChange={setView}
                chatLocked={awaitingAuth}
                onToggleHistory={toggleHistory}
                simpleChatTopBar={toSimpleChatTopBarProps(simpleChatController, modeControls.toggleSimpleHistory)}
            />
        </div>
    )
}
