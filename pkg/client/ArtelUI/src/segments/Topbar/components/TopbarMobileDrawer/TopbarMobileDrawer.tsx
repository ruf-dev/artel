import {createPortal} from "react-dom"
import {NavLink} from "react-router-dom"

import {Path} from "@/app/routing/Router.tsx"
import {cn} from "@/app/utils/cn.ts"
import useUser from "@/hooks/user/User.ts"
import TopbarDrawerCloseButton from "@/segments/Topbar/components/TopbarDrawerCloseButton/TopbarDrawerCloseButton.tsx"
import cls from "@/segments/Topbar/components/TopbarMobileDrawer/TopbarMobileDrawer.module.css"

interface TopbarMobileDrawerProps {
    open: boolean
    onClose: () => void
}

export default function TopbarMobileDrawer({open, onClose}: TopbarMobileDrawerProps) {
    const {isNotesEnabled} = useUser()

    function linkClass({isActive}: {isActive: boolean}) {
        return cn(cls.NavLink, isActive && cls.NavLinkActive)
    }

    return createPortal(
        <>
            <div className={cn(cls.Backdrop, open && cls.BackdropOpen)} onClick={onClose}/>
            <div className={cn(cls.DrawerPanel, open && cls.DrawerPanelOpen)}>
                <div className={cls.DrawerHeader}>
                    <TopbarDrawerCloseButton onClose={onClose}/>
                </div>
                <nav className={cls.NavList} onClick={onClose}>
                    <NavLink to={Path.HomePage} end className={linkClass}>
                        Vaults
                    </NavLink>
                    {isNotesEnabled ? (
                        <NavLink to={Path.NotesPage} className={linkClass}>
                            Notes
                        </NavLink>
                    ) : (
                        <span className={cls.NavLinkDisabled}>Notes</span>
                    )}
                    <NavLink to={Path.ConnectionsPage} className={linkClass}>
                        Connections
                    </NavLink>
                    <NavLink to={Path.ToolboxPage} className={linkClass}>
                        Toolbox
                    </NavLink>
                    <NavLink to={Path.TractsPage} className={linkClass}>
                        Tracts
                    </NavLink>
                </nav>
            </div>
        </>,
        document.body
    )
}
