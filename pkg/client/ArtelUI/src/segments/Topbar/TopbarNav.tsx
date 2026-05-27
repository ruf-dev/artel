import cls from "./TopbarNav.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {NavLink} from "react-router-dom"
import useUser from "@/hooks/user/User.ts"

export default function TopbarNav() {
    const {isEmailsEnabled} = useUser()

    return (
        <nav className={cls.Nav}>
            <NavLink
                to={Path.HomePage}
                end
                className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
            >
                Vaults
            </NavLink>
            {isEmailsEnabled ? (
                <NavLink
                    to={Path.EmailsPage}
                    className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
                >
                    Emails
                </NavLink>
            ) : (
                <div
                    data-tooltip-id={"root-tooltip"}
                    data-tooltip-content={"Emails are not enabled for your account"}
                >
                <span
                    className={cls.NavLinkDisabled}>Emails</span>
                </div>
            )}
        </nav>
    )
}
