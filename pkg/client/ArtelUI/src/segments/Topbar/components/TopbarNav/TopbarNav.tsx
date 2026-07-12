import {NavLink} from "react-router-dom"

import {Path} from "@/app/routing/Router.tsx"
import useUser from "@/hooks/user/User.ts"
import ConnectionsIcon from "@/segments/Topbar/components/icons/ConnectionsIcon.tsx"
import NotesIcon from "@/segments/Topbar/components/icons/NotesIcon.tsx"
import ToolboxIcon from "@/segments/Topbar/components/icons/ToolboxIcon.tsx"
import TractsIcon from "@/segments/Topbar/components/icons/TractsIcon.tsx"
import VaultsIcon from "@/segments/Topbar/components/icons/VaultsIcon.tsx"
import cls from "@/segments/Topbar/components/TopbarNav/TopbarNav.module.css"

export default function TopbarNav() {
    const {isNotesEnabled} = useUser()

    return (
        <nav className={cls.Nav}>
            <NavLink
                to={Path.HomePage}
                end
                className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
            >
                <VaultsIcon className={cls.NavIcon}/>
                Vaults
            </NavLink>
            {isNotesEnabled ? (
                <NavLink
                    to={Path.NotesPage}
                    className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
                >
                    <NotesIcon className={cls.NavIcon}/>
                    Notes
                </NavLink>
            ) : (
                <span
                    className={cls.NavLinkDisabled}
                    data-tooltip-id={"root-tooltip"}
                    data-tooltip-content={"Notes are not enabled for your account"}
                >
                    <NotesIcon className={cls.NavIcon}/>
                    Notes
                </span>
            )}
            <NavLink
                to={Path.ConnectionsPage}
                className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
            >
                <ConnectionsIcon className={cls.NavIcon}/>
                Connections
            </NavLink>
            <NavLink
                to={Path.ToolboxPage}
                className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
            >
                <ToolboxIcon className={cls.NavIcon}/>
                Toolbox
            </NavLink>
            <NavLink
                to={Path.TractsPage}
                className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
            >
                <TractsIcon className={cls.NavIcon}/>
                Tracts
            </NavLink>
        </nav>
    )
}
