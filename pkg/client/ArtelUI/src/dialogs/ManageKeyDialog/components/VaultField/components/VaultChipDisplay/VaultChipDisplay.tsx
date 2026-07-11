import {useVaults} from "@/app/hooks/Vaults.ts"
import cls
    from "@/dialogs/ManageKeyDialog/components/VaultField/components/VaultChipDisplay/VaultChipDisplay.module.css"

interface VaultChipDisplayProps {
    vault: ReturnType<typeof useVaults>["vaults"][number] | undefined
}

export default function VaultChipDisplay({vault}: VaultChipDisplayProps) {
    return (
        <span
            className={vault ? cls.VaultChip : `${cls.Chip} ${cls.ChipMuted}`}
            data-tooltip-id="root-tooltip"
            data-tooltip-content={vault ? `Vault: ${vault.name}` : "No vault assigned"}
        >
            {vault && <span className={cls.VaultBadge}>A</span>}
            {vault ? vault.name : "No vault"}
        </span>
    )
}
