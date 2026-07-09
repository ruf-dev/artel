import {useState} from "react"

import cls from "@/widgets/VaultCard/VaultCard.module.css"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import VaultCardFront from "@/components/VaultCard/VaultCardFront.tsx";
import VaultCardBack from "@/components/VaultCard/VaultCardBack.tsx";
import {cn} from "@/app/utils/cn";

interface Props {
    vault: VaultItem
    onEdit?: (id: string) => void
}

export default function VaultCard({vault, onEdit}: Props) {
    const [isFlipped, setIsFlipped] = useState(false)

    return (
        <article className={cls.VaultCardContainer} onClick={() => setIsFlipped(!isFlipped)}>
            <div className={cn(cls.VaultCardContentWrapper, isFlipped && cls.Flipped)}>
                <VaultCardFront vault={vault} onEdit={onEdit}/>
                <VaultCardBack/>
            </div>
        </article>
    )
}
