import {Input, Toggle} from "@vervstack/chures"

import cls
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/components/QuotaOverrideRow/QuotaOverrideRow.module.css"

interface QuotaOverrideRowProps {
    label: string
    overrideEnabled: boolean
    overrideValueMb: string
    planDefaultMb: number
    onToggleOverride: (enabled: boolean) => void
    onChangeValue: (mb: string) => void
}

export default function QuotaOverrideRow({
    label, overrideEnabled, overrideValueMb, planDefaultMb, onToggleOverride, onChangeValue,
}: QuotaOverrideRowProps) {
    return (
        <div className={cls.QuotaOverrideRowContainer}>
            <div className={cls.QuotaOverrideHead}>
                <span className={cls.QuotaLabel}>{label}</span>
                <Toggle checked={overrideEnabled} onChange={onToggleOverride} label="Override" labelPosition="left" />
            </div>
            {overrideEnabled ? (
                <Input value={overrideValueMb} setValue={onChangeValue} type="number" label="MB" />
            ) : (
                <span className={cls.QuotaDefault}>Plan default: {planDefaultMb} MB</span>
            )}
        </div>
    )
}
