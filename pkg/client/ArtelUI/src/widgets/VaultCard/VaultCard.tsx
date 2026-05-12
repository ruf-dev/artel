import {useState} from "react"
import cls from "@/widgets/VaultCard/VaultCard.module.css"

import {VaultItem} from "@/app/api/artel/vaults.pb.ts"

import VaultCardHeader from "@/components/VaultCard/VaultCardHeader.tsx";
import VaultCardStatus from "@/components/VaultCard/VaultCardStatus.tsx";
import VaultCardConnBar from "@/components/VaultCard/VaultCardConnBar.tsx";
import VaultCardQR from "@/components/VaultCard/VaultCardQR.tsx";
import VaultCardLinks from "@/components/VaultCard/VaultCardLinks.tsx";

interface Props {
    vault: VaultItem
    onEdit?: (id: string) => void
}

export default function VaultCard({vault, onEdit}: Props) {
    const [isFlipped, setIsFlipped] = useState(false)

    return (
        <article className={cls.VaultCardContainer} onClick={() => setIsFlipped(!isFlipped)}>
            <div className={`${cls.VaultCardContentWrapper} ${isFlipped ? cls.Flipped : ""}`}>
                <VaultCardFront vault={vault} onEdit={onEdit}/>
                <VaultCardBack/>
            </div>
        </article>
    )
}

function VaultCardFront({vault, onEdit}: Props) {
    return (
        <div className={cls.VaultCardFrontContainer}>
            <VaultCardHeader vault={vault} onEdit={onEdit}/>
            <VaultCardStatus/>
            <VaultCardConnBar dbUrl={vault.dbUrl ?? ""}/>
        </div>
    )
}


function VaultCardBack() {
    return (
        <div className={cls.VaultCardBackContainer}>
            <VaultCardQR/>
            <VaultCardLinks/>
        </div>
    )
}
