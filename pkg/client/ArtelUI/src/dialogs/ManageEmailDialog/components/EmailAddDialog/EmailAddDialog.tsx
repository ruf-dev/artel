import {useState, useRef} from "react"
import {Button} from "@vervstack/chures"

import Input from "@/components/atoms/Input/Input.tsx"
import {
    AddEmailConnectionRequest,
    ExternalConnectionsAPI,
} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import useUser from "@/hooks/user/User.ts"
import DialogHead from "@/dialogs/ManageEmailDialog/components/DialogHead/DialogHead.tsx"
import HostPortRow from "@/dialogs/ManageEmailDialog/components/HostPortRow/HostPortRow.tsx"
import cls from "@/dialogs/ManageEmailDialog/components/EmailAddDialog/EmailAddDialog.module.css"

export default function EmailAddDialog() {
    const [adding, setAdding] = useState(false)
    const [email, setEmail] = useState("")
    const [imapHost, setImapHost] = useState("")
    const [imapPort, setImapPort] = useState("")
    const [smtpHost, setSmtpHost] = useState("")
    const [smtpPort, setSmtpPort] = useState("")
    const [password, setPassword] = useState("")
    const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)

    const {addEmailConnection} = useExternalConnections()
    const {CloseDialog} = useDialog()
    const {auth} = useUser()
    const bakeError = useBakeError()

    function handleEmailChange(value: string) {
        setEmail(value)
        const domain = value.split("@")[1] ?? ""
        if (!domain) return

        if (debounceRef.current) clearTimeout(debounceRef.current)
        debounceRef.current = setTimeout(() => {
            ExternalConnectionsAPI.ListMailServerSuggestions({domain}, auth.getInitReq())
                .then(resp => {
                    if (resp.suggestions?.length === 1) {
                        const s = resp.suggestions[0]
                        setImapHost(s.imap ?? "")
                        setImapPort(s.imapPort ? String(s.imapPort) : "")
                        setSmtpHost(s.smtp ?? "")
                        setSmtpPort(s.smtpPort ? String(s.smtpPort) : "")
                    }
                })
                .catch(() => {})
        }, 300)
    }

    function handleAdd() {
        setAdding(true)
        const req: AddEmailConnectionRequest = {
            email,
            imapHost,
            imapPort: imapPort ? parseInt(imapPort, 10) : undefined,
            smtpHost,
            smtpPort: smtpPort ? parseInt(smtpPort, 10) : undefined,
            password,
        }
        addEmailConnection(req)
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
                host={{label: "IMAP host", value: imapHost, placeholder: "imap.example.com", onChange: setImapHost}}
                port={{label: "IMAP port", value: imapPort, placeholder: "993", onChange: setImapPort}}
                disabled={adding}
            />
            <HostPortRow
                host={{label: "SMTP host", value: smtpHost, placeholder: "smtp.example.com", onChange: setSmtpHost}}
                port={{label: "SMTP port", value: smtpPort, placeholder: "587", onChange: setSmtpPort}}
                disabled={adding}
            />
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Password</span>
                <Input type="password" placeholder="App-specific password" value={password}
                    setValue={setPassword} disabled={adding} autoComplete="new-password"/>
            </label>
            <div className={cls.ModalActions}>
                <Button variant="primary" onClick={handleAdd} disabled={adding}>
                    {adding ? "Adding…" : "Add account"}
                </Button>
            </div>
        </div>
    )
}
