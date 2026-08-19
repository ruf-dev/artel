import {Loader} from "@vervstack/chures"

import cls from "@/pages/workbench/components/Chat/components/ChatStatusBanner/ChatStatusBanner.module.css"
import {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"

interface Props {
    status: ChatConnectionStatus
}

const LABELS: Partial<Record<ChatConnectionStatus, string>> = {
    connecting: "Connecting…",
    reconnecting: "Reconnecting…",
    closed: "Disconnected",
}

export default function ChatStatusBanner({status}: Props) {
    if (status === "open") return null

    return (
        <div className={cls.ChatStatusBannerContainer}>
            <Loader variant="arcs" size="sm" color="var(--coral)"/>
            <span>{LABELS[status]}</span>
        </div>
    )
}
