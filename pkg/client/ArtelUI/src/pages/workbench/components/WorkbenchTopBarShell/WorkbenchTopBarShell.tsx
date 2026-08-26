import {ReactNode} from "react"

import cls from "@/pages/workbench/components/WorkbenchTopBarShell/WorkbenchTopBarShell.module.css"

interface Props {
    statusBadge?: ReactNode
    actions?: ReactNode
}

// Shared top-bar chrome (sticky positioning) used by the Docker workbench's
// WorkbenchToolbar — Simple Chat mode renders no header at all anymore.
export default function WorkbenchTopBarShell(props: Props) {
    return (
        <div className={cls.WorkbenchTopBarShellContainer}>
            <div className={cls.LeftSection}>
                {props.statusBadge}
            </div>
            {props.actions}
        </div>
    )
}
