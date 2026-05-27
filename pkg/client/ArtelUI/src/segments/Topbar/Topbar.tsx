import cls from "./Topbar.module.css"
import TopbarBrand from "./TopbarBrand.tsx"
import TopbarNav from "./TopbarNav.tsx"
import TopbarUserMenu from "./TopbarUserMenu.tsx"

export default function Topbar() {
    return (
        <header className={cls.Topbar}>
            <TopbarBrand />
            <TopbarNav />
            <TopbarUserMenu />
        </header>
    )
}
