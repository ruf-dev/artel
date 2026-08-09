import {Button} from "@vervstack/chures"

import cls
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/dialogs/ManageVaultDialog/widgets/VaultSettingsSection/components/PostgresSection/PostgresSection.module.css"

type StatusType = "not_enabled" | "provisioning" | "ready" | "error"

interface PostgresActionsProps {
    currentStatus: StatusType
    enabling: boolean
    disabling: boolean
    onEnable: () => void
    onDisable: () => void
}

export default function PostgresActions({
    currentStatus,
    enabling,
    disabling,
    onEnable,
    onDisable,
}: PostgresActionsProps) {
    return (
        <div className={cls.Actions}>
            {currentStatus === "not_enabled" && (
                <Button
                    variant="primary"
                    onClick={onEnable}
                    disabled={enabling}
                >
                    {enabling ? "Enabling…" : "Enable Postgres"}
                </Button>
            )}
            {currentStatus === "provisioning" && (
                <span className={cls.ProvisioningText}>Provisioning database…</span>
            )}
            {currentStatus === "ready" && (
                <Button
                    variant="danger"
                    onClick={onDisable}
                    disabled={disabling}
                >
                    {disabling ? "Disabling…" : "Disable"}
                </Button>
            )}
            {currentStatus === "error" && (
                <Button
                    variant="ghost"
                    onClick={onEnable}
                    disabled={enabling}
                >
                    {enabling ? "Retrying…" : "Retry"}
                </Button>
            )}
        </div>
    )
}
