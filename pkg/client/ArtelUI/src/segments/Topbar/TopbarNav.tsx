import cls from "./TopbarNav.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {NavLink} from "react-router-dom"
import useUser from "@/hooks/user/User.ts"

export default function TopbarNav() {
    const {isNotesEnabled} = useUser()

    return (
        <nav className={cls.Nav}>
            <NavLink
                to={Path.HomePage}
                end
                className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
            >
                Vaults
            </NavLink>
            {isNotesEnabled ? (
                <NavLink
                    to={Path.NotesPage}
                    className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
                >
                    Notes
                </NavLink>
            ) : (
                <div
                    data-tooltip-id={"root-tooltip"}
                    data-tooltip-content={"Notes are not enabled for your account"}
                >
                    <span className={cls.NavLinkDisabled}>Notes</span>
                </div>
            )}
            <NavLink
                to={Path.ConnectionsPage}
                className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
            >
                Connections
            </NavLink>
            <NavLink
                to={Path.ToolboxPage}
                className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
            >
                Toolbox
            </NavLink>
        </nav>
    )
}
