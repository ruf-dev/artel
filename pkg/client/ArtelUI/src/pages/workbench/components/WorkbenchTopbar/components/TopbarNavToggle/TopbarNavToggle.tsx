import {Button} from "@vervstack/chures"
import {MorphIcon} from "morphicons/react"
import {Menu} from "lucide"

import cls from "@/pages/workbench/components/WorkbenchTopbar/components/TopbarNavToggle/TopbarNavToggle.module.css"

interface Props {
    onToggle: () => void
}

// Mobile-only burger in the workbench top bar: opens the off-canvas WorkbenchSidebar
// drawer. Hidden above the mobile breakpoint via CSS (the desktop sidebar has its own
// in-rail toggle). Mirrors SidebarToggleButton's chures Button + MorphIcon shape.
export default function TopbarNavToggle({onToggle}: Props) {
    return (
        <div className={cls.TopbarNavToggleContainer}>
            <Button variant="unstyled" className={cls.Btn} aria-label="Open menu" onClick={onToggle}>
                <MorphIcon icon={Menu} size={20} strokeWidth={1.6} className={cls.Icon}/>
            </Button>
        </div>
    )
}
