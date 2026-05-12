import {useState} from "react"
import cls from "@/components/VaultCard/VaultCardConnBar.module.css"

interface Props {
    dbUrl: string
}

export default function VaultCardConnBar({dbUrl}: Props) {
    const [copied, setCopied] = useState(false)

    async function handleCopy() {
        if (navigator.clipboard?.writeText) {
            await navigator.clipboard.writeText(dbUrl).catch(() => {})
        } else {
            const ta = document.createElement("textarea")
            ta.value = dbUrl
            document.body.appendChild(ta)
            ta.select()
            try { document.execCommand("copy") } catch {}
            document.body.removeChild(ta)
        }
        setCopied(true)
        setTimeout(() => setCopied(false), 1500)
    }

    return (
        <div className={cls.VaultCardConnBarContainer}>
            <code className={cls.ConnString} title={dbUrl}>{dbUrl}</code>
            <button
                className={`${cls.CopyBtn} ${copied ? cls.CopyBtnCopied : ""}`}
                type="button"
                onClick={(e) => {
                    e.stopPropagation()
                    handleCopy()
                }}
            >
                <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
                    <rect x="9" y="9" width="13" height="13" rx="2"/>
                    <path d="M5 15V5a2 2 0 0 1 2-2h10"/>
                </svg>
                <span>{copied ? "Copied" : "Copy"}</span>
            </button>
        </div>
    )
}
