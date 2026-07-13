import {useState} from "react"
import {Button} from "@vervstack/chures"

import {CheckTrelloConnectionRequest, ExternalConnectionsAPI} from "@/app/api/artel/external_connections.pb.ts"
import {grpcErrorMessage, isGrpcError} from "@/processes/grpcErrors.ts"
import useUser from "@/hooks/user/User.ts"
import cls from "@/dialogs/ManageTrelloDialog/components/TrelloCheckButton/TrelloCheckButton.module.css"

export type CheckStatus = "idle" | "checking" | "ok" | "fail"

interface TrelloCheckButtonProps {
    req: CheckTrelloConnectionRequest
    disabled?: boolean
    onResult: (verified: boolean) => void
}

export default function TrelloCheckButton({req, disabled, onResult}: TrelloCheckButtonProps) {
    const [status, setStatus] = useState<CheckStatus>("idle")
    const [fullName, setFullName] = useState("")
    const [errorMessage, setErrorMessage] = useState<string | null>(null)
    const {auth} = useUser()

    function handleCheck() {
        setStatus("checking")
        setErrorMessage(null)
        ExternalConnectionsAPI.CheckTrelloConnection(req, auth.getInitReq())
            .then(resp => {
                setStatus("ok")
                setFullName(resp.fullName ?? "")
                onResult(true)
            })
            .catch(err => {
                setStatus("fail")
                setErrorMessage(isGrpcError(err) ? grpcErrorMessage(err) : "Connection check failed")
                onResult(false)
            })
    }

    return (
        <div className={cls.TrelloCheckButtonContainer}>
            {status === "ok" && <span className={cls.BadgeOk}>Connected as {fullName}</span>}
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
