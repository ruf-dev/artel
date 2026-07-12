import {createPortal} from "react-dom"
import {ModalClose} from "@vervstack/chures"

import cls from "@/dialogs/ManageVaultDialog/components/ExpertSettingsDrawer/ExpertSettingsDrawer.module.css"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {cn} from "@/app/utils/cn.ts"
import VaultSettingsSection from "@/dialogs/ManageVaultDialog/widgets/VaultSettingsSection/VaultSettingsSection.tsx"

interface Props {
    open: boolean
    vault: VaultItem
    onChanged: (patch: Partial<VaultItem>) => void
    onClose: () => void
}

export default function ExpertSettingsDrawer({open, vault, onChanged, onClose}: Props) {
    return createPortal(
        <>
            <div className={cn(cls.Backdrop, open && cls.BackdropOpen)} onClick={onClose}/>
            <div className={cn(cls.DrawerPanel, open && cls.DrawerPanelOpen)}>
                <div className={cls.DrawerHeader}>
                    <span className={cls.DrawerTitle}>Expert settings</span>
                    <ModalClose onClick={onClose} className={cls.DrawerClose}/>
                </div>
                <div className={cls.DrawerBody}>
                    <VaultSettingsSection vault={vault} onChanged={onChanged}/>
                </div>
            </div>
        </>,
        document.body
    )
}
