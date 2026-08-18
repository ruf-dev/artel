import {useEffect, useState} from "react"
import {Button, Loader} from "@vervstack/chures"

import cls from "@/pages/workbench/components/Terminal/Terminal.module.css"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

interface Props {
    vaultId: string
}

type TerminalState = "loading" | "ok" | "error"

export default function Terminal({vaultId}: Props) {
    const [state, setState] = useState<TerminalState>("loading")
    const [retryToken, setRetryToken] = useState(0)
    const bakeError = useBakeError()

    const url = `/api/vaults/workbench/${vaultId}/terminal/`

    useEffect(() => {
        let ignore = false

        fetch(url)
            .then(res => {
                if (ignore) return
                if (res.ok) {
                    setState("ok")
                } else {
                    return res.text().then(body => {
                        if (ignore) return
                        bakeError("Failed to connect to workbench terminal", new Error(body || `status ${res.status}`))
                        setState("error")
                    })
                }
            })
            .catch(err => {
                if (ignore) return
                bakeError("Failed to connect to workbench terminal", err)
                setState("error")
            })

        return () => {
            ignore = true
        }
    }, [vaultId, retryToken, url, bakeError])

    function handleRetry() {
        setState("loading")
        setRetryToken(t => t + 1)
    }

    return (
        <div className={cls.TerminalContainer}>
            {state === "loading" && (
                <Loader variant="arcs" size="sm" color="var(--coral)"/>
            )}
            {state === "error" && (
                <div className={cls.ErrorState}>
                    <p>Couldn&apos;t connect to the workbench terminal.</p>
                    <Button variant="primary" onClick={handleRetry}>Retry</Button>
                </div>
            )}
            {state === "ok" && (
                <iframe
                    className={cls.Iframe}
                    src={url}
                    title="Workbench terminal"
                    sandbox="allow-scripts allow-same-origin"
                />
            )}
        </div>
    )
}
