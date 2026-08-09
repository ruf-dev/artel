import {useEffect, useState} from "react"
import {ConfirmDialog} from "@vervstack/chures"

import cls
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/dialogs/ManageVaultDialog/widgets/VaultSettingsSection/components/PostgresSection/PostgresSection.module.css"
import {useDialog} from "@/app/hooks/Dialog.ts"
import {useVaultMutations} from "@/app/hooks/Vaults.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import PostgresActions
    from "@/dialogs/ManageVaultDialog/widgets/VaultSettingsSection/components/PostgresSection/PostgresActions.tsx"

interface PostgresSectionProps {
    vaultId: string
    postgresStatus?: string
    onStatusChange?: (status: string) => void
}

type StatusType = "not_enabled" | "provisioning" | "ready" | "error"

export default function PostgresSection({vaultId, postgresStatus, onStatusChange}: PostgresSectionProps) {
    const {OpenDialog, CloseDialog} = useDialog()
    const {enablePostgresDatabase, getPostgresDatabase, disablePostgresDatabase} = useVaultMutations()
    const bakeError = useBakeError()

    const [enabling, setEnabling] = useState(false)
    const [disabling, setDisabling] = useState(false)
    const [currentStatus, setCurrentStatus] = useState<StatusType>((postgresStatus as StatusType) || "not_enabled")
    const [errorMessage, setErrorMessage] = useState("")

    // Poll for status while provisioning
    useEffect(() => {
        if (currentStatus !== "provisioning") {
            return
        }

        const interval = setInterval(() => {
            getPostgresDatabase(vaultId)
                .then(res => {
                    setCurrentStatus((res.status as StatusType) || "not_enabled")
                    setErrorMessage(res.errorMessage ?? "")
                    onStatusChange?.(res.status ?? "not_enabled")
                })
                .catch(e => {
                    bakeError("Error checking Postgres status", e)
                })
        }, 2000)

        return () => clearInterval(interval)
    }, [currentStatus, vaultId, getPostgresDatabase, onStatusChange, bakeError])

    function handleEnable() {
        setEnabling(true)
        enablePostgresDatabase(vaultId)
            .then(res => {
                setCurrentStatus((res.status as StatusType) || "provisioning")
                setErrorMessage(res.errorMessage ?? "")
                onStatusChange?.(res.status ?? "provisioning")
            })
            .catch(e => bakeError("Failed to enable Postgres database", e))
            .finally(() => setEnabling(false))
    }

    function handleDisable() {
        OpenDialog(
            <ConfirmDialog
                title="Disable Postgres Database"
                message="This will remove the Postgres database and all its data. This action cannot be undone."
                confirmLabel="Disable"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() => {
                    setDisabling(true)
                    return disablePostgresDatabase(vaultId)
                        .then(() => {
                            setCurrentStatus("not_enabled")
                            onStatusChange?.("not_enabled")
                        })
                        .catch(e => bakeError("Failed to disable Postgres database", e))
                        .finally(() => setDisabling(false))
                }}
            />
        )
    }

    const statusLabel = {
        not_enabled: "Not enabled",
        provisioning: "Provisioning…",
        ready: "Ready",
        error: "Error",
    }[currentStatus]

    return (
        <div className={cls.PostgresSectionContainer}>
            <div className={cls.StatusHeader}>
                <span className={cls.Label}>PostgreSQL Database</span>
                <span className={cls.Status} data-status={currentStatus}>
                    {statusLabel}
                </span>
            </div>

            {currentStatus === "error" && errorMessage && (
                <p className={cls.ErrorMessage}>{errorMessage}</p>
            )}

            <PostgresActions
                currentStatus={currentStatus}
                enabling={enabling}
                disabling={disabling}
                onEnable={handleEnable}
                onDisable={handleDisable}
            />
        </div>
    )
}
