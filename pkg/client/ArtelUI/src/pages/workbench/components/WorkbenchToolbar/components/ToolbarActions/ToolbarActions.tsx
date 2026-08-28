import {Button} from "@vervstack/chures"
import {motion} from "framer-motion"
import {MorphIcon} from "morphicons/react"
import {Play, Square} from "lucide"

import {cn} from "@/app/utils/cn.ts"
import WorkbenchSettingsMenu
    from "@/pages/workbench/components/WorkbenchToolbar/components/WorkbenchSettingsMenu/WorkbenchSettingsMenu.tsx"
import {WorkbenchView} from "@/pages/workbench/components/WorkbenchToolbar/WorkbenchView.ts"
import cls from "@/pages/workbench/components/WorkbenchToolbar/components/ToolbarActions/ToolbarActions.module.css"

interface Props {
    isRunning: boolean
    view: WorkbenchView
    onViewChange: (view: WorkbenchView) => void
    chatLocked: boolean
    vaultId?: string
    onStart: () => void
    onStop: () => void
    stopping: boolean
    starting: boolean
}

// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function ToolbarActions(props: Props) {
    return (
        <div className={cls.ToolbarActionsContainer}>
            {props.isRunning && (
                <motion.div layout className={cls.ViewToggle}>
                    <Button
                        variant="secondary"
                        className={cn(
                            cls.ViewToggleButton,
                            props.view === "chat" && cls.ViewToggleButtonActive,
                        )}
                        onClick={() => props.onViewChange("chat")}
                        aria-pressed={props.view === "chat"}
                        disabled={props.chatLocked}
                        title={props.chatLocked ? "Log in via the terminal first" : undefined}
                    >
                        Chat
                    </Button>
                    <Button
                        variant="secondary"
                        className={cn(
                            cls.ViewToggleButton,
                            props.view === "terminal" && cls.ViewToggleButtonActive,
                        )}
                        onClick={() => props.onViewChange("terminal")}
                        aria-pressed={props.view === "terminal"}
                    >
                        Terminal
                    </Button>
                </motion.div>
            )}
            <motion.div layout className={cls.StartStopButtonWrapper}>
                <Button
                    variant="secondary"
                    className={cls.StartStopButton}
                    onClick={props.isRunning ? props.onStop : props.onStart}
                    disabled={props.stopping || props.starting}
                    aria-label={props.stopping ? "Stopping" : props.isRunning ? "Stop" : "Start"}
                >
                    <MorphIcon
                        icon={props.isRunning ? Square : Play}
                        size={20}
                        strokeWidth={1.6}
                        className={cn(cls.StartStopIcon, props.isRunning && cls.StartStopIconRunning)}
                    />
                </Button>
            </motion.div>
            {props.vaultId && (
                <motion.div layout className={cls.SettingsMenuWrapper}>
                    <WorkbenchSettingsMenu vaultId={props.vaultId}/>
                </motion.div>
            )}
        </div>
    )
}
