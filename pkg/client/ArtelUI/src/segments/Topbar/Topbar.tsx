import cls from "./Topbar.module.css"
import TopbarBrand from "./TopbarBrand.tsx"
import TopbarNav from "./TopbarNav.tsx"
import TopbarUserMenu from "./TopbarUserMenu.tsx"
import TopbarThemeToggle from "./TopbarThemeToggle.tsx"

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
