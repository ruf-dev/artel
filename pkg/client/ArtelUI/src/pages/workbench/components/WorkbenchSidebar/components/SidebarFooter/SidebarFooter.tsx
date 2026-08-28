import ThemeToggle from "@/components/ThemeToggle/ThemeToggle.tsx"
import {cn} from "@/app/utils/cn.ts"
import cls from "@/pages/workbench/components/WorkbenchSidebar/components/SidebarFooter/SidebarFooter.module.css"

interface Props {
    collapsed?: boolean
}

// The sidebar's pinned footer — the shared dark/light theme toggle. Right-aligned
// when expanded, centered in the narrow rail when collapsed.
export default function SidebarFooter({collapsed}: Props) {
    return (
        <div className={cn(cls.SidebarFooterContainer, collapsed && cls.Collapsed)}>
            <ThemeToggle/>
        </div>
    )
}
