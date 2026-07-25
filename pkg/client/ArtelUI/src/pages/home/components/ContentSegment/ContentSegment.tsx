import {useVaults} from "@/app/hooks/Vaults.ts"
import VaultCard from "@/widgets/VaultCard/VaultCard.tsx"
import cls from "@/pages/home/components/ContentSegment/ContentSegment.module.css"

interface ContentSegmentProps {
    onEditClick: (id: string) => void
    onWorkbenchClick: (id: string) => void
    isWorkbenchAvailable: boolean
}

export default function ContentSegment({onEditClick, onWorkbenchClick, isWorkbenchAvailable}: ContentSegmentProps) {
    const {isLoading, vaults} = useVaults()

    let loadingState = null

    if (isLoading) {
        loadingState = (
            <p className={cls.Empty}>Loading…</p>
        )
    }

    return (
        <div className={cls.ContentSegmentContainer}>
            {loadingState}
            {
                !isLoading && <div className={cls.Grid}>
                    {vaults.map(v => (
                        <VaultCard
                            key={v.id}
                            vault={v}
                            onEdit={onEditClick}
                            onWorkbench={onWorkbenchClick}
                            isWorkbenchAvailable={isWorkbenchAvailable}
                        />
                    ))}
                </div>
            }
        </div>
    )
}
