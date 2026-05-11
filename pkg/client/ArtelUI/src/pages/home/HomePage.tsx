import {useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/home/HomePage.module.css"

import {AuthMiddleware} from "@/processes/Auth.ts"
import {Path} from "@/app/routing/Router.tsx"

interface Props {
    auth: AuthMiddleware
    onLogout: () => void
}

export default function HomePage({auth, onLogout}: Props) {
    const navigate = useNavigate()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    function handleLogout() {
        onLogout()
        navigate(Path.InitPage)
    }

    return (
        <div className={cls.Root}>
            <header className={cls.Header}>
                <span className={cls.Logo}>artel</span>
                <button className={cls.LogoutBtn} onClick={handleLogout}>log out</button>
            </header>

            <div className={cls.Content}>
                <p>Your vaults will appear here.</p>
            </div>
        </div>
    )
}
