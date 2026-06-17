import {useEffect} from "react"

import cls from "@/widgets/McpKeyCard/McpKeyCard.module.css"

import {McpConnectorInfo, McpKeyInfo} from "@/app/api/artel/mcp_keys.pb.ts"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {ExternalConnectionInfo} from "@/app/api/artel/external_connections.pb.ts"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useVaults} from "@/app/hooks/Vaults.ts"
import {useDialog} from "@/app/hooks/Dialog"
import useUser from "@/hooks/user/User.ts"

import ConnectorChip from "@/components/ConnectorChip/ConnectorChip.tsx"
import ManageVaultDialog from "@/components/ManageVaultDialog/ManageVaultDialog.tsx"

interface Props {
    mcpKey: McpKeyInfo
    onManage: () => void
}

export default function McpKeyCard({mcpKey, onManage}: Props) {
    const {vaults} = useVaults()
    const {connectorsByKey, fetchConnectors} = useMcpKeys()
    const {connections} = useExternalConnections()

    const vault = vaults.find(v => v.id === mcpKey.vaultId)
    const connectors = mcpKey.id ? connectorsByKey[mcpKey.id] ?? [] : []

    useEffect(() => {
        if (mcpKey.id) {
            void fetchConnectors(mcpKey.id)
        }
    }, [mcpKey.id, fetchConnectors])

    return (
        <div className={cls.Card} onClick={onManage}>
            <div className={cls.CardMain}>
                <CardHeader name={mcpKey.name} keyPreview={mcpKey.keyPreview}/>
                <CardChips vault={vault} connectors={connectors} connections={connections}/>
                <CardMeta lastAccessedAt={mcpKey.lastAccessedAt}/>
            </div>
        </div>
    )
}

function CardHeader({name, keyPreview}: {name?: string, keyPreview?: string}) {
    return (
        <div className={cls.CardHeader}>
            <span className={cls.CardName}>{name}</span>
            <span className={cls.CardPreview}>{keyPreview}…</span>
        </div>
    )
}

function CardChips({vault, connectors, connections}: {
    vault?: VaultItem
    connectors: McpConnectorInfo[]
    connections: ExternalConnectionInfo[]
}) {
    const {OpenDialog} = useDialog()

    function openVaultManage() {
        if (!vault) return
        const {CloseDialog} = useDialog.getState()
        const {auth} = useUser.getState()
        OpenDialog(
            <ManageVaultDialog
                vault={vault}
                currentUserId={auth.userInfo?.id ?? ""}
                onClose={CloseDialog}
                onDeleted={CloseDialog}
            />
        )
    }

    return (
        <div className={cls.CardChips}>
            {vault ? (
                <span
                    className={cls.VaultChip}
                    onClick={e => { e.stopPropagation(); openVaultManage() }}
                    data-tooltip-id="root-tooltip"
                    data-tooltip-content={`Vault: ${vault.name}`}
                >
                    <span className={cls.VaultBadge}>A</span>
                    {vault.name}
                </span>
            ) : (
                <span
                    className={`${cls.Chip} ${cls.ChipMuted}`}
                    data-tooltip-id="root-tooltip"
                    data-tooltip-content="No vault assigned"
                >
                    No vault
                </span>
            )}

            {connectors.map(c => (
                <ConnectorChip key={c.mcpName} connector={c} connections={connections}/>
            ))}
        </div>
    )
}

function CardMeta({lastAccessedAt}: {lastAccessedAt?: string}) {
    return (
        <div className={cls.CardMeta}>
            Last accessed: {formatDate(lastAccessedAt)}
        </div>
    )
}


function formatDate(iso: string | undefined): string {
    if (!iso) return "Never"
    const d = new Date(iso)
    if (isNaN(d.getTime())) return "Never"
    return d.toLocaleDateString(undefined, {year: "numeric", month: "short", day: "numeric"})
}
