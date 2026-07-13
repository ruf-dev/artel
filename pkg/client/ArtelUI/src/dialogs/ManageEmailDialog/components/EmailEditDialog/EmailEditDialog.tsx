import {useState, useEffect} from "react"
import {Button} from "@vervstack/chures"

import Input from "@/components/atoms/Input/Input.tsx"
import {
    AddEmailConnectionRequest,
    ExternalConnectionInfo,
} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import DialogHead from "@/dialogs/ManageEmailDialog/components/DialogHead/DialogHead.tsx"
import EmailCheckButton from "@/dialogs/ManageEmailDialog/components/EmailCheckButton/EmailCheckButton.tsx"
import HostPortRow from "@/dialogs/ManageEmailDialog/components/HostPortRow/HostPortRow.tsx"
import PasswordField from "@/dialogs/ManageEmailDialog/components/PasswordField/PasswordField.tsx"
import {useMailServerSuggestion} from "@/dialogs/ManageEmailDialog/processes/useMailServerSuggestion.ts"
import cls from "@/dialogs/ManageEmailDialog/components/EmailEditDialog/EmailEditDialog.module.css"

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

export default function EmailEditDialog({connection}: {connection: ExternalConnectionInfo}) {
    const [saving, setSaving] = useState(false)
    const initialEmail = connection.generic?.fields?.username ?? connection.generic?.fields?.email ?? ""
    const [email, setEmail] = useState(initialEmail)
    const [password, setPassword] = useState("")
    const suggestion = useMailServerSuggestion()

    const {addEmailConnection} = useExternalConnections()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

    useEffect(() => {
        suggestion.seedFromEmail(initialEmail)
    }, [])

    function handleEmailChange(value: string) {
        setEmail(value)
        suggestion.notifyEmailChanged(value)
    }

    function handleSave() {
        setSaving(true)
        addEmailConnection(buildEmailRequest(
            email, suggestion.imapHost, suggestion.imapPort, suggestion.smtpHost, suggestion.smtpPort, password,
        ))
            .then(CloseDialog)
            .catch(e => bakeError("Error updating email", e))
            .finally(() => setSaving(false))
    }

    return (
        <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <DialogHead
                title="Edit email account"
                onClose={CloseDialog}
                disabled={saving}
                email={email}/>
            <p className={cls.ModalSub}>Update server settings. Leave the password blank to keep the current one.</p>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Email address</span>
                <Input value={email} setValue={handleEmailChange} disabled={saving} autoComplete="off"/>
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
                disabled={saving}
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
                disabled={saving}
            />
            <PasswordField value={password} onChange={setPassword} disabled={saving}
                appPasswordUrl={suggestion.appPasswordUrl} hasStoredValue={true}/>
            <div className={cls.ModalActions}>
                <EmailCheckButton
                    req={buildEmailRequest(
                        email, suggestion.imapHost, suggestion.imapPort,
                        suggestion.smtpHost, suggestion.smtpPort, password,
                    )}
                    disabled={saving}/>
                <Button variant="primary" onClick={handleSave} disabled={saving}>
                    {saving ? "Saving…" : "Save changes"}
                </Button>
            </div>
        </div>
    )
}
