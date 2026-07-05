import {useState} from "react"
import cls from "./TopbarUserMenu.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {useNavigate} from "react-router-dom"
import useUser from "@/hooks/user/User.ts"
import TopbarUserMenuPill from "./TopbarUserMenuPill.tsx"
import TopbarUserMenuList from "./TopbarUserMenuList.tsx"

export default function TopbarUserMenu() {
    const [menuOpen, setMenuOpen] = useState(false)
    const navigate = useNavigate()
    const {isAdmin, logout, photoUrl} = useUser()

    function handleLogout() {
        setMenuOpen(false)
        logout()
        navigate(Path.InitPage)
    }

    function handleAdmin() {
        setMenuOpen(false)
        navigate(Path.Admin)
    }

    function handleApiKeys() {
        setMenuOpen(false)
        navigate(Path.McpKeysPage)
    }

    return (
        <div className={cls.UserWrap}>
            <TopbarUserMenuPill menuOpen={menuOpen} photoUrl={photoUrl} onClick={e => {
                e.stopPropagation();
                setMenuOpen(v => !v)
            }}/>
            {menuOpen && <div className={cls.Backdrop} onClick={() => setMenuOpen(false)}/>}
            {menuOpen && <TopbarUserMenuList isAdmin={isAdmin} onAdmin={handleAdmin} onApiKeys={handleApiKeys} onLogout={handleLogout}/>}
        </div>
    )
}
