import {useState} from "react"
import cls from "./Topbar.module.css"
import {Path} from "@/app/routing/Router.tsx";
import {useNavigate, NavLink} from "react-router-dom";
import useUser from "@/hooks/user/User.ts";


export default function Topbar() {
    const [menuOpen, setMenuOpen] = useState(false)
    const navigate = useNavigate();
    const {isAdmin, logout} = useUser();

    function handleLogout() {
        logout()
        navigate(Path.InitPage)
    }

    return (
        <header className={cls.Topbar}>
            <a className={cls.Brand} href={Path.HomePage} aria-label="Artel home">
                <svg className={cls.BrandMark} viewBox="0 0 100 100" aria-hidden="true">
                    <defs>
                        <mask id="brand-a-cut">
                            <rect width="100" height="100" fill="white"/>
                            <path
                                d="M 50 22 L 28 78 L 38 78 L 43.5 64 L 56.5 64 L 62 78 L 72 78 Z M 46.5 56 L 50 38 L 53.5 56 Z"
                                fill="black"/>
                        </mask>
                    </defs>
                    <circle cx="50" cy="50" r="50" fill="#FF4B3E" mask="url(#brand-a-cut)"/>
                </svg>
                <span className={cls.BrandWord}>artel</span>
            </a>

            <nav className={cls.Nav}>
                <NavLink
                    to={Path.HomePage}
                    end
                    className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
                >
                    Vaults
                </NavLink>
                <NavLink
                    to={Path.EmailsPage}
                    className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
                >
                    Emails
                </NavLink>
            </nav>

            <div className={cls.Right}>
                <div className={cls.UserWrap}>
                    <button
                        className={cls.UserPill}
                        type="button"
                        aria-expanded={menuOpen}
                        aria-haspopup="menu"
                        onClick={e => {
                            e.stopPropagation();
                            setMenuOpen(v => !v)
                        }}
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

                    {menuOpen && (
                        <>
                            <div className={cls.Backdrop} onClick={() => setMenuOpen(false)}/>
                            <div className={cls.Menu} role="menu">
                                {isAdmin && (
                                    <button
                                        className={cls.MenuItem}
                                        role="menuitem"
                                        type="button"
                                        onClick={() => {
                                            setMenuOpen(false);
                                            navigate(Path.Admin)
                                        }}
                                    >
                                        <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
                                             strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                                            <circle cx="12" cy="12" r="3"/>
                                            <path d="M19.07 4.93a10 10 0 0 1 0 14.14M4.93 4.93a10 10 0 0 0 0 14.14"/>
                                        </svg>
                                        <span>Admin Panel</span>
                                    </button>
                                )}
                                <button
                                    className={`${cls.MenuItem} ${cls.MenuItemDanger}`}
                                    role="menuitem"
                                    type="button"
                                    onClick={() => {
                                        setMenuOpen(false);
                                        handleLogout()
                                    }}
                                >
                                    <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
                                         strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                                        <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                                        <polyline points="16 17 21 12 16 7"/>
                                        <line x1="21" y1="12" x2="9" y2="12"/>
                                    </svg>
                                    <span>Log out</span>
                                </button>
                            </div>
                        </>
                    )}
                </div>
            </div>
        </header>
    )
}
