import {useState} from "react"
import cls from "./TopbarUserMenu.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {useNavigate} from "react-router-dom"
import useUser from "@/hooks/user/User.ts"

interface UserPillProps {
    menuOpen: boolean
    photoUrl?: string
    onClick: (e: React.MouseEvent) => void
}

function UserPill({menuOpen, photoUrl, onClick}: UserPillProps) {
    return (
        <button
            className={cls.UserPill}
            type="button"
            aria-expanded={menuOpen}
            aria-haspopup="menu"
            onClick={onClick}
        >
            <span className={cls.Avatar}>
                {photoUrl ? (
                    <img src={photoUrl} alt="avatar" width={28} height={28}
                         style={{borderRadius: '50%', objectFit: 'cover', display: 'block'}}/>
                ) : (
                    <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="white" strokeWidth="2"
                         strokeLinecap="round" strokeLinejoin="round">
                        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                        <circle cx="12" cy="7" r="4"/>
                    </svg>
                )}
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

function ApiKeysItem({onClick}: {onClick: () => void}) {
    return (
        <button className={cls.MenuItem} role="menuitem" type="button" onClick={onClick}>
            <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
                 strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
            </svg>
            <span>API Keys</span>
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
    onApiKeys: () => void
    onLogout: () => void
}

function UserMenu({isAdmin, onAdmin, onApiKeys, onLogout}: UserMenuProps) {
    return (
        <div className={cls.Menu} role="menu">
            {isAdmin && <AdminPanelItem onClick={onAdmin}/>}
            <ApiKeysItem onClick={onApiKeys}/>
            <LogoutItem onClick={onLogout}/>
        </div>
    )
}

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
            <UserPill menuOpen={menuOpen} photoUrl={photoUrl} onClick={e => {
                e.stopPropagation();
                setMenuOpen(v => !v)
            }}/>
            {menuOpen && <div className={cls.Backdrop} onClick={() => setMenuOpen(false)}/>}
            {menuOpen && <UserMenu isAdmin={isAdmin} onAdmin={handleAdmin} onApiKeys={handleApiKeys} onLogout={handleLogout}/>}
        </div>
    )
}
