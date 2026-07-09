import {McpConnectorInfo} from "@/app/api/artel/mcp_keys.pb.ts"
import ConnectorRow from "@/dialogs/ManageKeyDialog/components/ConnectionsField/components/ConnectorRow/ConnectorRow.tsx"
import cls from "@/dialogs/ManageKeyDialog/components/ConnectionsField/ConnectionsField.module.css"

interface ConnectionsFieldProps {
    connectors: McpConnectorInfo[]
    mcpKeyId: string
    onEdit: (connector: McpConnectorInfo) => void
}

export default function ConnectionsField({connectors, mcpKeyId, onEdit}: ConnectionsFieldProps) {
    return (
        <div className={cls.Field}>
            <span className={cls.FieldLabel}>Connections</span>
            <div className={cls.OptionList}>
                {connectors.length == 0 ?
                    <p className={cls.Empty}>
                        No connections linked to this key yet.
                    </p>
                    :
                    connectors.map(c => (
                        <ConnectorRow
                            key={c.mcpName}
                            connector={c}
                            mcpKeyId={mcpKeyId}
                            onEdit={() => onEdit(c)}/>
                    ))}
            </div>
        </div>
    )
}
