import {Button} from "@vervstack/chures"

import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {cn} from "@/app/utils/cn.ts"
import cls from "@/pages/admin/components/SettingsTab/components/DefaultDocsVaultSection/components/PublishedVaultRow/PublishedVaultRow.module.css" // eslint-disable-line max-len

interface Props {
    vault: VaultItem
    isDefault: boolean
    onSelect: () => void
}

export default function PublishedVaultRow({vault, isDefault, onSelect}: Props) {
    return (
        <div className={cn(cls.PublishedVaultRowContainer, isDefault && cls.PublishedVaultRowDefault)}>
            <div className={cls.Info}>
                <span className={cls.Name}>{vault.name}</span>
                <span className={cls.Slug}>/docs/{vault.slug}</span>
            </div>
            <Button
                variant={isDefault ? "primary" : "secondary"}
                onClick={onSelect}
                disabled={isDefault}
            >
                {isDefault ? "Current default" : "Set as default"}
            </Button>
        </div>
    )
}
