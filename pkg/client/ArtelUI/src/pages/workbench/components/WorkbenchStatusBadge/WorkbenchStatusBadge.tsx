import cls from "@/pages/workbench/components/WorkbenchStatusBadge/WorkbenchStatusBadge.module.css"
import {cn} from "@/app/utils/cn.ts"

interface Props {
    status: string
}

const STATUS_LABELS: Record<string, string> = {
    not_configured: "Not configured",
    created: "Not started",
    configuring: "Configuring",
    running: "Running",
    stopped: "Stopped",
    removed: "Removed",
}

const STATUS_CLASSES: Record<string, string> = {
    not_configured: "StatusNotConfigured",
    created: "StatusCreated",
    configuring: "StatusConfiguring",
    running: "StatusRunning",
    stopped: "StatusStopped",
    removed: "StatusRemoved",
}

export default function WorkbenchStatusBadge({status}: Props) {
    const label = STATUS_LABELS[status] ?? "Unknown"
    const statusClass = STATUS_CLASSES[status] ?? "StatusUnknown"
    return <span className={cn(cls.WorkbenchStatusBadgeContainer, cls[statusClass])}>{label}</span>
}
