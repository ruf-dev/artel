import {useState} from "react"

import cls from "@/dialogs/ConnectionDetailDialog/ConnectionDetailDialog.module.css"

import {ExternalConnectionInfo, ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

import Button from "@/components/shared/Button/Button.tsx"
import ConfirmDialog from "@/components/ConfirmDialog/ConfirmDialog.tsx"
import ModalClose from "@/components/ModalClose/ModalClose.tsx"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"

import GenericConnectionContent from "@/widgets/GenericConnectionContent/GenericConnectionContent.tsx"
import GoogleConnectionContent from "@/widgets/GoogleConnectionContent/GoogleConnectionContent.tsx"
import GoogleSheetsConnectionContent from "@/widgets/GoogleSheetsConnectionContent/GoogleSheetsConnectionContent.tsx"

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

const PROVIDER_KEY: Partial<Record<ExternalProvider, string>> = {
    [ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS]: "google_sheets",
    [ExternalProvider.EXTERNAL_PROVIDER_TRELLO]: "trello",
    [ExternalProvider.EXTERNAL_PROVIDER_MIRO]: "miro",
    [ExternalProvider.EXTERNAL_PROVIDER_EMAIL]: "email",
}

function DialogHead({title, provider, onClose}: {
    title: string
    provider: ExternalProvider
    onClose: () => void
}) {
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

function NotConnectedContent({description, canConnect, name, connecting, onConnect}: {
    description: string
    canConnect: boolean
    name: string
    connecting: boolean
    onConnect: () => void
}) {
    return (
        <>
            <p className={cls.ModalDesc}>{description}</p>
            {canConnect ? (
                <div className={cls.ModalActions}>
                    <Button variant="primary" onClick={onConnect} disabled={connecting}>
                        {connecting ? "Redirecting…" : `Connect with ${name}`}
                    </Button>
                </div>
            ) : (
                <p className={cls.ComingSoon}>Coming soon — this integration is not yet available.</p>
            )}
        </>
    )
}

function ConnectedContent({provider, connection, onDisconnect}: {
    provider: ExternalProvider
    connection: ExternalConnectionInfo
    onDisconnect: () => void
}) {
    if (connection.google && provider === ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS) {
        return <GoogleSheetsConnectionContent connection={connection} onDisconnect={onDisconnect}/>
    }
    if (connection.google) {
        return <GoogleConnectionContent connection={connection} onDisconnect={onDisconnect}/>
    }
    return <GenericConnectionContent connection={connection} onDisconnect={onDisconnect}/>
}

export default function ConnectionDetailDialog({provider}: { provider: ExternalProvider }) {
    const {CloseDialog, OpenDialog} = useDialog()
    const {connections, disconnect, initiateGoogleOAuth} = useExternalConnections()
    const bakeError = useBakeError()
    const [connecting, setConnecting] = useState(false)

    const config = PROVIDER_CONFIG[provider]
    const connection = connections.find(c => c.provider === provider)

    function handleConnect() {
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
        OpenDialog(
            <ConfirmDialog
                title={`Disconnect ${config?.name ?? ""}`}
                message={`Remove the connection to ${config?.name ?? ""}? You can reconnect at any time.`}
                confirmLabel="Disconnect"
                cancelLabel="Cancel"
                danger
                onConfirm={() =>
                    disconnect(PROVIDER_KEY[provider] ?? "").catch(e => bakeError("Failed to disconnect", e))
                }
            />
        )
    }

    return (
        <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <DialogHead title={config?.name ?? ""} provider={provider} onClose={CloseDialog}/>
            {connection ? (
                <ConnectedContent provider={provider} connection={connection} onDisconnect={handleDisconnect}/>
            ) : (
                <NotConnectedContent
                    description={config?.description ?? ""}
                    canConnect={config?.canConnect ?? false}
                    name={config?.name ?? ""}
                    connecting={connecting}
                    onConnect={handleConnect}
                />
            )}
        </div>
    )
}
