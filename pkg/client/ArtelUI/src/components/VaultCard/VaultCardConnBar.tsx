import cls from "@/components/VaultCard/VaultCardConnBar.module.css"
import ConnBarRow from "@/components/VaultCard/components/ConnBarRow"

interface Props {
    dbUrl: string
    passphrase: string
}

export default function VaultCardConnBar({dbUrl, passphrase}: Props) {
    return (
        <div className={cls.VaultCardConnBarContainer}>
            <ConnBarRow value={dbUrl} label="uri"/>
            <ConnBarRow value={passphrase} label="key" hidden/>
        </div>
    )
}
