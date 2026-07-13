import {useState} from "react"
import {Button} from "@vervstack/chures"

import {CheckGitlabConnectionRequest, ExternalConnectionsAPI} from "@/app/api/artel/external_connections.pb.ts"
import {grpcErrorMessage, isGrpcError} from "@/processes/grpcErrors.ts"
import useUser from "@/hooks/user/User.ts"
import cls from "@/dialogs/ManageGitlabDialog/components/GitlabCheckButton/GitlabCheckButton.module.css"

export type CheckStatus = "idle" | "checking" | "ok" | "fail"

interface GitlabCheckButtonProps {
    req: CheckGitlabConnectionRequest
    disabled?: boolean
    onResult: (verified: boolean) => void
}

export default function GitlabCheckButton({req, disabled, onResult}: GitlabCheckButtonProps) {
    const [status, setStatus] = useState<CheckStatus>("idle")
    const [username, setUsername] = useState("")
    const [errorMessage, setErrorMessage] = useState<string | null>(null)
    const {auth} = useUser()

    function handleCheck() {
        setStatus("checking")
        setErrorMessage(null)
        ExternalConnectionsAPI.CheckGitlabConnection(req, auth.getInitReq())
            .then(resp => {
                setStatus("ok")
                setUsername(resp.username ?? "")
                onResult(true)
            })
            .catch(err => {
                setStatus("fail")
                setErrorMessage(isGrpcError(err) ? grpcErrorMessage(err) : "Connection check failed")
                onResult(false)
            })
    }

    return (
        <div className={cls.GitlabCheckButtonContainer}>
            {status === "ok" && <span className={cls.BadgeOk}>Connected as {username}</span>}
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
