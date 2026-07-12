import {useState} from "react"
import {Button, ModalClose} from "@vervstack/chures"

import cls from "@/dialogs/ManageVaultDialog/ManageVaultDialog.module.css"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import ExpertSettingsSection from "@/dialogs/ManageVaultDialog/widgets/ExpertSettingsSection/ExpertSettingsSection.tsx"
import MembersSection from "@/dialogs/ManageVaultDialog/widgets/MembersSection/MembersSection.tsx"
import InviteLinksSection from "@/dialogs/ManageVaultDialog/widgets/InviteLinksSection/InviteLinksSection.tsx"
import VaultDangerZone from "@/widgets/VaultDangerZone/VaultDangerZone.tsx"

interface Props {
    vault: VaultItem
    currentUserId: string
    onClose: () => void
    onDeleted: () => void
}

export default function ManageVaultDialog({vault, currentUserId, onClose, onDeleted}: Props) {
    const [settingsOverride, setSettingsOverride] = useState<Partial<VaultItem> | null>(null)

    const effectiveVault: VaultItem = {...vault, ...settingsOverride}

    return (
        <div className={cls.ManageVaultDialogContainer}
             onClick={e => e.stopPropagation()}
             role="dialog"
             aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>Manage vault</h2>
                <ModalClose onClick={onClose} className={cls.ModalClose}/>
            </div>
            <p className={cls.VaultName}>{vault.name}</p>

            {vault.id && <MembersSection vaultId={vault.id} currentUserId={currentUserId}/>}

            {vault.id && <InviteLinksSection vaultId={vault.id}/>}

            {vault.id && <VaultDangerZone vaultId={vault.id} onDeleted={onDeleted}/>}

            {vault.id && (
                <ExpertSettingsSection vault={effectiveVault} onChanged={setSettingsOverride}/>
            )}

            <div className={cls.ModalFooter}>
                <Button variant="ghost" onClick={onClose}>Close</Button>
            </div>
        </div>
    )
}
