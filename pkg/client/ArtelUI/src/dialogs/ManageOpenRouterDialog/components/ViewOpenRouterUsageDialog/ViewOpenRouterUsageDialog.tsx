import {useEffect, useState} from "react"
import {LoadingWrapper} from "@vervstack/chures"

import {ExternalProvider, OpenRouterStatistics} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import LlmKeyDialogHead from "@/components/LlmKeyDialogHead/LlmKeyDialogHead.tsx"
import UsageStat from
    "@/dialogs/ManageOpenRouterDialog/components/ViewOpenRouterUsageDialog/components/UsageStat/UsageStat.tsx"
import cls from
    "@/dialogs/ManageOpenRouterDialog/components/ViewOpenRouterUsageDialog/ViewOpenRouterUsageDialog.module.css"

function formatUsd(amount: number): string {
    return `$${amount.toFixed(2)}`
}

export default function ViewOpenRouterUsageDialog() {
    const {CloseDialog} = useDialog()
    const {getProviderStatistics} = useExternalConnections()
    const bakeError = useBakeError()

    const [loading, setLoading] = useState(true)
    const [stats, setStats] = useState<OpenRouterStatistics | null>(null)

    useEffect(() => {
        setLoading(true)
        getProviderStatistics(ExternalProvider.EXTERNAL_PROVIDER_OPENROUTER)
            .then(res => setStats(res.openrouter ?? null))
            .catch(err => bakeError("Failed to load usage", err))
            .finally(() => setLoading(false))
    }, [getProviderStatistics, bakeError])

    return (
        <div
            className={cls.ViewOpenRouterUsageDialogContainer}
            onClick={e => e.stopPropagation()}
            role="dialog"
            aria-modal="true"
        >
            <LlmKeyDialogHead
                provider={ExternalProvider.EXTERNAL_PROVIDER_OPENROUTER}
                title="OpenRouter usage"
                onClose={CloseDialog}
            />
            <LoadingWrapper isLoading={loading}>
                {stats && (
                    <>
                        <label className={cls.Field}>
                            <span className={cls.FieldLabel}>Credit limit</span>
                            <div className={cls.FieldValue}>
                                {stats.limitUnlimited ? "Unlimited" : formatUsd(stats.limit ?? 0)}
                            </div>
                        </label>
                        {!stats.limitUnlimited && (
                            <label className={cls.Field}>
                                <span className={cls.FieldLabel}>Remaining</span>
                                <div className={cls.FieldValue}>{formatUsd(stats.limitRemaining ?? 0)}</div>
                            </label>
                        )}
                        <div className={cls.UsageGrid}>
                            <UsageStat label="Total" value={formatUsd(stats.usageTotal ?? 0)}/>
                            <UsageStat label="Today" value={formatUsd(stats.usageDaily ?? 0)}/>
                            <UsageStat label="This week" value={formatUsd(stats.usageWeekly ?? 0)}/>
                            <UsageStat label="This month" value={formatUsd(stats.usageMonthly ?? 0)}/>
                        </div>
                        {stats.isFreeTier && (
                            <p className={cls.FreeTierNotice}>This key is on OpenRouter's free tier.</p>
                        )}
                    </>
                )}
            </LoadingWrapper>
        </div>
    )
}
