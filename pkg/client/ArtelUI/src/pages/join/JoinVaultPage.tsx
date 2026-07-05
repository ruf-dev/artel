import {useEffect, useState} from "react"
import {useNavigate, useParams} from "react-router-dom"
import {Loader} from "@vervstack/chures"

import cls from "@/pages/join/JoinVaultPage.module.css"

import {useVaultMutations} from "@/app/hooks/Vaults.ts"
import useUser from "@/hooks/user/User.ts"
import {Path, REDIRECT_AFTER_LOGIN_KEY} from "@/app/routing/Router.tsx"

export default function JoinVaultPage() {
    const {token} = useParams<{token: string}>()
    const navigate = useNavigate()
    const {auth} = useUser()
    const {acceptInvite} = useVaultMutations()
    const [status, setStatus] = useState<"joining" | "success" | "error">("joining")
    const [errorMsg, setErrorMsg] = useState("")

    useEffect(() => {
        if (!token) {
            navigate(Path.HomePage, {replace: true})
            return
        }

        if (!auth.isAuthenticated()) {
            localStorage.setItem(REDIRECT_AFTER_LOGIN_KEY, `/join/${token}`)
            navigate(Path.InitPage, {replace: true})
            return
        }

        acceptInvite(token)
            .then(() => {
                setStatus("success")
                setTimeout(() => navigate(Path.HomePage, {replace: true}), 1500)
            })
            .catch((err: unknown) => {
                setStatus("error")
                setErrorMsg(err instanceof Error ? err.message : "Failed to join vault")
            })
    }, [token, auth])

    return (
        <div className={cls.JoinContainer}>
            <div className={cls.Card}>
                {status === "joining" && (
                    <>
                        <Loader variant="arcs" size="sm" color="var(--coral)"/>
                        <p className={cls.Message}>Joining vault…</p>
                    </>
                )}
                {status === "success" && (
                    <>
                        <div className={cls.SuccessIcon}>✓</div>
                        <p className={cls.Message}>Joined! Redirecting…</p>
                    </>
                )}
                {status === "error" && (
                    <>
                        <p className={cls.ErrorMessage}>{errorMsg}</p>
                        <button className={cls.BtnHome} onClick={() => navigate(Path.HomePage)}>
                            Go home
                        </button>
                    </>
                )}
            </div>
        </div>
    )
}
