import {useState} from "react"
import {Button} from "@vervstack/chures"

import type {
    CheckAnthropicConnectionRequest,
    CheckAnthropicConnectionResponse,
} from "@/app/api/artel/external_connections.pb.ts"
import {grpcErrorMessage, isGrpcError} from "@/processes/grpcErrors.ts"
import cls from "@/components/LlmKeyCheckButton/LlmKeyCheckButton.module.css"

export type CheckStatus = "idle" | "checking" | "ok" | "fail"

interface LlmKeyCheckButtonProps {
    req: CheckAnthropicConnectionRequest
    disabled?: boolean
    onResult: (verified: boolean, recommendedDefaultModel?: string) => void
    onError: (message: string | null) => void
    checkConnection: (req: CheckAnthropicConnectionRequest) => Promise<CheckAnthropicConnectionResponse>
}

export default function LlmKeyCheckButton({req, disabled, onResult, onError, checkConnection}: LlmKeyCheckButtonProps) {
    const [status, setStatus] = useState<CheckStatus>("idle")
    const [modelCount, setModelCount] = useState(0)

    function handleCheck() {
        setStatus("checking")
        onError(null)
        checkConnection(req)
            .then(resp => {
                setStatus("ok")
                setModelCount(resp.availableModels?.length ?? 0)
                onResult(true, resp.recommendedDefaultModel)
            })
            .catch(err => {
                setStatus("fail")
                onError(isGrpcError(err) ? grpcErrorMessage(err) : "Connection check failed")
                onResult(false)
            })
    }

    return (
        <div className={cls.LlmKeyCheckButtonContainer}>
            {status === "ok" && <span className={cls.BadgeOk}>{modelCount} models available</span>}
            {status === "fail" && <span className={cls.BadgeFail}>Failed</span>}
            <Button variant="secondary" onClick={handleCheck} disabled={disabled || status === "checking"}>
                {status === "checking" ? "Testing…" : "Test"}
            </Button>
        </div>
    )
}
