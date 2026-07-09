import {Button} from "@vervstack/chures"

import cls from "@/dialogs/ManageKeyDialog/ManageKeyDialog.module.css"
import {McpKeyInfo, McpConnectorInfo} from "@/app/api/artel/mcp_keys.pb.ts"
import DialogHead from "@/dialogs/ManageKeyDialog/components/DialogHead/DialogHead.tsx"
import VaultField from "@/dialogs/ManageKeyDialog/components/VaultField/VaultField.tsx"
import ConnectionsField from "@/dialogs/ManageKeyDialog/components/ConnectionsField/ConnectionsField.tsx"

interface MainScreenProps {
    mcpKey: McpKeyInfo
    selectedVaultId: string
    connectors: McpConnectorInfo[]
    onChangeVault: () => void
    onEditConnector: (connector: McpConnectorInfo) => void
    onAddConnection: () => void
    onRevoke: () => void
}

export default function MainScreen(props: MainScreenProps) {
    return (
        <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
             aria-labelledby="manageKeyTitle">
            <DialogHead
                titleId="manageKeyTitle"
                title={<>Manage mcp key: <b className={cls.KeyName}>{props.mcpKey.name}</b></>}
            />

            <VaultField
                selectedVaultId={props.selectedVaultId}
                onChangeVault={props.onChangeVault}
            />

            <ConnectionsField
                connectors={props.connectors}
                mcpKeyId={props.mcpKey.id ?? ""}
                onEdit={props.onEditConnector}
            />

            <div className={cls.ModalFooter}>
                <Button variant="danger" onClick={props.onRevoke}>
                    Revoke key
                </Button>

                <Button variant="primary" onClick={props.onAddConnection}>
                    Add connection
                </Button>
            </div>
        </div>
    )
}
