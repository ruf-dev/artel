import {cn} from "@/app/utils/cn"
import cls from "@/pages/tract-canvas/components/RunStatusBadge/RunStatusBadge.module.css"

interface Props {
    running: boolean
    lastRunStatus?: string
}

export default function RunStatusBadge({running, lastRunStatus}: Props) {
    return (
        <span className={cls.RunBadge}>
            <span className={cn(cls.BadgeDot, running && cls.BadgeDotRunning)}/>
            {running ? "Running…" : lastRunStatus ? `Last run: ${lastRunStatus}` : "No runs yet"}
        </span>
    )
}
