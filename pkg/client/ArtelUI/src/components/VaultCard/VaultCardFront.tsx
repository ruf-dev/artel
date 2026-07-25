import cls from "@/components/VaultCard/VaultCardFront.module.css"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import VaultCardHeader from "@/components/VaultCard/VaultCardHeader.tsx";
import VaultCardStatus from "@/components/VaultCard/VaultCardStatus.tsx";
import VaultCardConnBar from "@/components/VaultCard/VaultCardConnBar.tsx";

interface Props {
    vault: VaultItem
    onEdit?: (id: string) => void
    onWorkbench?: (id: string) => void
    isWorkbenchAvailable?: boolean
}

export default function VaultCardFront({vault, onEdit, onWorkbench, isWorkbenchAvailable}: Props) {
    return (
        <div className={cls.VaultCardFrontContainer}>
            <VaultCardHeader
                vault={vault}
                onEdit={onEdit}
                onWorkbench={onWorkbench}
                isWorkbenchAvailable={isWorkbenchAvailable}
            />
            <VaultCardStatus/>
            <VaultCardConnBar dbUrl={vault.dbUrl ?? ""} passphrase={vault.livesyncPassphrase ?? ""}/>
        </div>
    )
}
