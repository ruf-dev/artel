import {useEffect, useState} from "react"
import {Button, Loader} from "@vervstack/chures"

import cls from "@/pages/workbench/components/LoginRelayScreen/LoginRelayScreen.module.css"
import Input from "@/components/atoms/Input/Input.tsx"
import {workbenchService} from "@/processes/Workbench.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

type LoginState = "pending" | "url_present" | "error" | "authorized"

const KNOWN_STATES = new Set<string>(["pending", "url_present", "error", "authorized"])

interface Props {
    vaultId: string
    onDone: () => void
}

export default function LoginRelayScreen({vaultId, onDone}: Props) {
    const bakeError = useBakeError()

    const [state, setState] = useState<LoginState>("pending")
    const [url, setUrl] = useState("")
    const [errorMessage, setErrorMessage] = useState("")
    const [code, setCode] = useState("")
    const [submitting, setSubmitting] = useState(false)

    useEffect(() => {
        const controller = new AbortController()

        workbenchService.watchLogin(vaultId, (s, u, e) => {
            setState(KNOWN_STATES.has(s) ? (s as LoginState) : "error")
            setUrl(u)
            setErrorMessage(e)
        }, controller.signal).catch((err: unknown) => {
            if (err instanceof DOMException && err.name === "AbortError") return
            bakeError("Lost connection to Workbench login", err)
        })

        return () => controller.abort()
    }, [vaultId])

    function handleSubmitCode() {
        setSubmitting(true)
        workbenchService.submitLoginCode(vaultId, code)
            .then(() => setCode(""))
            .catch(e => bakeError("Failed to submit login code", e))
            .finally(() => setSubmitting(false))
    }

    return (
        <div className={cls.LoginRelayScreenContainer}>
            {state === "pending" && (
                <>
                    <Loader variant="arcs" size="sm" color="var(--coral)"/>
                    <p className={cls.PendingText}>Starting Claude Code…</p>
                </>
            )}

            {(state === "url_present" || state === "error") && (
                <>
                    {url && (
                        <a className={cls.LoginLink} href={url} target="_blank" rel="noreferrer">
                            Open Claude sign-in page
                        </a>
                    )}
                    {state === "error" && errorMessage && <p className={cls.ErrorText}>{errorMessage}</p>}
                    <label className={cls.Field}>
                        <span className={cls.FieldLabel}>Verification code</span>
                        <Input value={code} setValue={setCode} placeholder="Paste the code here"
                            disabled={submitting}/>
                    </label>
                    <Button variant="primary" onClick={handleSubmitCode} disabled={submitting || !code}>
                        {submitting ? "Submitting…" : "Submit code"}
                    </Button>
                </>
            )}

            {state === "authorized" && (
                <>
                    <p className={cls.AuthorizedText}>Claude Code signed in successfully.</p>
                    <Button variant="primary" onClick={onDone}>Done</Button>
                </>
            )}
        </div>
    )
}
