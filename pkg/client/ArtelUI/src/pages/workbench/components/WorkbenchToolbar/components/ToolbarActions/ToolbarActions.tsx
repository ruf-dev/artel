import {Button} from "@vervstack/chures"
import {AnimatePresence, motion} from "framer-motion"

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
    onToggleHistory: () => void
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
            <AnimatePresence mode="popLayout">
                {props.isRunning && props.view === "chat" && props.vaultId && (
                    <motion.div
                        layout
                        className={cls.HistoryButtonWrapper}
                        initial={{opacity: 0, scale: 0.7, y: -6}}
                        animate={{opacity: 1, scale: 1, y: 0}}
                        exit={{opacity: 0, scale: 0.7, y: -6}}
                        transition={{duration: 0.18}}
                    >
                        <Button
                            variant="secondary"
                            className={cls.HistoryButton}
                            onClick={props.onToggleHistory}
                            aria-label="View chat history"
                            title="View chat history"
                        >
                            <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
                                 strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                                <circle cx="12" cy="12" r="10"/>
                                <polyline points="12 6 12 12 16 14"/>
                            </svg>
                        </Button>
                    </motion.div>
                )}
            </AnimatePresence>
            <motion.div layout className={cls.StartStopButtonWrapper}>
                <Button
                    variant="secondary"
                    className={cls.StartStopButton}
                    onClick={props.isRunning ? props.onStop : props.onStart}
                    disabled={props.stopping || props.starting}
                    aria-label={props.stopping ? "Stopping" : props.isRunning ? "Stop" : "Start"}
                >
                    {props.isRunning ? (
                        <svg viewBox="0 0 24 24" width="20" height="20" fill="var(--color-error)">
                            <rect x="5" y="5" width="14" height="14" rx="1.5"/>
                        </svg>
                    ) : (
                        <svg viewBox="0 0 24 24" width="20" height="20" fill="#a0dc8c">
                            <polygon points="6 4 20 12 6 20"/>
                        </svg>
                    )}
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
