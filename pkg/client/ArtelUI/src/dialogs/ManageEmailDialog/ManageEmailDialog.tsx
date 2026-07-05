import {useState, useRef, useEffect} from "react"

import {Button, ModalClose, ConfirmDialog} from "@vervstack/chures"
import cls from "@/dialogs/ManageEmailDialog/ManageEmailDialog.module.css"

import {
    AddEmailConnectionRequest,
    ExternalConnectionInfo,
    ExternalConnectionsAPI,
    ExternalProvider
} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import useUser from "@/hooks/user/User.ts"

import {mailProviderIcon} from "@/app/utils/mailProviderIcon"

import Input from "@/components/shared/Input/Input.tsx"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"

export default function ManageEmailDialog() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {connections, loading, disconnect} = useExternalConnections()
    const bakeError = useBakeError()

    const emailConnections = connections.filter(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_EMAIL)

    function handleEdit(conn: ExternalConnectionInfo) {
        OpenDialog(<EmailEditDialog connection={conn}/>)
    }

    function handleRemove(conn: ExternalConnectionInfo) {
        const email = conn.generic?.fields?.username ?? conn.generic?.fields?.email ?? ""
        OpenDialog(
            <ConfirmDialog
                title="Remove email account"
                message={`Remove "${email}"? You can re-add it at any time.`}
                confirmLabel="Remove"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() =>
                    disconnect("email")
                        .catch(e => bakeError("Failed to remove account", e))}
            />
        )
    }

    return (
        <div className={cls.ModalContainer}
             onClick={e => e.stopPropagation()}
             role="dialog"
             aria-modal="true">
            <DialogHead title="Email" onClose={CloseDialog}/>
            <AccountsSection
                loading={loading}
                connections={emailConnections}
                onAdd={() => OpenDialog(<EmailAddDialog/>)}
                onRemove={handleRemove}
                onEdit={handleEdit}
            />
        </div>
    )
}

function DialogHead({title, onClose, disabled, email}: { title: string; onClose: () => void; disabled?: boolean; email?: string }) {
    const icon = email ? mailProviderIcon(email) : null

    return (
        <div className={cls.ModalHead}>
            <div className={cls.ModalHeadLeft}>
                <div className={cls.ModalIcon}>
                    {icon
                        ? <img src={icon} className={cls.ModalProviderIcon} alt=""/>
                        : <ProviderIcon provider={ExternalProvider.EXTERNAL_PROVIDER_EMAIL}/>
                    }
                </div>
                <span className={cls.ModalTitle}>{title}</span>
            </div>
            <ModalClose onClick={onClose} disabled={disabled} className={cls.ModalClose}/>
        </div>
    )
}

function AccountsSection({loading, connections, onAdd, onRemove, onEdit}: {
    loading: boolean
    connections: ExternalConnectionInfo[]
    onAdd: () => void
    onRemove: (conn: ExternalConnectionInfo) => void
    onEdit: (conn: ExternalConnectionInfo) => void
}) {
    return (
        <div className={cls.AccountSection}>
            <div className={cls.AccountSectionHeader}>
                <span className={cls.AccountSectionLabel}>Accounts</span>
                <Button variant="ghost" onClick={onAdd}>+ Add</Button>
            </div>
            {loading ? (
                <p className={cls.Empty}>Loading…</p>
            ) : connections.length === 0 ? (
                <p className={cls.Empty}>No email accounts added yet.</p>
            ) : (
                <div className={cls.AccountList}>
                    {connections.map(conn => (
                        <EmailConnectionRow key={conn.id} connection={conn} onRemove={onRemove} onEdit={onEdit}/>
                    ))}
                </div>
            )}
        </div>
    )
}

function EmailConnectionRow({connection, onRemove, onEdit}: {
    connection: ExternalConnectionInfo
    onRemove: (conn: ExternalConnectionInfo) => void
    onEdit: (conn: ExternalConnectionInfo) => void
}) {
    const email = connection.generic?.fields?.username ?? connection.generic?.fields?.email ?? "—"
    const icon = mailProviderIcon(email)
    return (
        <div className={cls.AccountRow} onClick={() => onEdit(connection)} role="button" tabIndex={0}>
            <div className={cls.AccountRowLeft}>
                {icon
                    ? <img src={icon} className={cls.AccountEmailIcon} alt=""/>
                    : <span className={cls.AccountEmailAtBadge}>@</span>
                }
                <span className={cls.AccountEmail}>{email}</span>
            </div>
            <div className={cls.AccountRowActions}>
                <Button variant="iconDanger" onClick={e => { e.stopPropagation(); onRemove(connection) }} aria-label="Remove">×</Button>
            </div>
        </div>
    )
}

function EmailEditDialog({connection}: {connection: ExternalConnectionInfo}) {
    const [saving, setSaving] = useState(false)
    const initialEmail = connection.generic?.fields?.username ?? connection.generic?.fields?.email ?? ""
    const [email, setEmail] = useState(initialEmail)
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

    useEffect(() => {
        const domain = initialEmail.split("@")[1] ?? ""
        if (!domain) return
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
    }, [])

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

    function handleSave() {
        setSaving(true)
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
            <p className={cls.ModalSub}>Update server settings. Re-enter your password to apply changes.</p>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Email address</span>
                <Input value={email} onChange={e => handleEmailChange(e.target.value)} disabled={saving} autoComplete="off"/>
            </label>
            <HostPortRow
                host={{label: "IMAP host", value: imapHost, placeholder: "imap.example.com", onChange: setImapHost}}
                port={{label: "IMAP port", value: imapPort, placeholder: "993", onChange: setImapPort}}
                disabled={saving}
            />
            <HostPortRow
                host={{label: "SMTP host", value: smtpHost, placeholder: "smtp.example.com", onChange: setSmtpHost}}
                port={{label: "SMTP port", value: smtpPort, placeholder: "587", onChange: setSmtpPort}}
                disabled={saving}
            />
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Password</span>
                <Input type="password" placeholder="App-specific password"
                    onChange={e => setPassword(e.target.value)} disabled={saving} autoComplete="new-password"/>
            </label>
            <div className={cls.ModalActions}>
                <Button variant="primary" onClick={handleSave} disabled={saving}>
                    {saving ? "Saving…" : "Save changes"}
                </Button>
            </div>
        </div>
    )
}

function HostPortRow({host, port, disabled}: {
    host: { label: string; value: string; placeholder: string; onChange: (v: string) => void }
    port: { label: string; value: string; placeholder: string; onChange: (v: string) => void }
    disabled?: boolean
}) {
    return (
        <div className={cls.FieldRow}>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>{host.label}</span>
                <Input placeholder={host.placeholder} value={host.value}
                    onChange={e => host.onChange(e.target.value)} disabled={disabled} autoComplete="off"/>
            </label>
            <label className={cls.FieldNarrow}>
                <span className={cls.FieldLabel}>{port.label}</span>
                <Input placeholder={port.placeholder} value={port.value}
                    onChange={e => port.onChange(e.target.value)} disabled={disabled} maxLength={5}
                    autoComplete="off"/>
            </label>
        </div>
    )
}

function EmailAddDialog() {
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
                    onChange={e => handleEmailChange(e.target.value)} disabled={adding} autoComplete="off"/>
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
                <Input type="password" placeholder="App-specific password"
                    onChange={e => setPassword(e.target.value)} disabled={adding} autoComplete="new-password"/>
            </label>
            <div className={cls.ModalActions}>
                <Button variant="primary" onClick={handleAdd} disabled={adding}>
                    {adding ? "Adding…" : "Add account"}
                </Button>
            </div>
        </div>
    )
}
