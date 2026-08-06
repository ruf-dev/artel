import {useState} from "react"
import {Button} from "@vervstack/chures"

import {CheckTelegramConnectionRequest, ExternalConnectionsAPI} from "@/app/api/artel/external_connections.pb.ts"
import {grpcErrorMessage, isGrpcError} from "@/processes/grpcErrors.ts"
import useUser from "@/hooks/user/User.ts"
import cls from "@/dialogs/ManageTelegramDialog/components/TelegramCheckButton/TelegramCheckButton.module.css"

export type CheckStatus = "idle" | "checking" | "ok" | "fail"

interface TelegramCheckButtonProps {
    req: CheckTelegramConnectionRequest
    disabled?: boolean
    onResult: (verified: boolean) => void
}

export default function TelegramCheckButton({req, disabled, onResult}: TelegramCheckButtonProps) {
    const [status, setStatus] = useState<CheckStatus>("idle")
    const [botUsername, setBotUsername] = useState("")
    const [errorMessage, setErrorMessage] = useState<string | null>(null)
    const {auth} = useUser()

    function handleCheck() {
        setStatus("checking")
        setErrorMessage(null)
        ExternalConnectionsAPI.CheckTelegramConnection(req, auth.getInitReq())
            .then(resp => {
                setStatus("ok")
                setBotUsername(resp.botUsername ?? "")
                onResult(true)
            })
            .catch(err => {
                setStatus("fail")
                setErrorMessage(isGrpcError(err) ? grpcErrorMessage(err) : "Connection check failed")
                onResult(false)
            })
    }

    return (
        <div className={cls.TelegramCheckButtonContainer}>
            {status === "ok" && <span className={cls.BadgeOk}>Connected as @{botUsername}</span>}
            {status === "fail" && <span className={cls.BadgeFail}>Failed</span>}
            {status === "fail" && errorMessage && (
                <span className={cls.ErrorText} title={errorMessage}>{errorMessage}</span>
            )}
            <Button variant="secondary" onClick={handleCheck} disabled={disabled || status === "checking"}>
                {status === "checking" ? "Testing…" : "Test"}
            </Button>
        </div>
    )
}
