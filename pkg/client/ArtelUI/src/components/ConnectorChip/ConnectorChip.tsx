import {McpConnectorInfo} from "@/app/api/artel/mcp_keys.pb.ts"
import {ExternalConnectionInfo, ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import GenericChip from "@/components/GenericChip/GenericChip.tsx"
import EmailChip from "@/components/EmailChip/EmailChip.tsx"
import ProviderChip from "@/components/ProviderChip/ProviderChip.tsx"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"
import {PROVIDER_CHIP_CLASS} from "@/components/ProviderChip/providerChipClass.ts"

export default function ConnectorChip({connector, connections}: {
    connector: McpConnectorInfo
    connections: ExternalConnectionInfo[]
}) {
    const match = connections.find(c => c.id === connector.externalConnectionId)

    if (!match) {
        return <GenericChip mcpName={connector.mcpName}/>
    }

    if (match.provider === ExternalProvider.EXTERNAL_PROVIDER_EMAIL) {
        return <EmailChip label={connectionLabel(match)}/>
    }

    if (match.provider && PROVIDER_CHIP_CLASS[match.provider]) {
        return (
            <ProviderChip
                provider={match.provider}
                variantClass={PROVIDER_CHIP_CLASS[match.provider]!}
                label={connectionLabel(match)}
            />
        )
    }

    return <GenericChip mcpName={connector.mcpName}/>
}
