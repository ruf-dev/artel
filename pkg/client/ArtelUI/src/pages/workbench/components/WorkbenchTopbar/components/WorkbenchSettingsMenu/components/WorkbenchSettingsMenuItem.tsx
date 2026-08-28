import {ReactNode} from "react"
import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn"
// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import cls from "@/pages/workbench/components/WorkbenchTopbar/components/WorkbenchSettingsMenu/components/WorkbenchSettingsMenuItem.module.css"

interface WorkbenchSettingsMenuItemProps {
    icon: ReactNode
    label: string
    onClick: () => void
    danger?: boolean
}

export default function WorkbenchSettingsMenuItem({icon, label, onClick, danger}: WorkbenchSettingsMenuItemProps) {
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
