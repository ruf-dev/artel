import {useState} from "react"

import cls from "@/dialogs/ManageVaultDialog/ManageVaultDialog.module.css"

import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import useUser from "@/hooks/user/User.ts"
import {Button, ModalClose} from "@vervstack/chures"
import ManageVaultS3Section, {S3LinkPatch} from "@/components/ManageVaultS3Section/ManageVaultS3Section.tsx"

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
    const {isAdmin} = useUser()

    const [s3Override, setS3Override] = useState<S3LinkPatch | null>(null)

    const effectiveVault: VaultItem = s3Override ? {...vault, ...s3Override} : vault

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

            {isAdmin && (
                <ManageVaultS3Section vault={effectiveVault} onLinked={setS3Override}/>
            )}

            {vault.id && <VaultDangerZone vaultId={vault.id} onDeleted={onDeleted}/>}

            <div className={cls.ModalFooter}>
                <Button variant="ghost" onClick={onClose}>Close</Button>
            </div>
        </div>
    )
}
