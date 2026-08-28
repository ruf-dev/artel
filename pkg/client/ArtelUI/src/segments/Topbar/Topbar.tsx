import {useState} from "react"

import {useIsMobileNav} from "@/app/hooks/useIsMobileNav.ts"
import ThemeToggle from "@/components/ThemeToggle/ThemeToggle.tsx"
import TopbarBrand from "@/segments/Topbar/components/TopbarBrand/TopbarBrand.tsx"
import TopbarNav from "@/segments/Topbar/components/TopbarNav/TopbarNav.tsx"
import TopbarUserMenu from "@/segments/Topbar/components/TopbarUserMenu/TopbarUserMenu.tsx"
import TopbarMobileTrigger from "@/segments/Topbar/components/TopbarMobileTrigger/TopbarMobileTrigger.tsx"
import TopbarMobileDrawer from "@/segments/Topbar/components/TopbarMobileDrawer/TopbarMobileDrawer.tsx"
import cls from "@/segments/Topbar/Topbar.module.css"

export default function Topbar() {
    const isMobile = useIsMobileNav()
    const [drawerOpen, setDrawerOpen] = useState(false)

    return (
        <header className={cls.Topbar}>
            {!isMobile && (
                <div className={cls.LeadingWrapper}>
                    <TopbarBrand />
                </div>
            )}
            {!isMobile && <TopbarNav />}
            <div className={cls.ActionsWrapper}>
                <ThemeToggle />
                <TopbarUserMenu />
                {isMobile && <TopbarMobileTrigger open={drawerOpen} onClick={() => setDrawerOpen(true)}/>}
            </div>
            {isMobile && <TopbarMobileDrawer open={drawerOpen} onClose={() => setDrawerOpen(false)}/>}
        </header>
    )
}
