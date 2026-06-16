import {useEffect} from "react"

import cls from "@/widgets/McpKeyCard/McpKeyCard.module.css"

import {McpKeyInfo} from "@/app/api/artel/mcp_keys.pb.ts"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useVaults} from "@/app/hooks/Vaults.ts"

import ConnectorChip from "@/components/ConnectorChip/ConnectorChip.tsx"

interface Props {
    mcpKey: McpKeyInfo
    onRevoke: () => void
    onManage: () => void
}

export default function McpKeyCard({mcpKey, onRevoke, onManage}: Props) {
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

    function formatDate(iso: string | undefined): string {
        if (!iso) return "Never"
        const d = new Date(iso)
        if (isNaN(d.getTime())) return "Never"
        return d.toLocaleDateString(undefined, {year: "numeric", month: "short", day: "numeric"})
    }

    return (
        <div className={cls.Card}>
            <div className={cls.CardMain}>
                <div className={cls.CardHeader}>
                    <span className={cls.CardName}>{mcpKey.name}</span>
                    <span className={cls.CardPreview}>{mcpKey.keyPreview}…</span>
                </div>
                <div className={cls.CardChips}>
                    {vault ? (
                        <span className={cls.VaultChip} title={`Vault: ${vault.name}`}>
                            <span className={cls.VaultBadge}>A</span>
                            {vault.name}
                        </span>
                    ) : (
                        <span className={`${cls.Chip} ${cls.ChipMuted}`} title="No vault assigned">No vault</span>
                    )}
                    {connectors.map(c => (
                        <ConnectorChip key={c.mcpName} connector={c} connections={connections}/>
                    ))}
                </div>
                <div className={cls.CardMeta}>
                    Last accessed: {formatDate(mcpKey.lastAccessedAt)}
                </div>
            </div>
            <div className={cls.CardActions}>
                <button className={cls.BtnGhost} onClick={onManage} type="button">Manage</button>
                <button className={cls.BtnDanger} onClick={onRevoke} type="button">Revoke</button>
            </div>
        </div>
    )
}
