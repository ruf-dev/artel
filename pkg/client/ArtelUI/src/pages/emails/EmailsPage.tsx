import {useState, useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/emails/EmailsPage.module.css"

import {EmailAccountInfo, AddEmailAccountRequest} from "@/app/api/artel/email_accounts.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useEmailAccounts} from "@/app/hooks/EmailAccounts.ts"
import useUser from "@/hooks/user/User.ts"

import Topbar from "@/components/Topbar/Topbar.tsx"
import ModalClose from "@/components/ModalClose/ModalClose.tsx"
import FormField from "@/components/FormField/FormField.tsx"
import ModalActions from "@/components/ModalActions/ModalActions.tsx"

export default function EmailsPage() {
    const navigate = useNavigate()
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const {fetch: fetchAccounts} = useEmailAccounts()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetchAccounts()
        }
    }, [auth, fetchAccounts])

    return (
        <div className={cls.Root}>
            <Topbar/>
            <HeroSegment onAddClick={() => OpenDialog(<AddAccountDialog/>)}/>
            <ContentSegment/>
        </div>
    )
}

function HeroSegment({onAddClick}: { onAddClick: () => void }) {
    const {accounts, loading} = useEmailAccounts()

    return (
        <div className={cls.Hero}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>Workspace</div>
                <h1 className={cls.HeroTitle}>Your email accounts</h1>
                <p className={cls.HeroSub}>
                    <b>{loading ? "…" : `${accounts.length} ${accounts.length === 1 ? "account" : "accounts"}`}</b>
                    {" · "}<span>connected via IMAP / SMTP</span>
                </p>
            </div>
            <button className={cls.AddBtn} onClick={onAddClick}>
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2"
                     strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Add account
            </button>
        </div>
    )
}

function ContentSegment() {
    const {accounts, loading, remove} = useEmailAccounts()

    if (loading) {
        return (
            <div className={cls.Content}>
                <p className={cls.Empty}>Loading…</p>
            </div>
        )
    }

    return (
        <div className={cls.Content}>
            <div className={cls.List}>
                {accounts.map(a => (
                    <EmailAccountRow key={a.id} account={a} onRemove={remove}/>
                ))}
                {accounts.length === 0 && (
                    <p className={cls.Empty}>No email accounts yet. Add one to get started.</p>
                )}
            </div>
        </div>
    )
}

function EmailAccountRow({account, onRemove}: {
    account: EmailAccountInfo
    onRemove: (id: string) => Promise<void>
}) {
    const [removing, setRemoving] = useState(false)

    async function handleRemove() {
        if (!account.id) return
        setRemoving(true)
        try {
            await onRemove(account.id)
        } finally {
            setRemoving(false)
        }
    }

    return (
        <div className={cls.Row}>
            <div className={cls.RowInfo}>
                <span className={cls.RowEmail}>{account.email}</span>
                <span className={cls.RowMeta}>
                    IMAP {account.imapHost}:{account.imapPort} · SMTP {account.smtpHost}:{account.smtpPort}
                </span>
            </div>
            <button className={cls.BtnDanger} onClick={handleRemove} disabled={removing} type="button">
                {removing ? "Removing…" : "Remove"}
            </button>
        </div>
    )
}

function AddAccountDialog() {
    const [adding, setAdding] = useState(false)
    const [email, setEmail] = useState("")
    const [imapHost, setImapHost] = useState("")
    const [imapPort, setImapPort] = useState("")
    const [smtpHost, setSmtpHost] = useState("")
    const [smtpPort, setSmtpPort] = useState("")
    const [password, setPassword] = useState("")

    const {add} = useEmailAccounts()
    const {CloseDialog} = useDialog()

    async function handleAdd() {
        setAdding(true)
        const req: AddEmailAccountRequest = {
            email,
            imapHost,
            imapPort: imapPort ? parseInt(imapPort, 10) : undefined,
            smtpHost,
            smtpPort: smtpPort ? parseInt(smtpPort, 10) : undefined,
            password,
        }
        add(req)
            .then(CloseDialog)
            .catch(console.error)
            .finally(() => setAdding(false))
    }

    return (
        <div className={cls.Overlay}>
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                 aria-labelledby="addEmailModalTitle">
                <div className={cls.ModalHead}>
                    <h2 className={cls.ModalTitle} id="addEmailModalTitle">Add email account</h2>
                    <ModalClose onClick={CloseDialog} disabled={adding} className={cls.ModalClose}/>
                </div>
                <p className={cls.ModalSub}>Connect an email account via IMAP and SMTP.</p>

                <FormField
                    label="Email address"
                    placeholder="you@example.com"
                    onChange={setEmail}
                    disabled={adding}
                    fieldClassName={cls.Field}
                    labelClassName={cls.FieldLabel}
                    inputClassName={cls.Input}
                />
                <div className={cls.FieldRow}>
                    <FormField
                        label="IMAP host"
                        placeholder="imap.example.com"
                        onChange={setImapHost}
                        disabled={adding}
                        fieldClassName={cls.Field}
                        labelClassName={cls.FieldLabel}
                        inputClassName={cls.Input}
                    />
                    <FormField
                        label="IMAP port"
                        placeholder="993"
                        onChange={setImapPort}
                        disabled={adding}
                        maxLength={5}
                        fieldClassName={cls.FieldNarrow}
                        labelClassName={cls.FieldLabel}
                        inputClassName={cls.Input}
                    />
                </div>
                <div className={cls.FieldRow}>
                    <FormField
                        label="SMTP host"
                        placeholder="smtp.example.com"
                        onChange={setSmtpHost}
                        disabled={adding}
                        fieldClassName={cls.Field}
                        labelClassName={cls.FieldLabel}
                        inputClassName={cls.Input}
                    />
                    <FormField
                        label="SMTP port"
                        placeholder="587"
                        onChange={setSmtpPort}
                        disabled={adding}
                        maxLength={5}
                        fieldClassName={cls.FieldNarrow}
                        labelClassName={cls.FieldLabel}
                        inputClassName={cls.Input}
                    />
                </div>
                <div className={cls.Field}>
                    <span className={cls.FieldLabel}>Password</span>
                    <input
                        type="password"
                        placeholder="App-specific password"
                        className={cls.Input}
                        onChange={e => setPassword(e.target.value)}
                        disabled={adding}
                        autoComplete="new-password"
                    />
                </div>

                <ModalActions
                    containerClassName={cls.ModalActions}
                    buttons={[
                        {
                            label: adding ? "Adding…" : "Add account",
                            onClick: handleAdd,
                            className: cls.BtnPrimary,
                            disabled: adding,
                        }
                    ]}
                />
            </div>
        </div>
    )
}
