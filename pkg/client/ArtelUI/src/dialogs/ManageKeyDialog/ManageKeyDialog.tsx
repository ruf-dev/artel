import {McpKeyInfo} from "@/app/api/artel/mcp_keys.pb.ts"
import {useManageKeyDialog} from "@/dialogs/ManageKeyDialog/processes/useManageKeyDialog.tsx"
import MainScreen from "@/dialogs/ManageKeyDialog/screens/MainScreen.tsx"
import VaultScreen from "@/dialogs/ManageKeyDialog/screens/VaultScreen.tsx"
import AddConnectionScreen from "@/dialogs/ManageKeyDialog/screens/AddConnectionScreen.tsx"
import SelectConnectionScreen from "@/dialogs/ManageKeyDialog/screens/SelectConnectionScreen.tsx"

export default function ManageKeyDialog({mcpKey}: { mcpKey: McpKeyInfo }) {
    const dialog = useManageKeyDialog(mcpKey)

    if (dialog.step === "vault") {
        return (
            <VaultScreen
                selectedVaultId={dialog.selectedVaultId}
                saving={dialog.saving}
                onSelectVault={dialog.setSelectedVaultId}
                onBack={() => dialog.setStep("main")}
                onSave={dialog.handleSaveVault}
            />
        )
    }

    if (dialog.step === "addConnection") {
        return (
            <AddConnectionScreen
                connectors={dialog.connectors}
                onSelect={dialog.handleSelectCandidate}
                onBack={() => dialog.setStep("main")}
            />
        )
    }

    if (dialog.step === "selectConnection" && dialog.selectedCandidate) {
        return (
            <SelectConnectionScreen
                candidate={dialog.selectedCandidate}
                selectedExternalConnectionId={dialog.selectedExternalConnectionId}
                saving={dialog.saving}
                editing={!!dialog.editingConnectorName}
                onSelectConnection={dialog.setSelectedExternalConnectionId}
                onBack={() => dialog.setStep(dialog.editingConnectorName ? "main" : "addConnection")}
                onAdd={dialog.handleAddConnector}
            />
        )
    }

    return (
        <MainScreen
            mcpKey={mcpKey}
            selectedVaultId={dialog.selectedVaultId}
            connectors={dialog.connectors}
            onChangeVault={() => dialog.setStep("vault")}
            onEditConnector={dialog.handleEditConnector}
            onAddConnection={() => dialog.setStep("addConnection")}
            onRevoke={dialog.handleRevoke}
        />
    )
}
