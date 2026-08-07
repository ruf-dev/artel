import {useState} from "react"
import {Button} from "@vervstack/chures"

import {grpcErrorMessage, isGrpcError} from "@/processes/grpcErrors.ts"
import cls from "@/components/StorageCheckButton/StorageCheckButton.module.css"

export type StorageCheckStatus = "idle" | "checking" | "ok" | "fail"

interface StorageCheckButtonProps<TReq> {
    req: TReq
    disabled?: boolean
    onResult: (verified: boolean) => void
    onError: (message: string | null) => void
    checkConnection: (req: TReq) => Promise<unknown>
}

export default function StorageCheckButton<TReq>(
    {req, disabled, onResult, onError, checkConnection}: StorageCheckButtonProps<TReq>,
) {
    const [status, setStatus] = useState<StorageCheckStatus>("idle")

    function handleCheck() {
        setStatus("checking")
        onError(null)
        checkConnection(req)
            .then(() => {
                setStatus("ok")
                onResult(true)
            })
            .catch(err => {
                setStatus("fail")
                onError(isGrpcError(err) ? grpcErrorMessage(err) : "Connection check failed")
                onResult(false)
            })
    }

    return (
        <div className={cls.StorageCheckButtonContainer}>
            {status === "ok" && <span className={cls.BadgeOk}>Verified</span>}
            {status === "fail" && <span className={cls.BadgeFail}>Failed</span>}
            <Button variant="secondary" onClick={handleCheck} disabled={disabled || status === "checking"}>
                {status === "checking" ? "Testing…" : "Test"}
            </Button>
        </div>
    )
}
