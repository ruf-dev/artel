import {useEffect, useState} from "react"
import {useNavigate, useParams} from "react-router-dom"
import {Button, Loader} from "@vervstack/chures"

import cls from "@/pages/join/JoinVaultPage.module.css"
import {useVaultMutations} from "@/app/hooks/Vaults.ts"
import {Path} from "@/app/routing/Router.tsx"

export default function JoinVaultPage() {
    const {token} = useParams<{token: string}>()
    const navigate = useNavigate()
    const {acceptInvite} = useVaultMutations()
    const [status, setStatus] = useState<"joining" | "success" | "error">("joining")
    const [errorMsg, setErrorMsg] = useState("")

    useEffect(() => {
        if (!token) {
            navigate(Path.HomePage, {replace: true})
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
    }, [token])

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
                        <Button variant="secondary" onClick={() => navigate(Path.HomePage)}>
                            Go home
                        </Button>
                    </>
                )}
            </div>
        </div>
    )
}
