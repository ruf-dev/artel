import {useState} from "react"
import {Button} from "@vervstack/chures"

import Input from "@/components/atoms/Input/Input.tsx"
import {AddEmailConnectionRequest} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import DialogHead from "@/dialogs/ManageEmailDialog/components/DialogHead/DialogHead.tsx"
import EmailCheckButton from "@/dialogs/ManageEmailDialog/components/EmailCheckButton/EmailCheckButton.tsx"
import HostPortRow from "@/dialogs/ManageEmailDialog/components/HostPortRow/HostPortRow.tsx"
import PasswordField from "@/dialogs/ManageEmailDialog/components/PasswordField/PasswordField.tsx"
import {useMailServerSuggestion} from "@/dialogs/ManageEmailDialog/processes/useMailServerSuggestion.ts"
import cls from "@/dialogs/ManageEmailDialog/components/EmailAddDialog/EmailAddDialog.module.css"

function buildEmailRequest(
    email: string, imapHost: string, imapPort: string, smtpHost: string, smtpPort: string, password: string,
): AddEmailConnectionRequest {
    return {
        email,
        imapHost,
        imapPort: imapPort ? parseInt(imapPort, 10) : undefined,
        smtpHost,
        smtpPort: smtpPort ? parseInt(smtpPort, 10) : undefined,
        password,
    }
}

export default function EmailAddDialog() {
    const [adding, setAdding] = useState(false)
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const suggestion = useMailServerSuggestion()

    const {addEmailConnection} = useExternalConnections()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

    function handleEmailChange(value: string) {
        setEmail(value)
        suggestion.notifyEmailChanged(value)
    }

    function handleAdd() {
        setAdding(true)
        addEmailConnection(buildEmailRequest(
            email, suggestion.imapHost, suggestion.imapPort, suggestion.smtpHost, suggestion.smtpPort, password,
        ))
            .then(CloseDialog)
            .catch(e => bakeError("Error adding email", e))
            .finally(() => setAdding(false))
    }

    return (
        <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <DialogHead title="Add email account" onClose={CloseDialog} disabled={adding} email={email}/>
            <p className={cls.ModalSub}>Connect an email account via IMAP and SMTP.</p>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Email address</span>
                <Input placeholder="you@example.com" value={email}
                    setValue={handleEmailChange} disabled={adding} autoComplete="off"/>
            </label>
            <HostPortRow
                host={{
                    label: "IMAP host", value: suggestion.imapHost,
                    placeholder: "imap.example.com", onChange: suggestion.setImapHost,
                }}
                port={{
                    label: "IMAP port", value: suggestion.imapPort,
                    placeholder: "993", onChange: suggestion.setImapPort,
                }}
                disabled={adding}
            />
            <HostPortRow
                host={{
                    label: "SMTP host", value: suggestion.smtpHost,
                    placeholder: "smtp.example.com", onChange: suggestion.setSmtpHost,
                }}
                port={{
                    label: "SMTP port", value: suggestion.smtpPort,
                    placeholder: "587", onChange: suggestion.setSmtpPort,
                }}
                disabled={adding}
            />
            <PasswordField value={password} onChange={setPassword} disabled={adding}
                appPasswordUrl={suggestion.appPasswordUrl} hasStoredValue={false}/>
            <div className={cls.ModalActions}>
                <EmailCheckButton
                    req={buildEmailRequest(
                        email, suggestion.imapHost, suggestion.imapPort,
                        suggestion.smtpHost, suggestion.smtpPort, password,
                    )}
                    disabled={adding}/>
                <Button variant="primary" onClick={handleAdd} disabled={adding}>
                    {adding ? "Adding…" : "Add account"}
                </Button>
            </div>
        </div>
    )
}
