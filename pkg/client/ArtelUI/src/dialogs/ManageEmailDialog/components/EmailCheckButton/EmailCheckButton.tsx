import {useState} from "react"
import {Button} from "@vervstack/chures"

import {CheckEmailConnectionRequest, ExternalConnectionsAPI} from "@/app/api/artel/external_connections.pb.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import useUser from "@/hooks/user/User.ts"
import cls from "@/dialogs/ManageEmailDialog/components/EmailCheckButton/EmailCheckButton.module.css"

export type CheckStatus = "idle" | "checking" | "ok" | "fail"

interface EmailCheckButtonProps {
    req: CheckEmailConnectionRequest
    disabled?: boolean
}

export default function EmailCheckButton({req, disabled}: EmailCheckButtonProps) {
    const [status, setStatus] = useState<CheckStatus>("idle")
    const {auth} = useUser()
    const bakeError = useBakeError()

    function handleCheck() {
        setStatus("checking")
        ExternalConnectionsAPI.CheckEmailConnection(req, auth.getInitReq())
            .then(() => setStatus("ok"))
            .catch(err => {
                setStatus("fail")
                bakeError("Connection check failed", err)
            })
    }

    return (
        <div className={cls.EmailCheckButtonContainer}>
            {status === "ok" && <span className={cls.BadgeOk}>Connected</span>}
            {status === "fail" && <span className={cls.BadgeFail}>Failed</span>}
            <Button variant="secondary" onClick={handleCheck} disabled={disabled || status === "checking"}>
                {status === "checking" ? "Checking…" : "Check settings"}
            </Button>
        </div>
    )
}
