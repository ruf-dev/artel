import cls from "@/pages/tract-canvas/components/RunStatusDot/RunStatusDot.module.css"

import {cn} from "@/app/utils/cn.ts"
import {TractLastRun} from "@/processes/Tracts.ts"

interface Props {
    lastRun?: TractLastRun
}

export default function RunStatusDot({lastRun}: Props) {
    if (!lastRun) {
        return (
            <span className={cls.RunStatusDotContainer}>
                <span className={cn(cls.Dot, cls.DotNever)}/>never run
            </span>
        )
    }
    if (lastRun.status === "done") {
        return (
            <span className={cls.RunStatusDotContainer}>
                <span className={cn(cls.Dot, cls.DotOk)}/>{formatRelative(lastRun.at)}
            </span>
        )
    }
    if (lastRun.status === "failed") {
        return (
            <span className={cls.RunStatusDotContainer}>
                <span className={cn(cls.Dot, cls.DotFailed)}/>failed {formatRelative(lastRun.at)}
            </span>
        )
    }
    return (
        <span className={cls.RunStatusDotContainer}>
            <span className={cn(cls.Dot, cls.DotRunning)}/>running…
        </span>
    )
}

function formatRelative(iso: string): string {
    const then = new Date(iso).getTime()
    if (isNaN(then)) return ""
    const diffMs = Date.now() - then
    const min = Math.round(diffMs / 60000)
    if (min < 1) return "just now"
    if (min < 60) return `${min} min ago`
    const hr = Math.round(min / 60)
    if (hr < 24) return `${hr} h ago`
    return `${Math.round(hr / 24)} d ago`
}
