import WorkbenchStatusBadge from "@/pages/workbench/components/WorkbenchStatusBadge/WorkbenchStatusBadge.tsx"
import ToolbarActions from "@/pages/workbench/components/WorkbenchToolbar/components/ToolbarActions/ToolbarActions.tsx"
import WorkbenchTopBarShell from "@/pages/workbench/components/WorkbenchTopBarShell/WorkbenchTopBarShell.tsx"
import {WorkbenchView} from "@/pages/workbench/components/WorkbenchToolbar/WorkbenchView.ts"

export type {WorkbenchView} from "@/pages/workbench/components/WorkbenchToolbar/WorkbenchView.ts"

interface WorkbenchToolbarProps {
    status: string
    exists: boolean
    vaultId?: string
    onStart: () => void
    onStop: () => void
    stopping: boolean
    starting: boolean
    view: WorkbenchView
    onViewChange: (view: WorkbenchView) => void
    chatLocked: boolean
}

// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function WorkbenchToolbar(props: WorkbenchToolbarProps) {
    const isRunning = props.status === "running"

    return (
        <WorkbenchTopBarShell
            statusBadge={<WorkbenchStatusBadge status={props.exists ? props.status : "not_configured"}/>}
            actions={props.exists && (
                <ToolbarActions
                    isRunning={isRunning}
                    view={props.view}
                    onViewChange={props.onViewChange}
                    chatLocked={props.chatLocked}
                    vaultId={props.vaultId}
                    onStart={props.onStart}
                    onStop={props.onStop}
                    stopping={props.stopping}
                    starting={props.starting}
                />
            )}
        />
    )
}
