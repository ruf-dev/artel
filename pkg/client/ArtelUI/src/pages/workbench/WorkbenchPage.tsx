import {useState} from "react"
import {Link, useParams} from "react-router-dom"
import {Button, Loader} from "@vervstack/chures"

import cls from "@/pages/workbench/WorkbenchPage.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {useVaults} from "@/app/hooks/Vaults.ts"
import {useWorkbench, useWorkbenchMutations} from "@/app/hooks/Workbench.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import HeroSegment from "@/components/HeroSegment/HeroSegment.tsx"
import WorkbenchStatusBadge from "@/pages/workbench/components/WorkbenchStatusBadge/WorkbenchStatusBadge.tsx"
import PickAuthModeScreen from "@/pages/workbench/components/PickAuthModeScreen/PickAuthModeScreen.tsx"
import LoginRelayScreen from "@/pages/workbench/components/LoginRelayScreen/LoginRelayScreen.tsx"

type Stage = "pick" | "login"

export default function WorkbenchPage() {
    const {vaultId} = useParams()
    const {vaults} = useVaults()
    const {exists, status, isLoading, refetch} = useWorkbench(vaultId)
    const {create, stop} = useWorkbenchMutations(vaultId)
    const bakeError = useBakeError()

    const [stage, setStage] = useState<Stage>("pick")
    const [showSetup, setShowSetup] = useState(false)
    const [stopping, setStopping] = useState(false)
    const [creating, setCreating] = useState(false)

    const vault = vaults.find(v => v.id === vaultId)
    const vaultName = vault?.name ?? "Vault"

    function handleCreate() {
        setCreating(true)
        create()
            .then(() => setShowSetup(true))
            .catch(e => bakeError("Failed to set up workbench", e))
            .finally(() => setCreating(false))
    }

    function handleStarted() {
        setStage("login")
    }

    function handleStop() {
        setStopping(true)
        stop()
            .catch(e => bakeError("Failed to stop workbench", e))
            .finally(() => setStopping(false))
    }

    const displayStatus = exists ? status : "not_configured"

    return (
        <div className={cls.WorkbenchPageContainer}>
            <Link className={cls.BackLink} to={Path.HomePage}>← Back to vaults</Link>
            <HeroSegment
                eyebrow="Workbench"
                title={vaultName}
                subtitle={
                    <>
                        <WorkbenchStatusBadge status={displayStatus}/>
                        {" · "}<span>Run Claude Code against this vault in an isolated container.</span>
                    </>
                }
                action={status === "running" ? (
                    <Button variant="secondary" onClick={handleStop} disabled={stopping}>
                        {stopping ? "Stopping…" : "Stop"}
                    </Button>
                ) : undefined}
            />
            <div className={cls.Body}>
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
                {!isLoading && exists && status === "running" && (
                    <p className={cls.RunningState}>
                        Claude Code is running in an isolated container for this vault.
                    </p>
                )}
                {!isLoading && exists && status !== "running" && !showSetup && (
                    <Button variant="primary" onClick={() => setShowSetup(true)}>
                        Start Workbench
                    </Button>
                )}
                {!isLoading && exists && status !== "running" && showSetup && vaultId && (
                    stage === "pick"
                        ? <PickAuthModeScreen vaultId={vaultId} onStarted={handleStarted}/>
                        : <LoginRelayScreen vaultId={vaultId} onDone={refetch}/>
                )}
            </div>
        </div>
    )
}
