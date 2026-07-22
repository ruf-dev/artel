import {useState} from "react"
import {Button} from "@vervstack/chures"

import {CheckAnthropicConnectionRequest} from "@/app/api/artel/external_connections.pb.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {grpcErrorMessage, isGrpcError} from "@/processes/grpcErrors.ts"
import cls from "@/dialogs/ManageAnthropicDialog/components/AnthropicCheckButton/AnthropicCheckButton.module.css"

export type CheckStatus = "idle" | "checking" | "ok" | "fail"

interface AnthropicCheckButtonProps {
    req: CheckAnthropicConnectionRequest
    disabled?: boolean
    onResult: (verified: boolean, recommendedDefaultModel?: string) => void
    onError: (message: string | null) => void
}

export default function AnthropicCheckButton({req, disabled, onResult, onError}: AnthropicCheckButtonProps) {
    const [status, setStatus] = useState<CheckStatus>("idle")
    const [modelCount, setModelCount] = useState(0)
    const {checkAnthropicConnection} = useExternalConnections()

    function handleCheck() {
        setStatus("checking")
        onError(null)
        checkAnthropicConnection(req)
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
        <div className={cls.AnthropicCheckButtonContainer}>
            {status === "ok" && <span className={cls.BadgeOk}>{modelCount} models available</span>}
            {status === "fail" && <span className={cls.BadgeFail}>Failed</span>}
            <Button variant="secondary" onClick={handleCheck} disabled={disabled || status === "checking"}>
                {status === "checking" ? "Testing…" : "Test"}
            </Button>
        </div>
    )
}
