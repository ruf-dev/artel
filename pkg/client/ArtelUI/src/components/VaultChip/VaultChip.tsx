import cls from "@/components/VaultChip/VaultChip.module.css"
import {cn} from "@/app/utils/cn"

import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import useUser from "@/hooks/user/User.ts"

import ManageVaultDialog from "@/dialogs/ManageVaultDialog/ManageVaultDialog.tsx"

export default function VaultChip({vault}: {vault?: VaultItem}) {
    const {OpenDialog} = useDialog()

    if (!vault) {
        return (
            <span
                className={cn(cls.Chip, cls.ChipMuted)}
                data-tooltip-id="root-tooltip"
                data-tooltip-content="No vault assigned"
            >
                No vault
            </span>
        )
    }

    function openVaultManage(v: VaultItem) {
        const {CloseDialog} = useDialog.getState()
        const {auth} = useUser.getState()
        OpenDialog(
            <ManageVaultDialog
                vault={v}
                currentUserId={auth.userInfo?.id ?? ""}
                onClose={CloseDialog}
                onDeleted={CloseDialog}
            />
        )
    }

    return (
        <span
            className={cls.VaultChip}
            onClick={e => { e.stopPropagation(); openVaultManage(vault) }}
            data-tooltip-id="root-tooltip"
            data-tooltip-content={`Vault: ${vault.name}`}
        >
            <span className={cls.VaultBadge}>A</span>
            {vault.name}
        </span>
    )
}
