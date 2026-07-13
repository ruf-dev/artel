import {useRef, useState} from "react"

import {ExternalConnectionsAPI} from "@/app/api/artel/external_connections.pb.ts"
import useUser from "@/hooks/user/User.ts"

export function useMailServerSuggestion() {
    const [imapHost, setImapHost] = useState("")
    const [imapPort, setImapPort] = useState("")
    const [smtpHost, setSmtpHost] = useState("")
    const [smtpPort, setSmtpPort] = useState("")
    const [appPasswordUrl, setAppPasswordUrl] = useState("")
    const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)

    const {auth} = useUser()

    function applySuggestion(email: string) {
        const domain = email.split("@")[1] ?? ""
        if (!domain) return

        ExternalConnectionsAPI.ListMailServerSuggestions({domain}, auth.getInitReq())
            .then(resp => {
                if (resp.suggestions?.length === 1) {
                    const s = resp.suggestions[0]
                    setImapHost(s.imap ?? "")
                    setImapPort(s.imapPort ? String(s.imapPort) : "")
                    setSmtpHost(s.smtp ?? "")
                    setSmtpPort(s.smtpPort ? String(s.smtpPort) : "")
                    setAppPasswordUrl(s.appPasswordUrl ?? "")
                }
            })
            .catch(() => {})
    }

    function notifyEmailChanged(email: string) {
        if (debounceRef.current) clearTimeout(debounceRef.current)
        debounceRef.current = setTimeout(() => applySuggestion(email), 300)
    }

    function seedFromEmail(email: string) {
        applySuggestion(email)
    }

    return {
        imapHost, setImapHost,
        imapPort, setImapPort,
        smtpHost, setSmtpHost,
        smtpPort, setSmtpPort,
        appPasswordUrl,
        notifyEmailChanged,
        seedFromEmail,
    }
}
