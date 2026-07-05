import {ReactNode} from "react"
import cls from "./TopbarUserMenuItem.module.css"

interface TopbarUserMenuItemProps {
    icon: ReactNode
    label: string
    onClick: () => void
    danger?: boolean
}

export default function TopbarUserMenuItem({icon, label, onClick, danger}: TopbarUserMenuItemProps) {
    return (
        <button
            className={danger ? `${cls.MenuItem} ${cls.MenuItemDanger}` : cls.MenuItem}
            role="menuitem"
            type="button"
            onClick={onClick}
        >
            {icon}
            <span>{label}</span>
        </button>
    )
}
