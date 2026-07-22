import {ReactNode} from "react"

import cls from "@/widgets/ProviderCard/ProviderCard.module.css"
import {cn} from "@/app/utils/cn.ts"
import {ExternalConnectionInfo, ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"

export default function ProviderCard({provider, name, icon, connections, loading, onClick}: {
    provider: ExternalProvider
    name: string
    icon?: ReactNode
    connections: ExternalConnectionInfo[]
    loading: boolean
    onClick?: () => void
}) {
    const isConnected = connections.length > 0
    const single = connections.length === 1 ? connections[0] : undefined
    const accountLabel = connections.length > 1
        ? `${connections.length} accounts`
        : single?.google?.email ?? single?.generic?.fields?.["username"] ?? single?.generic?.fields?.["email"]

    return (
        <div
            className={cls.CardContainer}
            onClick={!loading ? onClick : undefined}
            role="button"
            tabIndex={0}
            onKeyDown={e => { if (e.key === "Enter" || e.key === " ") onClick?.() }}
        >
            <div className={cls.CardHeader}>
                <div className={cls.CardIcon}>{icon ?? <ProviderIcon provider={provider}/>}</div>
                <div className={cls.CardTitles}>
                    <div className={cls.CardName}>{name}</div>
                    {accountLabel && <div className={cls.CardAccount}>{accountLabel}</div>}
                </div>
            </div>
            <div className={cls.CardFooter}>
                <span className={cn(cls.StatusDot, isConnected ? cls.StatusDotConnected : cls.StatusDotDisconnected)}/>
                <span
                    className={cn(
                        cls.StatusLabel,
                        isConnected ? cls.StatusLabelConnected : cls.StatusLabelDisconnected,
                    )}
                >
                    {loading ? "…" : isConnected ? "Connected" : "Not connected"}
                </span>
            </div>
        </div>
    )
}
