import {createPortal} from "react-dom"
import {useNavigate} from "react-router-dom"

import {Path} from "@/app/routing/Router.tsx"
import {cn} from "@/app/utils/cn.ts"
import useUser from "@/hooks/user/User.ts"
import ConnectionsIcon from "@/segments/Topbar/components/icons/ConnectionsIcon.tsx"
import NotesIcon from "@/segments/Topbar/components/icons/NotesIcon.tsx"
import ToolboxIcon from "@/segments/Topbar/components/icons/ToolboxIcon.tsx"
import TractsIcon from "@/segments/Topbar/components/icons/TractsIcon.tsx"
import VaultsIcon from "@/segments/Topbar/components/icons/VaultsIcon.tsx"
import TopbarBrand from "@/segments/Topbar/components/TopbarBrand/TopbarBrand.tsx"
import TopbarDrawerCloseButton from "@/segments/Topbar/components/TopbarDrawerCloseButton/TopbarDrawerCloseButton.tsx"
import TopbarDrawerLogoutButton
    from "@/segments/Topbar/components/TopbarDrawerLogoutButton/TopbarDrawerLogoutButton.tsx"
import cls from "@/segments/Topbar/components/TopbarMobileDrawer/TopbarMobileDrawer.module.css"
import TopbarMobileNavLink from "@/segments/Topbar/components/TopbarMobileNavLink/TopbarMobileNavLink.tsx"

interface TopbarMobileDrawerProps {
    open: boolean
    onClose: () => void
}

export default function TopbarMobileDrawer({open, onClose}: TopbarMobileDrawerProps) {
    const {isNotesEnabled, logout} = useUser()
    const navigate = useNavigate()

    function handleLogout() {
        onClose()
        logout()
        navigate(Path.InitPage)
    }

    return createPortal(
        <>
            <div className={cn(cls.Backdrop, open && cls.BackdropOpen)} onClick={onClose}/>
            <div className={cn(cls.DrawerPanel, open && cls.DrawerPanelOpen)}>
                <div className={cls.DrawerHeader}>
                    <TopbarBrand/>
                    <TopbarDrawerCloseButton onClose={onClose}/>
                </div>
                <nav className={cls.NavList} onClick={onClose}>
                    <TopbarMobileNavLink to={Path.HomePage} end Icon={VaultsIcon} label="Vaults"/>
                    <TopbarMobileNavLink
                        to={Path.NotesPage}
                        Icon={NotesIcon}
                        label="Notes"
                        disabled={!isNotesEnabled}
                    />
                    <TopbarMobileNavLink to={Path.ConnectionsPage} Icon={ConnectionsIcon} label="Connections"/>
                    <TopbarMobileNavLink to={Path.ToolboxPage} Icon={ToolboxIcon} label="Toolbox"/>
                    <TopbarMobileNavLink to={Path.TractsPage} Icon={TractsIcon} label="Tracts"/>
                </nav>
                <div className={cls.DrawerFooter}>
                    <TopbarDrawerLogoutButton onLogout={handleLogout}/>
                </div>
            </div>
        </>,
        document.body
    )
}
