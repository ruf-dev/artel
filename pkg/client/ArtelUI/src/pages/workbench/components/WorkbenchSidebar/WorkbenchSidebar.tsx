import {useState} from "react"

import SegmentedControl from "@/components/atoms/SegmentedControl/SegmentedControl.tsx"
import {WorkbenchHistory} from "@/pages/workbench/processes/useWorkbenchHistory.ts"
import {VaultPaneBundle} from "@/pages/workbench/processes/useWorkbenchSidebar.ts"
import SidebarBrand from "@/pages/workbench/components/WorkbenchSidebar/components/SidebarBrand/SidebarBrand.tsx"
// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import SidebarToggleButton from "@/pages/workbench/components/WorkbenchSidebar/components/SidebarToggleButton/SidebarToggleButton.tsx"
import HistoryPane from "@/pages/workbench/components/WorkbenchSidebar/components/HistoryPane/HistoryPane.tsx"
import VaultPane from "@/pages/workbench/components/WorkbenchSidebar/components/VaultPane/VaultPane.tsx"
import SidebarFooter from "@/pages/workbench/components/WorkbenchSidebar/components/SidebarFooter/SidebarFooter.tsx"
import HistoryTabIcon from "@/pages/workbench/components/WorkbenchSidebar/components/icons/HistoryTabIcon.tsx"
import ToolsTabIcon from "@/pages/workbench/components/WorkbenchSidebar/components/icons/ToolsTabIcon.tsx"
import VaultTabIcon from "@/pages/workbench/components/WorkbenchSidebar/components/icons/VaultTabIcon.tsx"
import cls from "@/pages/workbench/components/WorkbenchSidebar/WorkbenchSidebar.module.css"
import {cn} from "@/app/utils/cn.ts"

interface Props {
    history: WorkbenchHistory
    vault: VaultPaneBundle
    navOpen: boolean
    onToggleNav: () => void
}

const SEGMENTS = [
    {key: "history", label: "History", icon: <HistoryTabIcon/>},
    {key: "tools", label: "Tools", icon: <ToolsTabIcon/>, disabled: true, tooltip: "Coming soon"},
    {key: "vault", label: "Vault", icon: <VaultTabIcon/>},
]

// Unified persistent left pane for both api mode and running-docker mode. The
// collapse toggle rides at the top so it stays visible when the column shrinks to
// the narrow rail. History/Vault tabs each swap the pane body; Tools is a disabled
// placeholder.
export default function WorkbenchSidebar({history, vault, navOpen, onToggleNav}: Props) {
    const [activeTab, setActiveTab] = useState("history")

    return (
        <>
            <div
                className={cn(cls.Backdrop, navOpen && cls.BackdropVisible)}
                data-testid="sidebar-backdrop"
                onClick={onToggleNav}
            />
            <aside className={cn(cls.WorkbenchSidebarContainer, navOpen && cls.NavOpen)}>
                <div className={cls.TopRail}>
                    {navOpen && <SidebarBrand/>}
                    <SidebarToggleButton open={navOpen} onToggle={onToggleNav}/>
                </div>
                <div className={cls.TabsRow}>
                    <SegmentedControl
                        options={SEGMENTS} active={activeTab} onChange={setActiveTab} collapsed={!navOpen}
                    />
                </div>
                {navOpen && activeTab === "vault" && <VaultPane {...vault}/>}
                {navOpen && activeTab !== "vault" && <HistoryPane history={history}/>}
                {!navOpen && <div className={cls.CollapsedFiller}/>}
                <SidebarFooter collapsed={!navOpen}/>
            </aside>
        </>
    )
}
