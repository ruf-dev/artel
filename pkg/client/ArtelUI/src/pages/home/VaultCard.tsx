import {useState} from "react"
import cls from "@/pages/home/VaultCard.module.css"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"

interface Props {
    vault: VaultItem
    onEdit?: (id: string) => void
}

export default function VaultCard({vault, onEdit}: Props) {
    const [copied, setCopied] = useState(false)

    async function handleCopy() {
        const text = vault.dbUrl ?? ""
        if (navigator.clipboard?.writeText) {
            await navigator.clipboard.writeText(text).catch(() => {})
        } else {
            const ta = document.createElement("textarea")
            ta.value = text
            document.body.appendChild(ta)
            ta.select()
            try { document.execCommand("copy") } catch {}
            document.body.removeChild(ta)
        }
        setCopied(true)
        setTimeout(() => setCopied(false), 1500)
    }

    return (
        <article className={cls.Card}>
            <div className={cls.Head}>
                <h3 className={cls.Name}>{vault.name}</h3>
                {onEdit && (
                    <button
                        className={cls.IconBtn}
                        type="button"
                        onClick={() => onEdit(vault.id ?? "")}
                        title="Edit vault"
                        aria-label="Edit vault"
                    >
                        <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M12 20h9"/><path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4z"/>
                        </svg>
                    </button>
                )}
            </div>

            <div className={cls.Meta}>
                <span className={cls.Chip}>
                    <span className={cls.StatusDot}/>
                    <span>Online</span>
                </span>
            </div>

            <div className={cls.ConnBar}>
                <code className={cls.ConnString} title={vault.dbUrl}>{vault.dbUrl}</code>
                <button
                    className={`${cls.CopyBtn} ${copied ? cls.CopyBtnCopied : ""}`}
                    type="button"
                    onClick={handleCopy}
                >
                    <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
                        <rect x="9" y="9" width="13" height="13" rx="2"/>
                        <path d="M5 15V5a2 2 0 0 1 2-2h10"/>
                    </svg>
                    <span>{copied ? "Copied" : "Copy"}</span>
                </button>
            </div>
        </article>
    )
}
