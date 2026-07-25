import cls from "@/widgets/VaultCard/VaultCard.module.css"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import VaultCardFront from "@/components/VaultCard/VaultCardFront.tsx";

interface Props {
    vault: VaultItem
    onEdit?: (id: string) => void
    onWorkbench?: (id: string) => void
    isWorkbenchAvailable?: boolean
}

export default function VaultCard({vault, onEdit, onWorkbench, isWorkbenchAvailable}: Props) {
    return (
        <article className={cls.VaultCardContainer}>
            <div className={cls.VaultCardContentWrapper}>
                <VaultCardFront
                    vault={vault}
                    onEdit={onEdit}
                    onWorkbench={onWorkbench}
                    isWorkbenchAvailable={isWorkbenchAvailable}
                />
            </div>
        </article>
    )
}
