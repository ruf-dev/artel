import {ReactNode} from "react"
import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn"
import cls from "@/segments/Topbar/components/TopbarUserMenuItem/TopbarUserMenuItem.module.css"

interface TopbarUserMenuItemProps {
    icon: ReactNode
    label: string
    onClick: () => void
    danger?: boolean
}

export default function TopbarUserMenuItem({icon, label, onClick, danger}: TopbarUserMenuItemProps) {
    return (
        <Button
            variant="ghost"
            className={cn(cls.MenuItem, danger && cls.MenuItemDanger)}
            role="menuitem"
            onClick={onClick}
        >
            {icon}
            <span>{label}</span>
        </Button>
    )
}
