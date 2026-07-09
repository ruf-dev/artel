import {ReactNode} from "react"
import {Button} from "@vervstack/chures"

import cls from "./TopbarUserMenuItem.module.css"

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
            className={danger ? `${cls.MenuItem} ${cls.MenuItemDanger}` : cls.MenuItem}
            role="menuitem"
            onClick={onClick}
        >
            {icon}
            <span>{label}</span>
        </Button>
    )
}
