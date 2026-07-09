import cls from "@/components/McpKeyCard/CardChips.module.css"
import {McpConnectorInfo} from "@/app/api/artel/mcp_keys.pb.ts"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {ExternalConnectionInfo} from "@/app/api/artel/external_connections.pb.ts"
import ConnectorChip from "@/components/ConnectorChip/ConnectorChip.tsx"
import VaultChip from "@/components/VaultChip/VaultChip.tsx"

interface Props {
    vault?: VaultItem
    connectors: McpConnectorInfo[]
    connections: ExternalConnectionInfo[]
}

export default function CardChips({vault, connectors, connections}: Props) {
    return (
        <div className={cls.CardChipsContainer}>
            <VaultChip vault={vault}/>

            {connectors.map(c => (
                <ConnectorChip key={c.mcpName} connector={c} connections={connections}/>
            ))}
        </div>
    )
}
