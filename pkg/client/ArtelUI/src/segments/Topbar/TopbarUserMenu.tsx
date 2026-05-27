import {useState} from "react"
import cls from "./TopbarUserMenu.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {useNavigate} from "react-router-dom"
import useUser from "@/hooks/user/User.ts"

interface UserPillProps {
    menuOpen: boolean
    onClick: (e: React.MouseEvent) => void
}

function UserPill({menuOpen, onClick}: UserPillProps) {
    return (
        <button
            className={cls.UserPill}
            type="button"
            aria-expanded={menuOpen}
            aria-haspopup="menu"
            onClick={onClick}
        >
            <span className={cls.Avatar}>
                <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="white" strokeWidth="2"
                     strokeLinecap="round" strokeLinejoin="round">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                    <circle cx="12" cy="7" r="4"/>
                </svg>
            </span>
            <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor"
                 strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"
                 style={{
                     transform: menuOpen ? 'rotate(180deg)' : 'none',
                     transition: 'transform 200ms ease',
                     color: 'rgba(255,255,255,0.4)'
                 }}>
                <polyline points="6 9 12 15 18 9"/>
            </svg>
        </button>
    )
}

interface AdminPanelItemProps {
    onClick: () => void
}

function AdminPanelItem({onClick}: AdminPanelItemProps) {
    return (
        <button className={cls.MenuItem} role="menuitem" type="button" onClick={onClick}>
            <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
                 strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                <circle cx="12" cy="12" r="3"/>
                <path d="M19.07 4.93a10 10 0 0 1 0 14.14M4.93 4.93a10 10 0 0 0 0 14.14"/>
            </svg>
            <span>Admin Panel</span>
        </button>
    )
}

interface LogoutItemProps {
    onClick: () => void
}

function LogoutItem({onClick}: LogoutItemProps) {
    return (
        <button className={`${cls.MenuItem} ${cls.MenuItemDanger}`} role="menuitem" type="button" onClick={onClick}>
            <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
                 strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                <polyline points="16 17 21 12 16 7"/>
                <line x1="21" y1="12" x2="9" y2="12"/>
            </svg>
            <span>Log out</span>
        </button>
    )
}

interface UserMenuProps {
    isAdmin: boolean
    onAdmin: () => void
    onLogout: () => void
}

function UserMenu({isAdmin, onAdmin, onLogout}: UserMenuProps) {
    return (
        <div className={cls.Menu} role="menu">
            {isAdmin && <AdminPanelItem onClick={onAdmin}/>}
            <LogoutItem onClick={onLogout}/>
        </div>
    )
}

export default function TopbarUserMenu() {
    const [menuOpen, setMenuOpen] = useState(false)
    const navigate = useNavigate()
    const {isAdmin, logout} = useUser()

    function handleLogout() {
        setMenuOpen(false)
        logout()
        navigate(Path.InitPage)
    }

    function handleAdmin() {
        setMenuOpen(false)
        navigate(Path.Admin)
    }
    
    return (
        <div className={cls.UserWrap}>
            <UserPill menuOpen={menuOpen} onClick={e => {
                e.stopPropagation();
                setMenuOpen(v => !v)
            }}/>
            {menuOpen && <div className={cls.Backdrop} onClick={() => setMenuOpen(false)}/>}
            {menuOpen && <UserMenu isAdmin={isAdmin} onAdmin={handleAdmin} onLogout={handleLogout}/>}
        </div>
    )
}
