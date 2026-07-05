import {useState} from "react"
import {createPortal} from "react-dom"
import cls from "@/components/VaultCard/FastSetupDialog.module.css"
import {PromptsAPI, PromptId} from "@/app/api/artel/prompts.pb.ts"
import useUser from "@/hooks/user/User.ts"
import {CheckmarkIcon} from "@vervstack/chures"

interface Props {
    onClose: () => void
}

export default function FastSetupDialog({onClose}: Props) {
    const [view, setView] = useState<"options" | "text">("options")
    const [copied, setCopied] = useState(false)
    const [promptText, setPromptText] = useState("")
    const [loading, setLoading] = useState(false)

    async function handlePersonalVault() {
        setLoading(true)
        setView("text")
        try {
            const {auth} = useUser.getState()
            const resp = await PromptsAPI.ListPrompts(
                {ids: [PromptId.PROMPT_ID_PERSONAL_VAULT_SETUP], page: 1, pageSize: 1},
                auth.getInitReq()
            )
            setPromptText(resp.prompts?.[0]?.prompt ?? "")
        } finally {
            setLoading(false)
        }
    }

    async function handleCopy() {
        if (navigator.clipboard?.writeText) {
            await navigator.clipboard.writeText(promptText).catch(() => {})
        } else {
            const ta = document.createElement("textarea")
            ta.value = promptText
            document.body.appendChild(ta)
            ta.select()
            try { document.execCommand("copy") } catch {}
            document.body.removeChild(ta)
        }
        setCopied(true)
        setTimeout(() => setCopied(false), 1500)
    }

    return createPortal(
        <div className={cls.Overlay} onClick={onClose}>
            <div className={cls.Dialog} onClick={(e) => e.stopPropagation()}>
                {view === "text" && (
                    <button
                        type="button"
                        className={`${cls.CopyBtn} ${copied ? cls.CopyBtnCopied : ""}`}
                        onClick={handleCopy}
                        disabled={loading}
                    >
                        {copied ? (
                            <CheckmarkIcon size={16} strokeWidth={2} shrink={false}/>
                        ) : (
                            <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
                                <rect x="9" y="9" width="13" height="13" rx="2"/>
                                <path d="M5 15V5a2 2 0 0 1 2-2h10"/>
                            </svg>
                        )}
                    </button>
                )}
                {view === "options" ? (
                    <button type="button" className={cls.OptionBtn} onClick={handlePersonalVault}>
                        Personal Vault
                    </button>
                ) : (
                    <pre className={cls.PromptText}>{loading ? "Loading…" : promptText}</pre>
                )}
            </div>
        </div>,
        document.body
    )
}
