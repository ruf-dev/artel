import {Dropdown, Toggle} from "@vervstack/chures"

import cls
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/components/SubscriptionForm/SubscriptionForm.module.css"
import {SubscriptionPlanEntry} from "@/app/api/artel/admin_subscriptions.pb.ts"
import FeatureOverrideRow, {FeatureOverrideState}
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/components/FeatureOverrideRow/FeatureOverrideRow.tsx"
import QuotaOverrideRow
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/components/QuotaOverrideRow/QuotaOverrideRow.tsx"
import {FEATURES, bytesToMb}
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/processes/subscriptionOverrides.ts"

interface FieldControl<T> {
    value: T
    onChange: (value: T) => void
}

interface QuotaControl {
    overrideOn: boolean
    overrideMb: string
    onToggle: (enabled: boolean) => void
    onChangeValue: (mb: string) => void
}

interface SubscriptionFormProps {
    plans: SubscriptionPlanEntry[]
    active: FieldControl<boolean>
    plan: FieldControl<string>
    featureOverrides: {
        states: Record<string, FeatureOverrideState>
        onChange: (key: string, state: FeatureOverrideState) => void
    }
    couchQuota: QuotaControl
    s3Quota: QuotaControl
}

export default function SubscriptionForm({
    plans, active, plan, featureOverrides, couchQuota, s3Quota,
}: SubscriptionFormProps) {
    const selectedPlan = plans.find(p => p.planKey === plan.value)
    const planOptions = plans.map(p => p.planKey ?? "")

    return (
        <div className={cls.SubscriptionFormContainer}>
            <div className={cls.Field}>
                <Toggle checked={active.value} onChange={active.onChange} label="Active" labelPosition="left" />
            </div>
            <div className={cls.Field}>
                <span className={cls.FieldLabel}>Plan</span>
                <Dropdown
                    options={planOptions}
                    value={plan.value ? [plan.value] : []}
                    onChange={v => plan.onChange(v[0] ?? "")}
                />
            </div>
            <div className={cls.Field}>
                <span className={cls.FieldLabel}>Feature overrides</span>
                {FEATURES.map(feature => (
                    <FeatureOverrideRow
                        key={feature.key}
                        label={feature.label}
                        state={featureOverrides.states[feature.key] ?? "inherit"}
                        onChange={state => featureOverrides.onChange(feature.key, state)}
                    />
                ))}
            </div>
            <div className={cls.Field}>
                <span className={cls.FieldLabel}>Quota overrides</span>
                <QuotaOverrideRow
                    label="CouchDB"
                    overrideEnabled={couchQuota.overrideOn}
                    overrideValueMb={couchQuota.overrideMb}
                    planDefaultMb={bytesToMb(selectedPlan?.couchQuotaBytes)}
                    onToggleOverride={couchQuota.onToggle}
                    onChangeValue={couchQuota.onChangeValue}
                />
                <QuotaOverrideRow
                    label="S3"
                    overrideEnabled={s3Quota.overrideOn}
                    overrideValueMb={s3Quota.overrideMb}
                    planDefaultMb={bytesToMb(selectedPlan?.s3QuotaBytes)}
                    onToggleOverride={s3Quota.onToggle}
                    onChangeValue={s3Quota.onChangeValue}
                />
            </div>
        </div>
    )
}
