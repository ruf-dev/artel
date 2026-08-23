import {useState} from "react"
import {Button} from "@vervstack/chures"

import Input from "@/components/atoms/Input/Input.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import type {CheckAnthropicConnectionResponse} from "@/app/api/artel/external_connections.pb.ts"
import LlmKeyCheckButton from "@/components/LlmKeyCheckButton/LlmKeyCheckButton.tsx"
import cls from "@/components/LlmKeyConnectForm/LlmKeyConnectForm.module.css"

type CheckRequest = {apiKey?: string; baseUrl?: string; defaultModel?: string}

interface LlmKeyConnectFormProps {
    providerName: string
    bodyCopy: string
    apiKeyPlaceholder: string
    baseUrlPlaceholder?: string
    modelPlaceholder?: string
    fixedBaseUrl?: string
    hideModelField?: boolean
    addConnection: (req: CheckRequest) => Promise<unknown>
    checkConnection: (req: CheckRequest) => Promise<CheckAnthropicConnectionResponse>
}

export default function LlmKeyConnectForm(props: LlmKeyConnectFormProps) {
    const [connecting, setConnecting] = useState(false)
    const [apiKey, setApiKey] = useState("")
    const [baseUrl, setBaseUrl] = useState(props.fixedBaseUrl ?? "")
    const [model, setModel] = useState("")
    const [verified, setVerified] = useState(false)
    const [checkError, setCheckError] = useState<string | null>(null)

    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

    function handleApiKeyChange(value: string) {
        setApiKey(value)
        setVerified(false)
    }

    function handleBaseUrlChange(value: string) {
        setBaseUrl(value)
        setVerified(false)
    }

    function handleModelChange(value: string) {
        setModel(value)
        setVerified(false)
    }

    function handleCheckResult(ok: boolean, recommended?: string) {
        setVerified(ok)
        if (ok && !model && recommended) {
            setModel(recommended)
        }
    }

    function handleConnect() {
        setConnecting(true)
        props.addConnection({apiKey, baseUrl, defaultModel: model})
            .then(CloseDialog)
            .catch(e => bakeError(`Failed to connect ${props.providerName}`, e))
            .finally(() => setConnecting(false))
    }

    return (
        <div className={cls.LlmKeyConnectFormContainer}>
            <p className={cls.ModalSub}>
                {props.bodyCopy}
            </p>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>API key</span>
                <Input
                    type="text"
                    inputClassName={cls.KeyInput}
                    placeholder={props.apiKeyPlaceholder}
                    value={apiKey}
                    setValue={handleApiKeyChange}
                    disabled={connecting}
                    autoComplete="off"
                />
            </label>
            {!props.fixedBaseUrl && (
                <label className={cls.Field}>
                    <span className={cls.FieldLabel}>Base URL (optional)</span>
                    <Input
                        placeholder={props.baseUrlPlaceholder}
                        value={baseUrl}
                        setValue={handleBaseUrlChange}
                        disabled={connecting}
                        autoComplete="off"
                    />
                </label>
            )}
            {!props.hideModelField && (
                <label className={cls.Field}>
                    <span className={cls.FieldLabel}>Model id (optional)</span>
                    <Input
                        placeholder={props.modelPlaceholder}
                        value={model}
                        setValue={handleModelChange}
                        disabled={connecting}
                        autoComplete="off"
                    />
                    <div className={cls.FieldHint}>
                        Used as the default model. If your provider doesn&apos;t support listing models
                        (common for non-official endpoints), it&apos;s also used to
                        verify the key.
                    </div>
                </label>
            )}
            <div className={cls.ModalActions}>
                <LlmKeyCheckButton
                    req={{apiKey, baseUrl, defaultModel: model}}
                    disabled={connecting || !apiKey}
                    onResult={handleCheckResult}
                    onError={setCheckError}
                    checkConnection={props.checkConnection}
                />
                <Button
                    variant="primary"
                    onClick={handleConnect}
                    disabled={connecting || !apiKey || !verified}
                >
                    {connecting ? "Connecting…" : "Connect"}
                </Button>
            </div>
            {checkError && <p className={cls.CheckErrorText}>{checkError}</p>}
        </div>
    )
}
