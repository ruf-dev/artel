import {useState} from "react"
import {Button} from "@vervstack/chures"

import SegmentedControl from "@/components/atoms/SegmentedControl/SegmentedControl.tsx"
import PlusIcon from "@/icons/common/PlusIcon.tsx"
import {WorkbenchHistory} from "@/pages/workbench/processes/useWorkbenchHistory.ts"
import SidebarBrand from "@/pages/workbench/components/WorkbenchSidebar/components/SidebarBrand/SidebarBrand.tsx"
import HistoryPane from "@/pages/workbench/components/WorkbenchSidebar/components/HistoryPane/HistoryPane.tsx"
import SidebarFooter from "@/pages/workbench/components/WorkbenchSidebar/components/SidebarFooter/SidebarFooter.tsx"
import HistoryTabIcon from "@/pages/workbench/components/WorkbenchSidebar/components/icons/HistoryTabIcon.tsx"
import ToolsTabIcon from "@/pages/workbench/components/WorkbenchSidebar/components/icons/ToolsTabIcon.tsx"
import VaultTabIcon from "@/pages/workbench/components/WorkbenchSidebar/components/icons/VaultTabIcon.tsx"
import cls from "@/pages/workbench/components/WorkbenchSidebar/WorkbenchSidebar.module.css"

interface Props {
    history: WorkbenchHistory
    showCloseButton?: boolean
}

const SEGMENTS = [
    {key: "history", label: "History", icon: <HistoryTabIcon/>},
    {key: "tools", label: "Tools", icon: <ToolsTabIcon/>, disabled: true, tooltip: "Coming soon"},
    {key: "vault", label: "Vault", icon: <VaultTabIcon/>, disabled: true, tooltip: "Coming soon"},
]

// Unified persistent left pane for both api mode and running-docker mode (Stage 2).
// The exit-to-home control renders only in api mode and rides here temporarily —
// Stage 3 relocates it. Tools/Vault tabs are disabled placeholders.
export default function WorkbenchSidebar({history, showCloseButton}: Props) {
    const [activeTab, setActiveTab] = useState("history")

    return (
        <aside className={cls.WorkbenchSidebarContainer}>
            <SidebarBrand showClose={showCloseButton}/>
            <Button variant="secondary" className={cls.NewChatButton} onClick={history.onNewChat}>
                <PlusIcon/>
                New chat
            </Button>
            <div className={cls.TabsRow}>
                <SegmentedControl options={SEGMENTS} active={activeTab} onChange={setActiveTab}/>
            </div>
            <HistoryPane history={history}/>
            <SidebarFooter/>
        </aside>
    )
}
