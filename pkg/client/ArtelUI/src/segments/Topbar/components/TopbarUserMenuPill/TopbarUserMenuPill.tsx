import {AvatarFallbackIcon, Button, ChevronDownIcon} from "@vervstack/chures"

import {cn} from "@/app/utils/cn"

import cls from "./TopbarUserMenuPill.module.css"


interface TopbarUserMenuPillProps {
    menuOpen: boolean
    photoUrl?: string
    onClick: (e: React.MouseEvent) => void
}

export default function TopbarUserMenuPill({menuOpen, photoUrl, onClick}: TopbarUserMenuPillProps) {
    return (
        <Button
            variant="ghost"
            className={cls.UserPill}
            aria-expanded={menuOpen}
            aria-haspopup="menu"
            onClick={onClick}
        >
            <span className={cls.Avatar}>
                {photoUrl ? (
                    <img src={photoUrl} alt="avatar" width={28} height={28}
                         style={{borderRadius: '50%', objectFit: 'cover', display: 'block'}}/>
                ) : (
                    <AvatarFallbackIcon/>
                )}
            </span>
            <ChevronDownIcon className={cn(cls.Chevron, menuOpen && cls.ChevronOpen)}/>
        </Button>
    )
}
