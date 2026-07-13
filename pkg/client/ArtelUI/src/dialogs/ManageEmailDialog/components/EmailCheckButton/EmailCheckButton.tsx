import {useState} from "react"
import {Button} from "@vervstack/chures"

import {CheckEmailConnectionRequest, ExternalConnectionsAPI} from "@/app/api/artel/external_connections.pb.ts"
import {grpcErrorMessage, isGrpcError} from "@/processes/grpcErrors.ts"
import useUser from "@/hooks/user/User.ts"
import cls from "@/dialogs/ManageEmailDialog/components/EmailCheckButton/EmailCheckButton.module.css"

export type CheckStatus = "idle" | "checking" | "ok" | "fail"

interface EmailCheckButtonProps {
    req: CheckEmailConnectionRequest
    disabled?: boolean
}

export default function EmailCheckButton({req, disabled}: EmailCheckButtonProps) {
    const [status, setStatus] = useState<CheckStatus>("idle")
    const [errorMessage, setErrorMessage] = useState<string | null>(null)
    const {auth} = useUser()

    function handleCheck() {
        setStatus("checking")
        setErrorMessage(null)
        ExternalConnectionsAPI.CheckEmailConnection(req, auth.getInitReq())
            .then(() => setStatus("ok"))
            .catch(err => {
                setStatus("fail")
                setErrorMessage(isGrpcError(err) ? grpcErrorMessage(err) : "Connection check failed")
            })
    }

    return (
        <div className={cls.EmailCheckButtonContainer}>
            <div className={cls.CheckRow}>
                {status === "ok" && <span className={cls.BadgeOk}>Connected</span>}
                {status === "fail" && <span className={cls.BadgeFail}>Failed</span>}
                <Button variant="secondary" onClick={handleCheck} disabled={disabled || status === "checking"}>
                    {status === "checking" ? "Checking…" : "Check settings"}
                </Button>
            </div>
            {status === "fail" && errorMessage && <p className={cls.ErrorText}>{errorMessage}</p>}
        </div>
    )
}
