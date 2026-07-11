import cls from "@/widgets/EmailCard/EmailCard.module.css"
import {cn} from "@/app/utils/cn.ts"
import {ExternalConnectionInfo, ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"

export default function EmailCard({connections, loading, onClick}: {
    connections: ExternalConnectionInfo[]
    loading: boolean
    onClick?: () => void
}) {
    const isConnected = connections.length > 0
    const accountLabel = connections.length === 1
        ? (connections[0].generic?.fields?.username ?? connections[0].generic?.fields?.email)
        : connections.length > 1 ? `${connections.length} accounts` : undefined

    return (
        <div
            className={cls.CardContainer}
            onClick={!loading ? onClick : undefined}
            role="button"
            tabIndex={0}
            onKeyDown={e => { if (e.key === "Enter" || e.key === " ") onClick?.() }}
        >
            <div className={cls.CardHeader}>
                <div className={cls.CardIcon}>
                    <ProviderIcon provider={ExternalProvider.EXTERNAL_PROVIDER_EMAIL}/>
                </div>
                <div className={cls.CardTitles}>
                    <div className={cls.CardName}>Email</div>
                    {accountLabel && <div className={cls.CardAccount}>{accountLabel}</div>}
                </div>
            </div>
            <div className={cls.CardFooter}>
                <span className={cn(cls.StatusDot, isConnected ? cls.StatusDotConnected : cls.StatusDotDisconnected)}/>
                <span
                    className={cn(
                        cls.StatusLabel, isConnected ? cls.StatusLabelConnected : cls.StatusLabelDisconnected
                    )}
                >
                    {loading ? "…" : isConnected ? "Connected" : "Not connected"}
                </span>
            </div>
        </div>
    )
}
