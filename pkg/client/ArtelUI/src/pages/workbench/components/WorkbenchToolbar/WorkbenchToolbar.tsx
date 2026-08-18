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
                        className={cls.StartStopButton}
                        onClick={isRunning ? props.onStop : props.onStart}
                        disabled={props.stopping}
                        aria-label={props.stopping ? "Stopping" : isRunning ? "Stop" : "Start"}
                    >
                        {isRunning ? (
                            <svg viewBox="0 0 24 24" width="20" height="20" fill="var(--color-error)">
                                <rect x="5" y="5" width="14" height="14" rx="1.5"/>
                            </svg>
                        ) : (
                            <svg viewBox="0 0 24 24" width="20" height="20" fill="#a0dc8c">
                                <polygon points="6 4 20 12 6 20"/>
                            </svg>
                        )}
                    </Button>
                    {props.vaultId && <WorkbenchSettingsMenu vaultId={props.vaultId}/>}
                </div>
            )}
        </div>
    )
}
