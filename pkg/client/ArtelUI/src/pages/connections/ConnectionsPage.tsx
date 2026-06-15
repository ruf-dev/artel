import {useEffect, useRef, useState} from "react"
import {useNavigate, useSearchParams} from "react-router-dom"

import cls from "@/pages/connections/ConnectionsPage.module.css"

import {ExternalConnectionInfo, ExternalProvider, Spreadsheet} from "@/app/api/artel/external_connections.pb.ts"
import {EmailAccountInfo, AddEmailAccountRequest, EmailAccountsAPI} from "@/app/api/artel/email_accounts.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useEmailAccounts} from "@/app/hooks/EmailAccounts.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import useUser from "@/hooks/user/User.ts"

import ConfirmDialog from "@/components/ConfirmDialog/ConfirmDialog.tsx"
import ModalClose from "@/components/ModalClose/ModalClose.tsx"

type GapiWindow = Window & {gapi: {load: (lib: string, cb: () => void) => void}}

let gapiPromise: Promise<void> | null = null

function loadGapi(): Promise<void> {
    if (gapiPromise) return gapiPromise
    gapiPromise = new Promise<void>((resolve, reject) => {
        const script = document.createElement("script")
        script.src = "https://apis.google.com/js/api.js"
        script.onload = () => (window as unknown as GapiWindow).gapi.load("picker", resolve)
        script.onerror = reject
        document.head.appendChild(script)
    })
    return gapiPromise
}

type ProviderDef = {
    provider: ExternalProvider
    name: string
    icon: React.JSX.Element
    description: string
    canConnect: boolean
}

const PROVIDERS: ProviderDef[] = [
    {
        provider: ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS,
        name: "Google Sheets",
        icon: <GoogleSheetsIcon/>,
        description: "Read and write data from your Google Sheets spreadsheets.",
        canConnect: true,
    },
    {
        provider: ExternalProvider.EXTERNAL_PROVIDER_TRELLO,
        name: "Trello",
        icon: <TrelloIcon/>,
        description: "Sync tasks and boards from your Trello workspace.",
        canConnect: false,
    },
    {
        provider: ExternalProvider.EXTERNAL_PROVIDER_MIRO,
        name: "Miro",
        icon: <MiroIcon/>,
        description: "Access and embed your Miro boards.",
        canConnect: false,
    },
]

function providerKey(p: ExternalProvider): string {
    const map: Record<ExternalProvider, string> = {
        [ExternalProvider.EXTERNAL_PROVIDER_UNSPECIFIED]: "",
        [ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS]: "google_sheets",
        [ExternalProvider.EXTERNAL_PROVIDER_TRELLO]: "trello",
        [ExternalProvider.EXTERNAL_PROVIDER_MIRO]: "miro",
    }
    return map[p] ?? ""
}

export default function ConnectionsPage() {
    const navigate = useNavigate()
    const {auth} = useUser()
    const {fetch: fetchConnections} = useExternalConnections()
    const {fetch: fetchAccounts} = useEmailAccounts()
    const bakeError = useBakeError()
    const [searchParams, setSearchParams] = useSearchParams()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    useEffect(() => {
        if (!auth.isAuthenticated()) return

        const status = searchParams.get("status")
        if (status === "error") {
            bakeError("Connection failed", new Error("Google OAuth authorization was denied or failed."))
        }
        if (status) {
            setSearchParams({}, {replace: true})
        }

        void fetchConnections()
        void fetchAccounts()
    }, [auth])

    return (
        <div className={cls.Root}>
            <HeroSegment/>
            <ContentSegment/>
        </div>
    )
}

function HeroSegment() {
    const {connections, loading} = useExternalConnections()
    const connectedCount = connections.length

    return (
        <div className={cls.Hero}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>Workspace</div>
                <h1 className={cls.HeroTitle}>Connected services</h1>
                <p className={cls.HeroSub}>
                    <b>{loading ? "…" : `${connectedCount} ${connectedCount === 1 ? "service" : "services"}`}</b>
                    {" · "}<span>linked external integrations</span>
                </p>
            </div>
        </div>
    )
}

function ContentSegment() {
    const {connections, loading} = useExternalConnections()
    const {accounts, loading: emailsLoading} = useEmailAccounts()
    const {OpenDialog} = useDialog()

    if (loading) {
        return (
            <div className={cls.Content}>
                <div className={cls.Grid}>
                    {PROVIDERS.map(def => (
                        <ProviderCard key={def.provider} def={def} connection={undefined} loading/>
                    ))}
                    <EmailCard accounts={[]} loading/>
                </div>
            </div>
        )
    }

    function findConnection(p: ExternalProvider): ExternalConnectionInfo | undefined {
        return connections.find(c => c.provider === p)
    }

    return (
        <div className={cls.Content}>
            <div className={cls.Grid}>
                {PROVIDERS.map(def => {
                    const connection = findConnection(def.provider)
                    return (
                        <ProviderCard
                            key={def.provider}
                            def={def}
                            connection={connection}
                            loading={false}
                            onClick={() => OpenDialog(<ConnectionDetailDialog def={def} connection={connection}/>)}
                        />
                    )
                })}
                <EmailCard
                    accounts={accounts}
                    loading={emailsLoading}
                    onClick={() => OpenDialog(<EmailDetailDialog/>)}
                />
            </div>
        </div>
    )
}

function ProviderCard({def, connection, loading, onClick}: {
    def: ProviderDef
    connection: ExternalConnectionInfo | undefined
    loading: boolean
    onClick?: () => void
}) {
    const isConnected = !!connection
    const accountLabel = connection?.google?.email ?? undefined

    return (
        <div className={cls.Card} onClick={!loading ? onClick : undefined} role="button" tabIndex={0}
             onKeyDown={e => { if (e.key === "Enter" || e.key === " ") onClick?.() }}>
            <div className={cls.CardHeader}>
                <div className={cls.CardIcon}>{def.icon}</div>
                <div className={cls.CardTitles}>
                    <div className={cls.CardName}>{def.name}</div>
                    {accountLabel && <div className={cls.CardAccount}>{accountLabel}</div>}
                </div>
            </div>
            <div className={cls.CardFooter}>
                <span className={`${cls.StatusDot} ${isConnected ? cls.StatusDotConnected : cls.StatusDotDisconnected}`}/>
                <span className={`${cls.StatusLabel} ${isConnected ? cls.StatusLabelConnected : cls.StatusLabelDisconnected}`}>
                    {loading ? "…" : isConnected ? "Connected" : "Not connected"}
                </span>
            </div>
        </div>
    )
}

function ConnectionDetailDialog({def, connection}: {
    def: ProviderDef
    connection: ExternalConnectionInfo | undefined
}) {
    const {CloseDialog, OpenDialog} = useDialog()
    const {disconnect, initiateGoogleOAuth} = useExternalConnections()
    const bakeError = useBakeError()
    const [connecting, setConnecting] = useState(false)

    async function handleConnectGoogle() {
        setConnecting(true)
        try {
            const authUrl = await initiateGoogleOAuth()
            window.location.href = authUrl
        } catch (e) {
            bakeError("Failed to start OAuth", e)
            setConnecting(false)
        }
    }

    function handleDisconnect() {
        const key = providerKey(def.provider)
        OpenDialog(
            <ConfirmDialog
                title={`Disconnect ${def.name}`}
                message={`Remove the connection to ${def.name}? You can reconnect at any time.`}
                confirmLabel="Disconnect"
                cancelLabel="Cancel"
                danger
                onConfirm={async () => { await disconnect(key) }}
            />
        )
    }

    return (
        <div className={cls.Overlay}>
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
                <div className={cls.ModalHead}>
                    <div className={cls.ModalHeadLeft}>
                        <div className={cls.ModalIcon}>{def.icon}</div>
                        <span className={cls.ModalTitle}>{def.name}</span>
                    </div>
                    <button className={cls.ModalClose} onClick={CloseDialog} aria-label="Close">
                        <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor"
                             strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <line x1="18" y1="6" x2="6" y2="18"/>
                            <line x1="6" y1="6" x2="18" y2="18"/>
                        </svg>
                    </button>
                </div>

                {connection ? (
                    def.provider === ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS
                        ? <GoogleSheetsConnectedContent connection={connection} onDisconnect={handleDisconnect}/>
                        : <ConnectedContent connection={connection} onDisconnect={handleDisconnect}/>
                ) : (
                    <NotConnectedContent def={def} connecting={connecting} onConnectGoogle={handleConnectGoogle}/>
                )}
            </div>
        </div>
    )
}

const SCOPE_INFO: Record<string, string> = {
    "spreadsheets": "Read and write your Google Sheets spreadsheets",
    "spreadsheets.readonly": "Read-only access to your Google Sheets",
    "drive.file": "Access files Artel creates or opens in Google Drive",
    "drive.readonly": "Read-only access to your Google Drive files",
    "gmail.readonly": "Read your Gmail messages",
    "calendar.readonly": "Read your Google Calendar events",
    "calendar": "Read and write your Google Calendar events",
}

function trimScope(s: string): string {
    const idx = s.lastIndexOf("/")
    return idx >= 0 ? s.slice(idx + 1) : s
}

function parseScopeList(scopes: string): string[] {
    return scopes.split(/[ ,]+/).filter(Boolean)
}

function ConnectedContent({connection, onDisconnect}: {
    connection: ExternalConnectionInfo
    onDisconnect: () => void
}) {
    const email = connection.google?.email
    const scopes = connection.google?.scopes
    const connectedAt = connection.createdAt
        ? new Date(connection.createdAt).toLocaleDateString(undefined, {year: "numeric", month: "short", day: "numeric"})
        : "—"

    const scopeList = scopes ? parseScopeList(scopes) : []
    const scopeTooltipHtml = scopeList
        .map(s => {
            const name = trimScope(s)
            const desc = SCOPE_INFO[name]
            return desc ? `<b>${name}</b>: ${desc}` : `<b>${name}</b>`
        })
        .join("<br/>")

    return (
        <>
            <div className={cls.InfoRows}>
                {email && (
                    <div className={cls.InfoRow}>
                        <span className={cls.InfoLabel}>Account</span>
                        <span className={cls.InfoValue}>{email}</span>
                    </div>
                )}
                {scopeList.length > 0 && (
                    <div className={cls.InfoRow}>
                        <span className={`${cls.InfoLabel} ${cls.ScopesLabel}`}>
                            Scopes
                            <button
                                type="button"
                                className={cls.ScopeHelp}
                                data-tooltip-id="root-tooltip"
                                data-tooltip-html={scopeTooltipHtml}
                                aria-label="What are scopes?"
                            >?</button>
                        </span>
                        <span className={cls.InfoValue}>{scopeList.map(trimScope).join(", ")}</span>
                    </div>
                )}
                <div className={cls.InfoRow}>
                    <span className={cls.InfoLabel}>Connected</span>
                    <span className={cls.InfoValue}>{connectedAt}</span>
                </div>
            </div>
            <div className={cls.ModalActions}>
                <button className={cls.BtnDanger} onClick={onDisconnect}>Disconnect</button>
            </div>
        </>
    )
}

function GoogleSheetsConnectedContent({connection, onDisconnect}: {
    connection: ExternalConnectionInfo
    onDisconnect: () => void
}) {
    const {OpenDialog} = useDialog()
    const bakeError = useBakeError()
    const {spreadsheets, spreadsheetsLoading, fetchSpreadsheets, addSpreadsheet, removeSpreadsheet, getPickerToken} = useExternalConnections()
    const pickerOpenRef = useRef(false)

    const email = connection.google?.email
    const scopes = connection.google?.scopes
    const connectedAt = connection.createdAt
        ? new Date(connection.createdAt).toLocaleDateString(undefined, {year: "numeric", month: "short", day: "numeric"})
        : "—"

    const scopeList = scopes ? parseScopeList(scopes) : []
    const scopeTooltipHtml = scopeList
        .map(s => {
            const name = trimScope(s)
            const desc = SCOPE_INFO[name]
            return desc ? `<b>${name}</b>: ${desc}` : `<b>${name}</b>`
        })
        .join("<br/>")

    useEffect(() => {
        void fetchSpreadsheets()
    }, [])

    async function openPicker() {
        if (pickerOpenRef.current) return
        pickerOpenRef.current = true
        try {
            const token = await getPickerToken()
            await loadGapi()
            const picker = new google.picker.PickerBuilder()
                .addView(new google.picker.DocsView(google.picker.ViewId.SPREADSHEETS))
                .setOAuthToken(token)
                .setDeveloperKey(import.meta.env.VITE_GOOGLE_API_KEY ?? "")
                .setCallback(async (data: google.picker.ResponseObject) => {
                    const action = data[google.picker.Response.ACTION]
                    if (action === google.picker.Action.PICKED) {
                        const docs = data[google.picker.Response.DOCUMENTS]
                        const file = docs?.[0]
                        if (file) {
                            try {
                                await addSpreadsheet(
                                    file[google.picker.Document.ID],
                                    file[google.picker.Document.NAME] ?? file[google.picker.Document.ID],
                                )
                            } catch (e) {
                                bakeError("Failed to add spreadsheet", e)
                            }
                        }
                    }
                    if (action === google.picker.Action.PICKED || action === google.picker.Action.CANCEL) {
                        pickerOpenRef.current = false
                    }
                })
                .build()
            picker.setVisible(true)
        } catch (e) {
            pickerOpenRef.current = false
            bakeError("Failed to open picker", e)
        }
    }

    function handleRemove(sheet: Spreadsheet) {
        OpenDialog(
            <ConfirmDialog
                title="Remove spreadsheet"
                message={`Remove "${sheet.name}" from MCP access? The file is not deleted from Google Drive.`}
                confirmLabel="Remove"
                cancelLabel="Cancel"
                danger
                onConfirm={async () => { await removeSpreadsheet(sheet.id ?? "") }}
            />
        )
    }

    return (
        <>
            <div className={cls.InfoRows}>
                {email && (
                    <div className={cls.InfoRow}>
                        <span className={cls.InfoLabel}>Account</span>
                        <span className={cls.InfoValue}>{email}</span>
                    </div>
                )}
                {scopeList.length > 0 && (
                    <div className={cls.InfoRow}>
                        <span className={`${cls.InfoLabel} ${cls.ScopesLabel}`}>
                            Scopes
                            <button
                                type="button"
                                className={cls.ScopeHelp}
                                data-tooltip-id="root-tooltip"
                                data-tooltip-html={scopeTooltipHtml}
                                aria-label="What are scopes?"
                            >?</button>
                        </span>
                        <span className={cls.InfoValue}>{scopeList.map(trimScope).join(", ")}</span>
                    </div>
                )}
                <div className={cls.InfoRow}>
                    <span className={cls.InfoLabel}>Connected</span>
                    <span className={cls.InfoValue}>{connectedAt}</span>
                </div>
            </div>

            <div className={cls.SpreadsheetSection}>
                <div className={cls.SpreadsheetSectionHeader}>
                    <span className={cls.SpreadsheetSectionLabel}>Spreadsheets</span>
                    <button className={cls.BtnGhost} onClick={openPicker}>+ Add</button>
                </div>
                {spreadsheetsLoading ? (
                    <p className={cls.EmptySpreadsheets}>Loading…</p>
                ) : spreadsheets.length === 0 ? (
                    <p className={cls.EmptySpreadsheets}>No spreadsheets added yet.</p>
                ) : (
                    <div className={cls.SpreadsheetList}>
                        {spreadsheets.map(sheet => (
                            <div key={sheet.id} className={cls.SpreadsheetRow}>
                                <span className={cls.SpreadsheetName}>{sheet.name}</span>
                                <button className={cls.BtnIconRemove} onClick={() => handleRemove(sheet)} aria-label="Remove">
                                    ×
                                </button>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            <div className={cls.ModalActions}>
                <button className={cls.BtnDanger} onClick={onDisconnect}>Disconnect</button>
            </div>
        </>
    )
}

function NotConnectedContent({def, connecting, onConnectGoogle}: {
    def: ProviderDef
    connecting: boolean
    onConnectGoogle: () => void
}) {
    return (
        <>
            <p className={cls.ModalDesc}>{def.description}</p>
            {def.canConnect ? (
                <div className={cls.ModalActions}>
                    <button className={cls.BtnPrimary} onClick={onConnectGoogle} disabled={connecting}>
                        {connecting ? "Redirecting…" : `Connect with ${def.name}`}
                    </button>
                </div>
            ) : (
                <p className={cls.ComingSoon}>Coming soon — this integration is not yet available.</p>
            )}
        </>
    )
}

/* ── Email connector ────────────────────────────────────────── */

function EmailCard({accounts, loading, onClick}: {
    accounts: EmailAccountInfo[]
    loading: boolean
    onClick?: () => void
}) {
    const isConnected = accounts.length > 0
    const accountLabel = accounts.length === 1
        ? accounts[0].email
        : accounts.length > 1 ? `${accounts.length} accounts` : undefined

    return (
        <div
            className={cls.Card}
            onClick={!loading ? onClick : undefined}
            role="button"
            tabIndex={0}
            onKeyDown={e => { if (e.key === "Enter" || e.key === " ") onClick?.() }}
        >
            <div className={cls.CardHeader}>
                <div className={cls.CardIcon}><EmailIcon/></div>
                <div className={cls.CardTitles}>
                    <div className={cls.CardName}>Email</div>
                    {accountLabel && <div className={cls.CardAccount}>{accountLabel}</div>}
                </div>
            </div>
            <div className={cls.CardFooter}>
                <span className={`${cls.StatusDot} ${isConnected ? cls.StatusDotConnected : cls.StatusDotDisconnected}`}/>
                <span className={`${cls.StatusLabel} ${isConnected ? cls.StatusLabelConnected : cls.StatusLabelDisconnected}`}>
                    {loading ? "…" : isConnected ? "Connected" : "Not connected"}
                </span>
            </div>
        </div>
    )
}

function EmailDetailDialog() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {accounts, loading, remove} = useEmailAccounts()
    const {isEmailsEnabled} = useUser()
    const bakeError = useBakeError()

    function handleRemove(account: EmailAccountInfo) {
        OpenDialog(
            <ConfirmDialog
                title="Remove email account"
                message={`Remove "${account.email}"? You can re-add it at any time.`}
                confirmLabel="Remove"
                cancelLabel="Cancel"
                danger
                onConfirm={async () => {
                    try {
                        await remove(account.id ?? "")
                    } catch (e) {
                        bakeError("Failed to remove account", e)
                    }
                }}
            />
        )
    }

    return (
        <div className={cls.Overlay}>
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
                <div className={cls.ModalHead}>
                    <div className={cls.ModalHeadLeft}>
                        <div className={cls.ModalIcon}><EmailIcon/></div>
                        <span className={cls.ModalTitle}>Email</span>
                    </div>
                    <button className={cls.ModalClose} onClick={CloseDialog} aria-label="Close">
                        <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor"
                             strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <line x1="18" y1="6" x2="6" y2="18"/>
                            <line x1="6" y1="6" x2="18" y2="18"/>
                        </svg>
                    </button>
                </div>

                {!isEmailsEnabled ? (
                    <p className={cls.ModalDesc}>Your account does not have access to the Emails feature.</p>
                ) : (
                    <div className={cls.EmailAccountSection}>
                        <div className={cls.EmailAccountSectionHeader}>
                            <span className={cls.EmailAccountSectionLabel}>Accounts</span>
                            <button className={cls.BtnGhost} onClick={() => OpenDialog(<EmailAddDialog/>)}>
                                + Add
                            </button>
                        </div>
                        {loading ? (
                            <p className={cls.EmptySpreadsheets}>Loading…</p>
                        ) : accounts.length === 0 ? (
                            <p className={cls.EmptySpreadsheets}>No email accounts added yet.</p>
                        ) : (
                            <div className={cls.EmailAccountList}>
                                {accounts.map(account => (
                                    <EmailAccountRow key={account.id} account={account} onRemove={handleRemove}/>
                                ))}
                            </div>
                        )}
                    </div>
                )}
            </div>
        </div>
    )
}

function EmailAccountRow({account, onRemove}: {
    account: EmailAccountInfo
    onRemove: (account: EmailAccountInfo) => void
}) {
    return (
        <div className={cls.EmailAccountRow}>
            <div className={cls.EmailAccountInfo}>
                <span className={cls.EmailAccountEmail}>{account.email}</span>
                <span className={cls.EmailAccountMeta}>
                    IMAP {account.imapHost}:{account.imapPort} · SMTP {account.smtpHost}:{account.smtpPort}
                </span>
            </div>
            <button className={cls.BtnIconRemove} onClick={() => onRemove(account)} aria-label="Remove">
                ×
            </button>
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

    const {add} = useEmailAccounts()
    const {CloseDialog} = useDialog()
    const {auth} = useUser()
    const bakeError = useBakeError()

    function handleEmailChange(value: string) {
        setEmail(value)
        const domain = value.split("@")[1] ?? ""
        if (!domain) return

        if (debounceRef.current) clearTimeout(debounceRef.current)
        debounceRef.current = setTimeout(async () => {
            try {
                const resp = await EmailAccountsAPI.ListMailServerSuggestions({domain}, auth.getInitReq())
                if (resp.suggestions?.length === 1) {
                    const s = resp.suggestions[0]
                    setImapHost(s.imap ?? "")
                    setImapPort(s.imapPort ? String(s.imapPort) : "")
                    setSmtpHost(s.smtp ?? "")
                    setSmtpPort(s.smtpPort ? String(s.smtpPort) : "")
                }
            } catch {
                // suggestion lookup is best-effort
            }
        }, 300)
    }

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
            .catch(e => bakeError("Error adding email", e))
            .finally(() => setAdding(false))
    }

    return (
        <div className={cls.Overlay}>
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                 aria-labelledby="addEmailModalTitle">
                <div className={cls.ModalHead}>
                    <div className={cls.ModalHeadLeft}>
                        <div className={cls.ModalIcon}><EmailIcon/></div>
                        <span className={cls.ModalTitle} id="addEmailModalTitle">Add email account</span>
                    </div>
                    <ModalClose onClick={CloseDialog} disabled={adding} className={cls.ModalClose}/>
                </div>
                <p className={cls.ModalSub}>Connect an email account via IMAP and SMTP.</p>

                <label className={cls.Field}>
                    <span className={cls.FieldLabel}>Email address</span>
                    <input
                        className={cls.Input}
                        placeholder="you@example.com"
                        value={email}
                        onChange={e => handleEmailChange(e.target.value)}
                        disabled={adding}
                        autoComplete="off"
                    />
                </label>
                <div className={cls.FieldRow}>
                    <label className={cls.Field}>
                        <span className={cls.FieldLabel}>IMAP host</span>
                        <input
                            className={cls.Input}
                            placeholder="imap.example.com"
                            value={imapHost}
                            onChange={e => setImapHost(e.target.value)}
                            disabled={adding}
                            autoComplete="off"
                        />
                    </label>
                    <label className={cls.FieldNarrow}>
                        <span className={cls.FieldLabel}>IMAP port</span>
                        <input
                            className={cls.Input}
                            placeholder="993"
                            value={imapPort}
                            onChange={e => setImapPort(e.target.value)}
                            disabled={adding}
                            maxLength={5}
                            autoComplete="off"
                        />
                    </label>
                </div>
                <div className={cls.FieldRow}>
                    <label className={cls.Field}>
                        <span className={cls.FieldLabel}>SMTP host</span>
                        <input
                            className={cls.Input}
                            placeholder="smtp.example.com"
                            value={smtpHost}
                            onChange={e => setSmtpHost(e.target.value)}
                            disabled={adding}
                            autoComplete="off"
                        />
                    </label>
                    <label className={cls.FieldNarrow}>
                        <span className={cls.FieldLabel}>SMTP port</span>
                        <input
                            className={cls.Input}
                            placeholder="587"
                            value={smtpPort}
                            onChange={e => setSmtpPort(e.target.value)}
                            disabled={adding}
                            maxLength={5}
                            autoComplete="off"
                        />
                    </label>
                </div>
                <label className={cls.Field}>
                    <span className={cls.FieldLabel}>Password</span>
                    <input
                        type="password"
                        placeholder="App-specific password"
                        className={cls.Input}
                        onChange={e => setPassword(e.target.value)}
                        disabled={adding}
                        autoComplete="new-password"
                    />
                </label>

                <div className={cls.ModalActions}>
                    <button
                        className={cls.BtnPrimary}
                        onClick={handleAdd}
                        disabled={adding}
                    >
                        {adding ? "Adding…" : "Add account"}
                    </button>
                </div>
            </div>
        </div>
    )
}

/* ── Provider icons ─────────────────────────────────────────── */

function GoogleSheetsIcon() {
    return (
        <svg width="20" height="20" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect x="10" y="2" width="28" height="44" rx="2" fill="#23A566"/>
            <rect x="26" y="2" width="12" height="12" fill="#159C53"/>
            <path d="M26 2L38 14H26V2Z" fill="#82C9A7"/>
            <rect x="16" y="20" width="16" height="2" rx="1" fill="white" fillOpacity="0.9"/>
            <rect x="16" y="25" width="16" height="2" rx="1" fill="white" fillOpacity="0.9"/>
            <rect x="16" y="30" width="10" height="2" rx="1" fill="white" fillOpacity="0.9"/>
        </svg>
    )
}

function TrelloIcon() {
    return (
        <svg width="20" height="20" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect width="48" height="48" rx="6" fill="#0052CC"/>
            <rect x="7" y="9" width="14" height="20" rx="2" fill="white"/>
            <rect x="27" y="9" width="14" height="13" rx="2" fill="white"/>
        </svg>
    )
}

function MiroIcon() {
    return (
        <svg width="20" height="20" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect width="48" height="48" rx="6" fill="#FFD02F"/>
            <path d="M32 10H26L30 24L22 10H16L20 24L12 10H6L16 38H22L18 24L26 38H32L28 24L36 38H42L32 10Z" fill="#050038"/>
        </svg>
    )
}

function EmailIcon() {
    return (
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"
             stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
            <rect x="2" y="4" width="20" height="16" rx="2"/>
            <path d="M2 7l10 7 10-7"/>
        </svg>
    )
}
