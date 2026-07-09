import TopbarBrand from "@/segments/Topbar/components/TopbarBrand/TopbarBrand.tsx"
import TopbarNav from "@/segments/Topbar/components/TopbarNav/TopbarNav.tsx"
import TopbarUserMenu from "@/segments/Topbar/components/TopbarUserMenu/TopbarUserMenu.tsx"
import TopbarThemeToggle from "@/segments/Topbar/components/TopbarThemeToggle/TopbarThemeToggle.tsx"
import cls from "@/segments/Topbar/Topbar.module.css"

export default function Topbar() {
    return (
        <header className={cls.Topbar}>
            <TopbarBrand />
            <TopbarNav />
            <div className={cls.ActionsWrapper}>
                <TopbarThemeToggle />
                <TopbarUserMenu />
            </div>
        </header>
    )
}
