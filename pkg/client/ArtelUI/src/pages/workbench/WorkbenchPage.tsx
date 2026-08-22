import {useEffect, useState} from "react"
import {Link, useParams} from "react-router-dom"
import {Button, Loader} from "@vervstack/chures"

import cls from "@/pages/workbench/WorkbenchPage.module.css"
import {cn} from "@/app/utils/cn.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useVaults} from "@/app/hooks/Vaults.ts"
import {
    useWorkbench,
    useWorkbenchTerminalTabs,
    useWorkbenchTerminalTabMutations,
} from "@/app/hooks/Workbench.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import PickAuthModeScreen from "@/pages/workbench/components/PickAuthModeScreen/PickAuthModeScreen.tsx"
import Chat from "@/pages/workbench/components/Chat/Chat.tsx"
import TerminalView from "@/pages/workbench/components/TerminalView/TerminalView.tsx"
import WorkbenchToolbar, {
    type WorkbenchView,
} from "@/pages/workbench/components/WorkbenchToolbar/WorkbenchToolbar.tsx"
import {useChatSession} from "@/pages/workbench/processes/useChatSession.ts"
import {useWorkbenchLifecycle} from "@/pages/workbench/processes/useWorkbenchLifecycle.ts"

export default function WorkbenchPage() {
    const {vaultId} = useParams()
    const {vaults} = useVaults()
    const {exists, status, isLoading, pendingTerminalAuthLink} = useWorkbench(vaultId)
    const lifecycle = useWorkbenchLifecycle(vaultId)
    const {tabs} = useWorkbenchTerminalTabs(vaultId, status === "running")
    const {create: createTab, select: selectTab, close: closeTab} = useWorkbenchTerminalTabMutations(vaultId)
    const bakeError = useBakeError()
    const {
        items,
        status: chatStatus,
        authComplete,
        sendMessage,
        sendPermissionDecision,
        sendAuthCode,
    } = useChatSession(status === "running" ? vaultId : undefined)

    const [view, setView] = useState<WorkbenchView>("chat")

    const vault = vaults.find(v => v.id === vaultId)
    const vaultName = vault?.name ?? "Vault"

    const awaitingAuth = !authComplete &&
        (items.some(i => i.kind === "auth_link" || i.kind === "auth_code_needed") ||
            lifecycle.pendingAuthMode === "subscription_login")

    const genericCentered = isLoading || !exists || (status !== "running" && !lifecycle.showSetup)
    const terminalViewActive = status === "running" && view === "terminal"

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
                {!isLoading && !exists && (
                    <>
                        <p className={cls.EmptyState}>
                            This vault doesn&apos;t have a Workbench. It may predate the feature, or Docker
                            isn&apos;t configured for this deployment.
                        </p>
                        <Button variant="primary" onClick={lifecycle.handleCreate} disabled={lifecycle.creating}>
                            {lifecycle.creating ? "Setting up…" : "Set up Workbench"}
                        </Button>
                    </>
                )}
                {!isLoading && exists && status === "running" && vaultId && !awaitingAuth && view === "chat" && (
                    <Chat items={items} status={chatStatus} sendMessage={sendMessage}
                          sendPermissionDecision={sendPermissionDecision} sendAuthCode={sendAuthCode}/>
                )}
                {!isLoading && exists && status === "running" && vaultId && view === "terminal" && (
                    <TerminalView
                        vaultId={vaultId}
                        tabs={tabs}
                        pendingTerminalAuthLink={pendingTerminalAuthLink}
                        onSelectTab={(tabId) => {
                            selectTab(tabId)
                                .catch(e => bakeError("Failed to select tab", e))
                        }}
                        onCreateTab={() => {
                            createTab()
                                .catch(e => bakeError("Failed to create tab", e))
                        }}
                        onCloseTab={(tabId) => {
                            closeTab(tabId)
                                .catch(e => bakeError("Failed to close tab", e))
                        }}
                    />
                )}
                {!isLoading && exists && status !== "running" && lifecycle.showSetup && vaultId && (
                    <PickAuthModeScreen
                        onStart={lifecycle.handleStartWorkbench}
                        starting={lifecycle.startingWorkbench}
                    />
                )}
            </div>
            {!isLoading && exists && vaultId ? (
                <WorkbenchToolbar
                    vaultName={vaultName}
                    status={status}
                    exists={exists}
                    vaultId={vaultId}
                    onStart={lifecycle.handleStartClick}
                    onStop={lifecycle.handleStop}
                    stopping={lifecycle.stopping}
                    starting={lifecycle.resuming}
                    view={view}
                    onViewChange={setView}
                    chatLocked={awaitingAuth}
                />
            ) : (
                <Link className={cls.BackLink} to={Path.HomePage}>← Back to vaults</Link>
            )}
        </div>
    )
}
