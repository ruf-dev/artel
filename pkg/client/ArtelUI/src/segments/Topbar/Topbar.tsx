import cls from "./Topbar.module.css"
import TopbarBrand from "./components/TopbarBrand/TopbarBrand.tsx"
import TopbarNav from "./components/TopbarNav/TopbarNav.tsx"
import TopbarUserMenu from "./components/TopbarUserMenu/TopbarUserMenu.tsx"
import TopbarThemeToggle from "./components/TopbarThemeToggle/TopbarThemeToggle.tsx"

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
