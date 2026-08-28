import {useNavigate} from "react-router-dom"
import {Button} from "@vervstack/chures"

import {Path} from "@/app/routing/Router.tsx"
import cls from "@/pages/workbench/components/WorkbenchSidebar/components/SidebarFooter/SidebarFooter.module.css"
import SidebarGearIcon from "@/pages/workbench/components/WorkbenchSidebar/components/icons/SidebarGearIcon.tsx"

// The sidebar's pinned footer — jumps to the connections/settings page.
export default function SidebarFooter() {
    const navigate = useNavigate()

    return (
        <Button
            variant="unstyled"
            className={cls.SidebarFooterContainer}
            onClick={() => navigate(Path.ConnectionsPage)}
        >
            <SidebarGearIcon className={cls.Gear}/>
            <span className={cls.Label}>Settings &amp; connections</span>
        </Button>
    )
}
