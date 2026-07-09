import {useVaults} from "@/app/hooks/Vaults.ts"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import cls from "@/dialogs/ManageKeyDialog/components/VaultOptionList/VaultOptionList.module.css"

interface VaultOptionListProps {
    selectedVaultId: string
    onSelect: (id: string) => void
}

export default function VaultOptionList({selectedVaultId, onSelect}: VaultOptionListProps) {
    const {vaults} = useVaults()
    return (
        <div className={cls.OptionList}>
            {vaults.map(v => (
                <SelectOption
                    key={v.id}
                    label={v.name ?? ""}
                    selected={selectedVaultId === v.id}
                    onSelect={() => onSelect(v.id ?? "")}
                />
            ))}
            {vaults.length === 0 && <p className={cls.Empty}>No vaults found.</p>}
        </div>
    )
}
