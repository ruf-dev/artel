import cls from "./TopbarNav.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {NavLink} from "react-router-dom"
import useUser from "@/hooks/user/User.ts"

export default function TopbarNav() {
    const {isEmailsEnabled, isTaskTrackersEnabled} = useUser()

    return (
        <nav className={cls.Nav}>
            <NavLink
                to={Path.HomePage}
                end
                className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
            >
                Vaults
            </NavLink>
            {isTaskTrackersEnabled ? (
                <NavLink
                    to={Path.TaskTrackersPage}
                    className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
                >
                    Tasks
                </NavLink>
            ) : (
                <div
                    data-tooltip-id={"root-tooltip"}
                    data-tooltip-content={"Task trackers are not enabled for your account"}
                >
                    <span className={cls.NavLinkDisabled}>Tasks</span>
                </div>
            )}
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
            <NavLink
                to={Path.NotesPage}
                className={({isActive}) => isActive ? `${cls.NavLink} ${cls.NavLinkActive}` : cls.NavLink}
            >
                Notes
            </NavLink>
        </nav>
    )
}
