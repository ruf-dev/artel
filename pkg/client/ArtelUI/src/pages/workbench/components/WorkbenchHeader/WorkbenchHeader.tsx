import WorkbenchToolbar, {type WorkbenchView} from "@/pages/workbench/components/WorkbenchToolbar/WorkbenchToolbar.tsx"
import type {WorkbenchMode} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"

interface Props {
    isLoading: boolean
    effectiveMode: WorkbenchMode | "picking"
    exists: boolean
    vaultId?: string
    status: string
    onStart: () => void
    onStop: () => void
    stopping: boolean
    starting: boolean
    view: WorkbenchView
    onViewChange: (view: WorkbenchView) => void
    chatLocked: boolean
    onToggleHistory: () => void
}

// Renders the docker toolbar for docker mode, nothing otherwise (Simple Chat mode
// has no header — its model switcher/new-chat button live in its pinned history
// sidebar instead). Split out purely to keep WorkbenchPage.tsx's render function
// under the max-lines-per-function lint limit — same rationale as WorkbenchPanels.tsx.
//
// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function WorkbenchHeader(props: Props) {
    if (!props.isLoading && props.effectiveMode === "docker" && props.exists && props.vaultId) {
        return (
            <WorkbenchToolbar
                status={props.status}
                exists={props.exists}
                vaultId={props.vaultId}
                onStart={props.onStart}
                onStop={props.onStop}
                stopping={props.stopping}
                starting={props.starting}
                view={props.view}
                onViewChange={props.onViewChange}
                chatLocked={props.chatLocked}
                onToggleHistory={props.onToggleHistory}
            />
        )
    }

    return null
}
