import {useState} from "react"
import {Button} from "@vervstack/chures"
import cls from "@/components/VaultCard/VaultCardConnBar.module.css"
import {cn} from "@/app/utils/cn"

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

interface RowProps {
    value: string
    label: string
    hidden?: boolean
}

function ConnBarRow({value, label, hidden = false}: RowProps) {
    const [copied, setCopied] = useState(false)
    const [revealed, setRevealed] = useState(false)

    async function handleCopy() {
        if (navigator.clipboard?.writeText) {
            await navigator.clipboard.writeText(value).catch(() => {})
        } else {
            const ta = document.createElement("textarea")
            ta.value = value
            document.body.appendChild(ta)
            ta.select()
            try { document.execCommand("copy") } catch {}
            document.body.removeChild(ta)
        }
        setCopied(true)
        setTimeout(() => setCopied(false), 1500)
    }

    const displayValue = hidden && !revealed ? "•".repeat(Math.min(value.length, 16)) : value

    return (
        <div className={cls.ConnBarRowWrapper}>
            <span className={cls.ConnBarLabel}>{label}</span>
            <code className={cls.ConnString} title={revealed || !hidden ? value : undefined}>{displayValue}</code>
            {hidden && (
                <Button
                    className={cls.CopyBtn}
                    onClick={(e) => {
                        e.stopPropagation()
                        setRevealed(r => !r)
                    }}
                >
                    {revealed ? (
                        <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94"/>
                            <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19"/>
                            <line x1="1" y1="1" x2="23" y2="23"/>
                        </svg>
                    ) : (
                        <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                            <circle cx="12" cy="12" r="3"/>
                        </svg>
                    )}
                </Button>
            )}
            <Button
                className={cn(cls.CopyBtn, copied && cls.CopyBtnCopied)}
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
            </Button>
        </div>
    )
}
