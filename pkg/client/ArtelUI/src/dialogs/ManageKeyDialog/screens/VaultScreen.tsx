import {Button} from "@vervstack/chures"

import cls from "@/dialogs/ManageKeyDialog/ManageKeyDialog.module.css"
import DialogHead from "@/dialogs/ManageKeyDialog/components/DialogHead/DialogHead.tsx"
import VaultOptionList from "@/dialogs/ManageKeyDialog/components/VaultOptionList/VaultOptionList.tsx"

interface VaultScreenProps {
    selectedVaultId: string
    saving: boolean
    onSelectVault: (id: string) => void
    onBack: () => void
    onSave: () => void
}

export default function VaultScreen({selectedVaultId, saving, onSelectVault, onBack, onSave}: VaultScreenProps) {
    return (
        <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
             aria-labelledby="manageVaultTitle">
            <DialogHead titleId="manageVaultTitle" title="Select vault" disabled={saving}/>
            <p className={cls.ModalSub}>Choose which vault this key connects to.</p>
            <VaultOptionList selectedVaultId={selectedVaultId} onSelect={onSelectVault}/>
            <div className={cls.ModalActions}>
                <Button variant="ghost" onClick={onBack} disabled={saving}>Back</Button>
                <Button variant="primary" onClick={onSave} disabled={saving || !selectedVaultId}>
                    {saving ? "Saving…" : "Save"}
                </Button>
            </div>
        </div>
    )
}
