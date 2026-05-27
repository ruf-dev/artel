import cls from "./TopbarNav.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {NavLink} from "react-router-dom"

export default function TopbarNav() {
    return (
        <nav className={cls.Nav}>
            <NavLink
                to={Path.HomePage}
                end
                className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
            >
                Vaults
            </NavLink>
            <NavLink
                to={Path.EmailsPage}
                className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
            >
                Emails
            </NavLink>
        </nav>
    )
}
