import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import cls from "@/pages/admin/components/SettingsTab/components/DefaultDocsVaultSection/components/PrivateVaultsGroup/PrivateVaultsGroup.module.css" // eslint-disable-line max-len
import PrivateVaultRow from "@/pages/admin/components/SettingsTab/components/DefaultDocsVaultSection/components/PrivateVaultRow/PrivateVaultRow.tsx" // eslint-disable-line max-len

interface Props {
    vaults: VaultItem[]
    onPublished: (vault: VaultItem) => void
}

export default function PrivateVaultsGroup({vaults, onPublished}: Props) {
    return (
        <div className={cls.PrivateVaultsGroupContainer}>
            <h3 className={cls.GroupTitle}>Private</h3>
            {vaults.length === 0 ? (
                <p className={cls.Empty}>All your vaults are already published.</p>
            ) : (
                <div className={cls.RowList}>
                    {vaults.map(vault => (
                        <PrivateVaultRow key={vault.id} vault={vault} onPublished={onPublished}/>
                    ))}
                </div>
            )}
        </div>
    )
}
