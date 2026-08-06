import {useState} from "react"
import {Button} from "@vervstack/chures"

import Input from "@/components/atoms/Input/Input.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import TelegramCheckButton from "@/dialogs/ManageTelegramDialog/components/TelegramCheckButton/TelegramCheckButton.tsx"
import cls from "@/dialogs/ManageTelegramDialog/components/ConnectForm/ConnectForm.module.css"

export default function ConnectForm() {
    const [connecting, setConnecting] = useState(false)
    const [botToken, setBotToken] = useState("")
    const [verified, setVerified] = useState(false)

    const {addTelegramConnection} = useExternalConnections()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

    function handleTokenChange(value: string) {
        setBotToken(value)
        setVerified(false)
    }

    function handleConnect() {
        setConnecting(true)
        addTelegramConnection({botToken})
            .then(CloseDialog)
            .catch(e => bakeError("Failed to connect Telegram", e))
            .finally(() => setConnecting(false))
    }

    return (
        <div className={cls.ConnectFormContainer}>
            <p className={cls.ModalSub}>
                Connect your Telegram bot to send and receive messages through Telegram.
                We&apos;ll verify the token before saving it.
            </p>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Bot API Token</span>
                <Input type="text"
                       inputClassName={cls.TokenInput}
                       placeholder="123456789:ABCdefGHIjklmnoPQRstuvWXYz-1234567890"
                       value={botToken}
                    setValue={handleTokenChange} disabled={connecting} autoComplete="off"/>
                <div className={cls.TokenHint}>
                    <a href="https://t.me/botfather" target="_blank" rel="noopener noreferrer">
                        Get a token from @BotFather
                    </a>
                </div>
            </label>
            <div className={cls.ModalActions}>
                <TelegramCheckButton req={{botToken}}
                    disabled={connecting || !botToken} onResult={setVerified}/>
                <Button variant="primary" onClick={handleConnect}
                    disabled={connecting || !botToken || !verified}>
                    {connecting ? "Verifying…" : "Connect"}
                </Button>
            </div>
        </div>
    )
}
