import cls from
    "@/dialogs/ManageOpenRouterDialog/components/ViewOpenRouterUsageDialog/components/UsageStat/UsageStat.module.css"

interface UsageStatProps {
    label: string
    value: string
}

export default function UsageStat({label, value}: UsageStatProps) {
    return (
        <div className={cls.UsageStatContainer}>
            <span className={cls.UsageStatLabel}>{label}</span>
            <span className={cls.UsageStatValue}>{value}</span>
        </div>
    )
}
