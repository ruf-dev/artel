import {useEffect, useState} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/closed-alpha/ClosedAlphaPage.module.css"
import {Path} from "@/app/routing/Router.tsx"
import useUser from "@/hooks/user/User.ts"
import {useVaults} from "@/app/hooks/Vaults.ts"

export default function ClosedAlphaPage() {
    const navigate = useNavigate()
    const {logout} = useUser()
    const {fetch: fetchVaults, forbidden} = useVaults()
    const [fetched, setFetched] = useState(false)

    useEffect(() => {
        fetchVaults().then(() => setFetched(true))
    }, [fetchVaults])

    useEffect(() => {
        if (fetched && !forbidden) navigate(Path.HomePage)
    }, [fetched, forbidden, navigate])

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
                    <p className={cls.Sub}>Your account is not yet approved for access. Contact the team to request an invite.</p>
                </div>
                <button className={cls.LogoutBtn} onClick={handleLogout}>Log out</button>
            </div>
        </div>
    )
}
