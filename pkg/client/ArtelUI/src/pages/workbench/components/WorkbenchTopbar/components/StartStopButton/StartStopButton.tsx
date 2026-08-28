import {Button} from "@vervstack/chures"
import {MorphIcon} from "morphicons/react"
import {Play, Square} from "lucide"

import {cn} from "@/app/utils/cn.ts"
import cls from "@/pages/workbench/components/WorkbenchTopbar/components/StartStopButton/StartStopButton.module.css"

interface Props {
    isRunning: boolean
    onStart: () => void
    onStop: () => void
    stopping: boolean
    starting: boolean
}

// Start/stop the Docker workbench container. Carried over from the old
// ToolbarActions — a Play glyph morphs to a Square once the container is running.
export default function StartStopButton({isRunning, onStart, onStop, stopping, starting}: Props) {
    return (
        <div className={cls.StartStopButtonContainer}>
            <Button
                variant="secondary"
                className={cls.Btn}
                onClick={isRunning ? onStop : onStart}
                disabled={stopping || starting}
                aria-label={stopping ? "Stopping" : isRunning ? "Stop" : "Start"}
            >
                <MorphIcon
                    icon={isRunning ? Square : Play}
                    size={20}
                    strokeWidth={1.6}
                    className={cn(cls.Icon, isRunning && cls.IconRunning)}
                />
            </Button>
        </div>
    )
}
