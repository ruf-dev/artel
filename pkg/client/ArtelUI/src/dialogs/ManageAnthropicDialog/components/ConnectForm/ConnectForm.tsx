import {useState} from "react"
import {Button} from "@vervstack/chures"

import Input from "@/components/atoms/Input/Input.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import AnthropicCheckButton
    from "@/dialogs/ManageAnthropicDialog/components/AnthropicCheckButton/AnthropicCheckButton.tsx"
import cls from "@/dialogs/ManageAnthropicDialog/components/ConnectForm/ConnectForm.module.css"

export default function ConnectForm() {
    const [connecting, setConnecting] = useState(false)
    const [apiKey, setApiKey] = useState("")
    const [baseUrl, setBaseUrl] = useState("")
    const [model, setModel] = useState("")
    const [verified, setVerified] = useState(false)
    const [checkError, setCheckError] = useState<string | null>(null)

    const {addAnthropicConnection} = useExternalConnections()
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
        addAnthropicConnection({apiKey, baseUrl, defaultModel: model})
            .then(CloseDialog)
            .catch(e => bakeError("Failed to connect Claude (Anthropic)", e))
            .finally(() => setConnecting(false))
    }

    return (
        <div className={cls.ConnectFormContainer}>
            <p className={cls.ModalSub}>
                Connect your Anthropic API key to let Claude access it as a BYOK LLM provider.
                We&apos;ll verify the key against Anthropic before saving it.
            </p>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>API key</span>
                <Input type="text" inputClassName={cls.KeyInput} placeholder="sk-ant-…" value={apiKey}
                    setValue={handleApiKeyChange} disabled={connecting} autoComplete="off"/>
            </label>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Base URL (optional)</span>
                <Input placeholder="https://api.anthropic.com" value={baseUrl}
                    setValue={handleBaseUrlChange} disabled={connecting} autoComplete="off"/>
            </label>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Model id (optional)</span>
                <Input placeholder="claude-opus-4-8" value={model}
                    setValue={handleModelChange} disabled={connecting} autoComplete="off"/>
                <div className={cls.FieldHint}>
                    Used as the default model. If your provider doesn&apos;t support listing models
                    (common for non-official Anthropic-compatible endpoints), it&apos;s also used to
                    verify the key.
                </div>
            </label>
            <div className={cls.ModalActions}>
                <AnthropicCheckButton req={{apiKey, baseUrl, defaultModel: model}}
                    disabled={connecting || !apiKey} onResult={handleCheckResult} onError={setCheckError}/>
                <Button variant="primary" onClick={handleConnect}
                    disabled={connecting || !apiKey || !verified}>
                    {connecting ? "Connecting…" : "Connect"}
                </Button>
            </div>
            {checkError && <p className={cls.CheckErrorText}>{checkError}</p>}
        </div>
    )
}
