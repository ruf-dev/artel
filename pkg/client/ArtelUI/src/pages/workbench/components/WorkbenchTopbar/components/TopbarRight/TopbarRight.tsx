import SegmentedControl from "@/components/atoms/SegmentedControl/SegmentedControl.tsx"
import IconToggleButton
    from "@/pages/workbench/components/WorkbenchTopbar/components/IconToggleButton/IconToggleButton.tsx"
import StartStopButton
    from "@/pages/workbench/components/WorkbenchTopbar/components/StartStopButton/StartStopButton.tsx"
import WorkbenchSettingsMenu
    from "@/pages/workbench/components/WorkbenchTopbar/components/WorkbenchSettingsMenu/WorkbenchSettingsMenu.tsx"
import TweaksIcon from "@/pages/workbench/components/WorkbenchTopbar/components/icons/TweaksIcon.tsx"
import type {WorkbenchMode} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"
import type {WorkbenchContext} from "@/pages/workbench/processes/workbenchContext.ts"
import type {WorkbenchView} from "@/pages/workbench/processes/workbenchView.ts"
import cls from "@/pages/workbench/components/WorkbenchTopbar/components/TopbarRight/TopbarRight.module.css"

interface Props {
    effectiveMode: WorkbenchMode | "picking"
    exists: boolean
    status: string
    vaultId?: string
    view: WorkbenchView
    onViewChange: (view: WorkbenchView) => void
    chatLocked: boolean
    onStart: () => void
    onStop: () => void
    stopping: boolean
    starting: boolean
    ctx: WorkbenchContext
}

// Right cluster of WorkbenchTopbar: the Chat/Terminal switch (running docker only),
// start/stop + settings (any provisioned docker workbench), and the always-present
// tweaks toggle that opens/closes the WorkbenchTweaksPanel.
//
// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function TopbarRight(props: Props) {
    const showViewSwitch = props.effectiveMode === "docker" && props.status === "running"
    const showLifecycle = props.effectiveMode === "docker" && props.exists

    return (
        <div className={cls.TopbarRightContainer}>
            {showViewSwitch && (
                <SegmentedControl
                    options={[
                        {
                            key: "chat",
                            label: "Chat",
                            disabled: props.chatLocked,
                            tooltip: props.chatLocked ? "Log in via the terminal first" : undefined,
                        },
                        {key: "terminal", label: "Terminal"},
                    ]}
                    active={props.view}
                    onChange={key => props.onViewChange(key as WorkbenchView)}
                />
            )}
            {showLifecycle && (
                <StartStopButton
                    isRunning={props.status === "running"}
                    onStart={props.onStart}
                    onStop={props.onStop}
                    stopping={props.stopping}
                    starting={props.starting}
                />
            )}
            {showLifecycle && props.vaultId && <WorkbenchSettingsMenu vaultId={props.vaultId}/>}
            <IconToggleButton
                icon={<TweaksIcon/>}
                label="Tweaks"
                active={props.ctx.tweaksOpen}
                onClick={() => props.ctx.tweaksOpen ? props.ctx.closeTweaks() : props.ctx.openTweaks()}
            />
        </div>
    )
}
