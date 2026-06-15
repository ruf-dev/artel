import {useEffect, useState} from "react"
import {useNavigate, useSearchParams} from "react-router-dom"

import cls from "@/pages/connections/ConnectionsPage.module.css"

import {ExternalConnectionInfo, ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import useUser from "@/hooks/user/User.ts"

import ConfirmDialog from "@/components/ConfirmDialog/ConfirmDialog.tsx"

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
    const {OpenDialog} = useDialog()

    if (loading) {
        return (
            <div className={cls.Content}>
                <div className={cls.Grid}>
                    {PROVIDERS.map(def => (
                        <ProviderCard key={def.provider} def={def} connection={undefined} loading/>
                    ))}
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
                    <ConnectedContent connection={connection} onDisconnect={handleDisconnect}/>
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
