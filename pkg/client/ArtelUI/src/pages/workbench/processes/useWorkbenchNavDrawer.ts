import {useEffect, useMemo, useState} from "react"

import {useIsMobileNav} from "@/app/hooks/useIsMobileNav.ts"
import {WorkbenchHistory} from "@/pages/workbench/processes/useWorkbenchHistory.ts"

// Owns the workbench sidebar's open/closed state. On desktop the sidebar is an
// in-flow column that starts open; on mobile it's a fully off-canvas drawer that
// starts closed, and selecting/creating a chat also closes it. Keeping this out of
// WorkbenchPage keeps that component under the max-lines lint budget.
export function useWorkbenchNavDrawer(history: WorkbenchHistory) {
    const isMobile = useIsMobileNav()
    const [navOpen, setNavOpen] = useState(!isMobile)

    useEffect(() => {
        setNavOpen(!isMobile)
    }, [isMobile])

    function toggleNav() {
        setNavOpen(v => !v)
    }

    const sidebarHistory = useMemo(() => (isMobile ? {
        ...history,
        onSelect: (id: string) => {
            history.onSelect(id)
            setNavOpen(false)
        },
        onNewChat: () => {
            history.onNewChat()
            setNavOpen(false)
        },
    } : history), [isMobile, history])

    return {navOpen, toggleNav, sidebarHistory}
}
