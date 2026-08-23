import {Button} from "@vervstack/chures"
import {useNavigate} from "react-router-dom"

import {Path} from "@/app/routing/Router.tsx"
import WorkbenchStatusBadge from "@/pages/workbench/components/WorkbenchStatusBadge/WorkbenchStatusBadge.tsx"
import BackIcon from "@/pages/workbench/components/WorkbenchToolbar/components/BackIcon/BackIcon.tsx"
import ToolbarActions from "@/pages/workbench/components/WorkbenchToolbar/components/ToolbarActions/ToolbarActions.tsx"
import cls from "@/pages/workbench/components/WorkbenchToolbar/WorkbenchToolbar.module.css"
import {WorkbenchView} from "@/pages/workbench/components/WorkbenchToolbar/WorkbenchView.ts"

export type {WorkbenchView} from "@/pages/workbench/components/WorkbenchToolbar/WorkbenchView.ts"

interface WorkbenchToolbarProps {
    vaultName: string
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
    onToggleHistory: () => void
}

// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function WorkbenchToolbar(props: WorkbenchToolbarProps) {
    const isRunning = props.status === "running"
    const navigate = useNavigate()

    return (
        <div className={cls.WorkbenchToolbarContainer}>
            <div className={cls.LeftSection}>
                <Button
                    variant="secondary"
                    className={cls.BackButton}
                    onClick={() => navigate(Path.HomePage)}
                    aria-label="Back to vaults"
                    title="Back to vaults"
                >
                    <BackIcon/>
                </Button>
                <span className={cls.VaultName}>{props.vaultName}</span>
                <WorkbenchStatusBadge status={props.exists ? props.status : "not_configured"}/>
            </div>
            {props.exists && (
                <ToolbarActions
                    isRunning={isRunning}
                    view={props.view}
                    onViewChange={props.onViewChange}
                    chatLocked={props.chatLocked}
                    vaultId={props.vaultId}
                    onToggleHistory={props.onToggleHistory}
                    onStart={props.onStart}
                    onStop={props.onStop}
                    stopping={props.stopping}
                    starting={props.starting}
                />
            )}
        </div>
    )
}
