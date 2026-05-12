import {useState} from "react"
import cls from "@/pages/home/VaultCard.module.css"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"

interface Props {
    vault: VaultItem
}

export default function VaultCard({vault}: Props) {
    const [copied, setCopied] = useState(false)

    async function handleCopy() {
        await navigator.clipboard.writeText(vault.dbUrl ?? "")
        setCopied(true)
        setTimeout(() => setCopied(false), 1500)
    }

    return (
        <div className={cls.Card}>
            <span className={cls.Name}>{vault.name}</span>
            <div className={cls.Footer}>
                <span className={cls.Url}>{vault.dbUrl}</span>
                <button className={cls.CopyBtn} onClick={handleCopy} title="Copy connection string">
                    {copied ? "✓" : "copy"}
                </button>
            </div>
        </div>
    )
}
