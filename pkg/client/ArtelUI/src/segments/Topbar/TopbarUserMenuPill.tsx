import cls from "./TopbarUserMenuPill.module.css"

interface TopbarUserMenuPillProps {
    menuOpen: boolean
    photoUrl?: string
    onClick: (e: React.MouseEvent) => void
}

export default function TopbarUserMenuPill({menuOpen, photoUrl, onClick}: TopbarUserMenuPillProps) {
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
