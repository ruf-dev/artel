import {Button} from "@vervstack/chures"

import {McpConnectorInfo} from "@/app/api/artel/mcp_keys.pb.ts"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import ConnectorChip from "@/components/ConnectorChip/ConnectorChip.tsx"
import cls from "@/dialogs/ManageKeyDialog/components/ConnectionsField/components/ConnectorRow/ConnectorRow.module.css"

interface ConnectorRowProps {
    connector: McpConnectorInfo
    mcpKeyId: string
    onEdit: () => void
}

export default function ConnectorRow({connector, mcpKeyId, onEdit}: ConnectorRowProps) {
    const {connections} = useExternalConnections()
    const {removeConnector} = useMcpKeys()

    async function onRevokeClick() {
        await removeConnector(mcpKeyId, connector.mcpName ?? "")
    }

    return (
        <div
            className={cls.ConnectorRowWrapper}
            onClick={onEdit}
        >
            <ConnectorChip
                connector={connector}
                connections={connections}
            />
            <div
                className={cls.ConnectorRowActions}
            >
                <Button
                    variant="ghost"
                    className={cls.BtnSmall}
                    onClick={onRevokeClick}>
                    Remove
                </Button>
            </div>
        </div>
    )
}
