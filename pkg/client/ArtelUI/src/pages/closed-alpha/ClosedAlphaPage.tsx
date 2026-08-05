import {useEffect} from "react"
import {useNavigate} from "react-router-dom"
import {Button} from "@vervstack/chures"

import cls from "@/pages/closed-alpha/ClosedAlphaPage.module.css"
import {Path} from "@/app/routing/Router.tsx"
import useUser from "@/hooks/user/User.ts"

export default function ClosedAlphaPage() {
    const navigate = useNavigate()
    const {logout, isAdmin, isEmailsEnabled, isNotesEnabled, isTaskTrackersEnabled} = useUser()

    useEffect(() => {
        // if (isAdmin || isEmailsEnabled || isNotesEnabled || isTaskTrackersEnabled) navigate(Path.HomePage)
    }, [isAdmin, isEmailsEnabled, isNotesEnabled, isTaskTrackersEnabled])

    function handleLogout() {
        logout()
        navigate(Path.InitPage)
    }

    return (
        <div className={cls.ClosedAlphaContainer}>
            <div className={cls.Card}>
                <div className={cls.Logo}>artel</div>
                <div className={cls.Body}>
                    <h2 className={cls.Title}>Service in closed alpha</h2>
                    <p className={cls.Sub}>Your account is not yet approved for access. Contact the team to request an
                        invite.</p>
                </div>
                <Button variant="secondary" onClick={handleLogout}>Log out</Button>
            </div>
        </div>
    )
}
