import {ReactNode} from "react"

import {cn} from "@/app/utils/cn.ts"
import cls from "@/pages/workbench/components/HistorySidebarShell/HistorySidebarShell.module.css"

interface Props {
    open: boolean
    children: ReactNode
}

// Shared positioned/animated outer shell for the left-anchored history sidebar,
// reused by the Docker workbench's ChatHistorySidebar and Simple Chat's
// SimpleChatHistorySidebar — each mode renders its own list/detail screens as
// children and keeps its own data source.
export default function HistorySidebarShell({open, children}: Props) {
    return (
        <div className={cn(cls.HistorySidebarShellContainer, open && cls.HistorySidebarShellOpen)}>
            {children}
        </div>
    )
}
