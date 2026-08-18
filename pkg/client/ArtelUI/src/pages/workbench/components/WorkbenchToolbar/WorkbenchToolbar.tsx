import {Button} from "@vervstack/chures"

import WorkbenchStatusBadge from "@/pages/workbench/components/WorkbenchStatusBadge/WorkbenchStatusBadge.tsx"
import WorkbenchSettingsMenu
    from "@/pages/workbench/components/WorkbenchToolbar/components/WorkbenchSettingsMenu/WorkbenchSettingsMenu.tsx"
import cls from "@/pages/workbench/components/WorkbenchToolbar/WorkbenchToolbar.module.css"

interface WorkbenchToolbarProps {
    vaultName: string
    status: string
    exists: boolean
    vaultId?: string
    onStart: () => void
    onStop: () => void
    stopping: boolean
}

// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function WorkbenchToolbar(props: WorkbenchToolbarProps) {
    const isRunning = props.status === "running"

    return (
        <div className={cls.WorkbenchToolbarContainer}>
            <div className={cls.LeftSection}>
                <span className={cls.VaultName}>{props.vaultName}</span>
                <WorkbenchStatusBadge status={props.exists ? props.status : "not_configured"}/>
            </div>
            {props.exists && (
                <div className={cls.RightSection}>
                    <Button
                        variant="secondary"
                        onClick={isRunning ? props.onStop : props.onStart}
                        disabled={props.stopping}
                    >
                        {props.stopping ? "Stopping…" : isRunning ? "Stop" : "Start"}
                    </Button>
                    {props.vaultId && <WorkbenchSettingsMenu vaultId={props.vaultId}/>}
                </div>
            )}
        </div>
    )
}
