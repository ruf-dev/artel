import {useState} from "react"

import {useIsMobileNav} from "@/app/hooks/useIsMobileNav.ts"
import TopbarBrand from "@/segments/Topbar/components/TopbarBrand/TopbarBrand.tsx"
import TopbarNav from "@/segments/Topbar/components/TopbarNav/TopbarNav.tsx"
import TopbarUserMenu from "@/segments/Topbar/components/TopbarUserMenu/TopbarUserMenu.tsx"
import TopbarThemeToggle from "@/segments/Topbar/components/TopbarThemeToggle/TopbarThemeToggle.tsx"
import TopbarMobileTrigger from "@/segments/Topbar/components/TopbarMobileTrigger/TopbarMobileTrigger.tsx"
import TopbarMobileDrawer from "@/segments/Topbar/components/TopbarMobileDrawer/TopbarMobileDrawer.tsx"
import cls from "@/segments/Topbar/Topbar.module.css"

export default function Topbar() {
    const isMobile = useIsMobileNav()
    const [drawerOpen, setDrawerOpen] = useState(false)

    return (
        <header className={cls.Topbar}>
            <div className={cls.LeadingWrapper}>
                {isMobile && <TopbarMobileTrigger open={drawerOpen} onClick={() => setDrawerOpen(true)}/>}
                <TopbarBrand />
            </div>
            {!isMobile && <TopbarNav />}
            <div className={cls.ActionsWrapper}>
                <TopbarThemeToggle />
                <TopbarUserMenu />
            </div>
            {isMobile && <TopbarMobileDrawer open={drawerOpen} onClose={() => setDrawerOpen(false)}/>}
        </header>
    )
}
