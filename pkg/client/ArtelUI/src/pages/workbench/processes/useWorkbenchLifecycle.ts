import {useState} from "react"

import {useWorkbenchMutations} from "@/app/hooks/Workbench.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

// Bundles the create/stop mutation calls for a workbench together with the
// local loading-state flags and setup-screen gating they drive, keeping
// WorkbenchPage.tsx focused on layout/branching rather than this orchestration.
export function useWorkbenchLifecycle(vaultId: string | undefined) {
    const {create, stop} = useWorkbenchMutations(vaultId)
    const bakeError = useBakeError()

    const [showSetup, setShowSetup] = useState(false)
    const [stopping, setStopping] = useState(false)
    const [creating, setCreating] = useState(false)
    const [resuming, setResuming] = useState(false)

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

    return {
        showSetup,
        stopping,
        creating,
        resuming,
        handleCreate,
        handleStartClick,
        handleStop,
    }
}
