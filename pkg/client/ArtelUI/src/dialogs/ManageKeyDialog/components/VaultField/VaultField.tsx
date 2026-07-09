import {Button} from "@vervstack/chures"

import {useVaults} from "@/app/hooks/Vaults.ts"
import VaultChipDisplay from "@/dialogs/ManageKeyDialog/components/VaultField/components/VaultChipDisplay/VaultChipDisplay.tsx"
import cls from "@/dialogs/ManageKeyDialog/components/VaultField/VaultField.module.css"

interface VaultFieldProps {
    selectedVaultId: string
    onChangeVault: () => void
}

export default function VaultField({selectedVaultId, onChangeVault}: VaultFieldProps) {
    const {vaults} = useVaults()
    const vault = vaults.find(v => v.id === selectedVaultId)
    return (
        <div className={cls.Field}>
            <span className={cls.FieldLabel}>Vault</span>
            <div className={cls.ConnectorRowWrapper}>
                <VaultChipDisplay vault={vault}/>
                <Button variant="ghost" onClick={onChangeVault}>Change</Button>
            </div>
        </div>
    )
}
