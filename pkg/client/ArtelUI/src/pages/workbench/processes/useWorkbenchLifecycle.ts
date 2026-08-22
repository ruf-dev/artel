import {useEffect, useState} from "react"

import {useWorkbenchMutations} from "@/app/hooks/Workbench.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"

// Bundles the create/stop mutation calls for a workbench together with the
// local loading-state flags and setup-screen gating they drive, keeping
// WorkbenchPage.tsx focused on layout/branching rather than this orchestration.
export function useWorkbenchLifecycle(vaultId: string | undefined) {
    const {create, stop, start} = useWorkbenchMutations(vaultId)
    const {connections, fetch: fetchConnections} = useExternalConnections()
    const bakeError = useBakeError()

    const [showSetup, setShowSetup] = useState(false)
    const [stopping, setStopping] = useState(false)
    const [creating, setCreating] = useState(false)
    const [resuming, setResuming] = useState(false)
    const [pendingAuthMode, setPendingAuthMode] = useState<string | undefined>(undefined)
    const [startingWorkbench, setStartingWorkbench] = useState(false)

    useEffect(() => {
        void fetchConnections()
    }, [fetchConnections])

    const hasAnthropicConnection = connections
        .some(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_ANTHROPIC)

    // Terminal login and the chat's own sign-in flow both land in the same
    // in-container credentials file, so when there's no API-key connection to
    // choose between, skip the picker entirely and go straight to subscription
    // login instead of showing a redundant "how do you want to auth" screen.
    function handleCreate() {
        setCreating(true)
        create()
            .then(() => {
                if (hasAnthropicConnection) {
                    setShowSetup(true)
                } else {
                    void handleStartWorkbench("subscription_login")
                }
            })
            .catch(e => bakeError("Failed to set up workbench", e))
            .finally(() => setCreating(false))
    }

    function handleStartClick() {
        setResuming(true)
        create()
            .then(() => {
                if (hasAnthropicConnection) {
                    setShowSetup(true)
                } else {
                    void handleStartWorkbench("subscription_login")
                }
            })
            .catch(e => bakeError("Failed to prepare workbench", e))
            .finally(() => setResuming(false))
    }

    function handleStop() {
        setStopping(true)
        stop()
            .catch(e => bakeError("Failed to stop workbench", e))
            .finally(() => setStopping(false))
    }

    function handleStartWorkbench(authMode: string): Promise<void> {
        setStartingWorkbench(true)
        setPendingAuthMode(authMode)
        return start(authMode)
            .then(() => {})
            .catch(e => {
                setPendingAuthMode(undefined)
                bakeError("Failed to start workbench", e)
            })
            .finally(() => setStartingWorkbench(false))
    }

    return {
        showSetup,
        stopping,
        creating,
        resuming,
        pendingAuthMode,
        startingWorkbench,
        handleCreate,
        handleStartClick,
        handleStop,
        handleStartWorkbench,
    }
}
