import {useState} from "react"
import {Link, useParams} from "react-router-dom"
import {Button, Loader} from "@vervstack/chures"

import cls from "@/pages/workbench/WorkbenchPage.module.css"
import {cn} from "@/app/utils/cn.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useVaults} from "@/app/hooks/Vaults.ts"
import {
    useWorkbench,
    useWorkbenchMutations,
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

export default function WorkbenchPage() {
    const {vaultId} = useParams()
    const {vaults} = useVaults()
    const {exists, status, isLoading} = useWorkbench(vaultId)
    const {create, stop} = useWorkbenchMutations(vaultId)
    const {tabs} = useWorkbenchTerminalTabs(vaultId, status === "running")
    const {create: createTab, select: selectTab, close: closeTab} = useWorkbenchTerminalTabMutations(vaultId)
    const bakeError = useBakeError()

    const [showSetup, setShowSetup] = useState(false)
    const [stopping, setStopping] = useState(false)
    const [creating, setCreating] = useState(false)
    const [resuming, setResuming] = useState(false)
    const [view, setView] = useState<WorkbenchView>("chat")

    const vault = vaults.find(v => v.id === vaultId)
    const vaultName = vault?.name ?? "Vault"

    function handleCreate() {
        setCreating(true)
        create()
            .then(() => setShowSetup(true))
            .catch(e => bakeError("Failed to set up workbench", e))
            .finally(() => setCreating(false))
    }

    function handleStartClick() {
        setResuming(true)
        create()
            .then(() => setShowSetup(true))
            .catch(e => bakeError("Failed to prepare workbench", e))
            .finally(() => setResuming(false))
    }

    function handleStop() {
        setStopping(true)
        stop()
            .catch(e => bakeError("Failed to stop workbench", e))
            .finally(() => setStopping(false))
    }

    const bodyCentered = isLoading || !exists || (status !== "running" && !showSetup)
    const terminalViewActive = status === "running" && view === "terminal"

    return (
        <div className={cn(cls.WorkbenchPageContainer, terminalViewActive && cls.TerminalViewActive)}>
            <Link className={cls.BackLink} to={Path.HomePage}>← Back to vaults</Link>
            {!isLoading && exists && vaultId && (
                <WorkbenchToolbar
                    vaultName={vaultName}
                    status={status}
                    exists={exists}
                    vaultId={vaultId}
                    onStart={handleStartClick}
                    onStop={handleStop}
                    stopping={stopping}
                    starting={resuming}
                    view={view}
                    onViewChange={setView}
                />
            )}
            <div className={cn(cls.Body, bodyCentered && cls.BodyCentered)}>
                {isLoading && (
                    <Loader variant="arcs" size="sm" color="var(--coral)"/>
                )}
                {!isLoading && !exists && (
                    <>
                        <p className={cls.EmptyState}>
                            This vault doesn&apos;t have a Workbench. It may predate the feature, or Docker
                            isn&apos;t configured for this deployment.
                        </p>
                        <Button variant="primary" onClick={handleCreate} disabled={creating}>
                            {creating ? "Setting up…" : "Set up Workbench"}
                        </Button>
                    </>
                )}
                {!isLoading && exists && status === "running" && vaultId && view === "chat" && (
                    <Chat vaultId={vaultId}/>
                )}
                {!isLoading && exists && status === "running" && vaultId && view === "terminal" && (
                    <TerminalView
                        vaultId={vaultId}
                        tabs={tabs}
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
                {!isLoading && exists && status !== "running" && showSetup && vaultId && (
                    <PickAuthModeScreen vaultId={vaultId}/>
                )}
            </div>
        </div>
    )
}
