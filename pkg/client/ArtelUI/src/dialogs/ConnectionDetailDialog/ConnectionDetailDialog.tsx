import {useEffect, useRef, useState} from "react"

import cls from "@/dialogs/ConnectionDetailDialog/ConnectionDetailDialog.module.css"

import {ExternalConnectionInfo, ExternalProvider, Spreadsheet} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

import Button from "@/components/shared/Button/Button.tsx"
import ConfirmDialog from "@/components/ConfirmDialog/ConfirmDialog.tsx"
import ModalClose from "@/components/ModalClose/ModalClose.tsx"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"

type GapiWindow = Window & { gapi: { load: (lib: string, cb: () => void) => void } }

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

type ProviderConfig = {
    name: string
    description: string
    canConnect: boolean
}

const PROVIDER_CONFIG: Partial<Record<ExternalProvider, ProviderConfig>> = {
    [ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS]: {
        name: "Google Sheets",
        description: "Read and write data from your Google Sheets spreadsheets.",
        canConnect: true,
    },
    [ExternalProvider.EXTERNAL_PROVIDER_TRELLO]: {
        name: "Trello",
        description: "Sync tasks and boards from your Trello workspace.",
        canConnect: false,
    },
    [ExternalProvider.EXTERNAL_PROVIDER_MIRO]: {
        name: "Miro",
        description: "Access and embed your Miro boards.",
        canConnect: false,
    },
}

function providerKey(p: ExternalProvider): string {
    const map: Partial<Record<ExternalProvider, string>> = {
        [ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS]: "google_sheets",
        [ExternalProvider.EXTERNAL_PROVIDER_TRELLO]: "trello",
        [ExternalProvider.EXTERNAL_PROVIDER_MIRO]: "miro",
        [ExternalProvider.EXTERNAL_PROVIDER_EMAIL]: "email",
    }
    return map[p] ?? ""
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

function DialogHead({title, provider, onClose}: { title: string; provider: ExternalProvider; onClose: () => void }) {
    return (
        <div className={cls.ModalHead}>
            <div className={cls.ModalHeadLeft}>
                <div className={cls.ModalIcon}>
                    <ProviderIcon provider={provider}/>
                </div>
                <span className={cls.ModalTitle}>{title}</span>
            </div>
            <ModalClose onClick={onClose} className={cls.ModalClose}/>
        </div>
    )
}

export default function ConnectionDetailDialog({provider}: { provider: ExternalProvider }) {
    const {CloseDialog, OpenDialog} = useDialog()
    const {connections, disconnect, initiateGoogleOAuth} = useExternalConnections()
    const bakeError = useBakeError()
    const [connecting, setConnecting] = useState(false)

    const config = PROVIDER_CONFIG[provider]
    const connection = connections.find(c => c.provider === provider)

    function handleConnectGoogle() {
        setConnecting(true)
        initiateGoogleOAuth()
            .then(authUrl => {
                window.location.href = authUrl
            })
            .catch(e => {
                bakeError("Failed to start OAuth", e)
                setConnecting(false)
            })
    }

    function handleDisconnect() {
        const key = providerKey(provider)
        OpenDialog(
            <ConfirmDialog
                title={`Disconnect ${config?.name ?? ""}`}
                message={`Remove the connection to ${config?.name ?? ""}? You can reconnect at any time.`}
                confirmLabel="Disconnect"
                cancelLabel="Cancel"
                danger
                onConfirm={() => disconnect(key).catch(e => bakeError("Failed to disconnect", e))}
            />
        )
    }

    return (
        <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <DialogHead title={config?.name ?? ""} provider={provider} onClose={CloseDialog}/>

            {connection ? (
                provider === ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS
                    ? <GoogleSheetsConnectedContent connection={connection} onDisconnect={handleDisconnect}/>
                    : <ConnectedContent connection={connection} onDisconnect={handleDisconnect}/>
            ) : (
                <NotConnectedContent
                    description={config?.description ?? ""}
                    canConnect={config?.canConnect ?? false}
                    name={config?.name ?? ""}
                    connecting={connecting}
                    onConnectGoogle={handleConnectGoogle}
                />
            )}
        </div>
    )
}

function ConnectedContent({connection, onDisconnect}: {
    connection: ExternalConnectionInfo
    onDisconnect: () => void
}) {
    const email = connection.google?.email
    const scopes = connection.google?.scopes
    const connectedAt = connection.createdAt
        ? new Date(connection.createdAt).toLocaleDateString(undefined, {
            year: "numeric",
            month: "short",
            day: "numeric"
        })
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
                <Button variant="danger" onClick={onDisconnect}>Disconnect</Button>
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
    const {
        spreadsheets,
        spreadsheetsLoading,
        fetchSpreadsheets,
        addSpreadsheet,
        removeSpreadsheet,
        getPickerToken
    } = useExternalConnections()
    const pickerOpenRef = useRef(false)

    const email = connection.google?.email
    const scopes = connection.google?.scopes
    const connectedAt = connection.createdAt
        ? new Date(connection.createdAt).toLocaleDateString(undefined, {
            year: "numeric",
            month: "short",
            day: "numeric"
        })
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

    function openPicker() {
        if (pickerOpenRef.current) return
        pickerOpenRef.current = true

        getPickerToken()
            .then(token => loadGapi().then(() => token))
            .then(token => {
                const picker = new google.picker.PickerBuilder()
                    .addView(new google.picker.DocsView(google.picker.ViewId.SPREADSHEETS))
                    .setOAuthToken(token)
                    .setDeveloperKey(import.meta.env.VITE_GOOGLE_API_KEY ?? "")
                    .setCallback((data: google.picker.ResponseObject) => {
                        const action = data[google.picker.Response.ACTION]
                        if (action === google.picker.Action.PICKED) {
                            const docs = data[google.picker.Response.DOCUMENTS]
                            const file = docs?.[0]
                            if (file) {
                                addSpreadsheet(
                                    file[google.picker.Document.ID],
                                    file[google.picker.Document.NAME] ?? file[google.picker.Document.ID],
                                ).catch(e => bakeError("Failed to add spreadsheet", e))
                            }
                        }
                        if (action === google.picker.Action.PICKED || action === google.picker.Action.CANCEL) {
                            pickerOpenRef.current = false
                        }
                    })
                    .build()
                picker.setVisible(true)
            })
            .catch(e => {
                pickerOpenRef.current = false
                bakeError("Failed to open picker", e)
            })
    }

    function handleRemove(sheet: Spreadsheet) {
        OpenDialog(
            <ConfirmDialog
                title="Remove spreadsheet"
                message={`Remove "${sheet.name}" from MCP access? The file is not deleted from Google Drive.`}
                confirmLabel="Remove"
                cancelLabel="Cancel"
                danger
                onConfirm={() => removeSpreadsheet(sheet.id ?? "").catch(e => bakeError("Failed to remove spreadsheet", e))}
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
                    <Button variant="ghost" onClick={openPicker}>+ Add</Button>
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
                                <Button variant="iconDanger" onClick={() => handleRemove(sheet)}
                                        aria-label="Remove">×</Button>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            <div className={cls.ModalActions}>
                <Button variant="danger" onClick={onDisconnect}>Disconnect</Button>
            </div>
        </>
    )
}

function NotConnectedContent({description, canConnect, name, connecting, onConnectGoogle}: {
    description: string
    canConnect: boolean
    name: string
    connecting: boolean
    onConnectGoogle: () => void
}) {
    return (
        <>
            <p className={cls.ModalDesc}>{description}</p>
            {canConnect ? (
                <div className={cls.ModalActions}>
                    <Button variant="primary" onClick={onConnectGoogle} disabled={connecting}>
                        {connecting ? "Redirecting…" : `Connect with ${name}`}
                    </Button>
                </div>
            ) : (
                <p className={cls.ComingSoon}>Coming soon — this integration is not yet available.</p>
            )}
        </>
    )
}
