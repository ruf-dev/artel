import cls from "@/components/VaultCard/VaultCardHeader.module.css"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"

interface Props {
    vault: VaultItem
    onEdit?: (id: string) => void
}

export default function VaultCardHeader({vault, onEdit}: Props) {
    return (
        <div className={cls.VaultCardHeaderContainer}>
            <h3 className={cls.Name}>{vault.name}</h3>
            {onEdit && (
                <button
                    className={cls.IconBtn}
                    type="button"
                    onClick={(e) => {
                        e.stopPropagation()
                        onEdit(vault.id ?? "")
                    }}
                    title="Edit vault"
                    aria-label="Edit vault"
                >
                    <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M12 20h9"/><path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4z"/>
                    </svg>
                </button>
            )}
        </div>
    )
}
