import {Button} from "@vervstack/chures"

import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import cls from "@/pages/admin/components/SettingsTab/components/DefaultDocsVaultSection/components/PublishedVaultsGroup/PublishedVaultsGroup.module.css" // eslint-disable-line max-len
import PublishedVaultRow from "@/pages/admin/components/SettingsTab/components/DefaultDocsVaultSection/components/PublishedVaultRow/PublishedVaultRow.tsx" // eslint-disable-line max-len

interface Props {
    vaults: VaultItem[]
    defaultDocsVaultId: string
    onSelect: (vaultId: string) => void
    onClear: () => void
}

export default function PublishedVaultsGroup({vaults, defaultDocsVaultId, onSelect, onClear}: Props) {
    return (
        <div className={cls.PublishedVaultsGroupContainer}>
            <h3 className={cls.GroupTitle}>Published</h3>
            {vaults.length === 0 ? (
                <p className={cls.Empty}>No vaults are published yet — publish one below to set it as the default.</p>
            ) : (
                <div className={cls.RowList}>
                    {vaults.map(vault => (
                        <PublishedVaultRow
                            key={vault.id}
                            vault={vault}
                            isDefault={vault.id === defaultDocsVaultId}
                            onSelect={() => onSelect(vault.id ?? "")}
                        />
                    ))}
                </div>
            )}
            {defaultDocsVaultId && vaults.length > 0 && (
                <Button variant="ghost" onClick={onClear} className={cls.ClearButton}>
                    Clear default
                </Button>
            )}
        </div>
    )
}
