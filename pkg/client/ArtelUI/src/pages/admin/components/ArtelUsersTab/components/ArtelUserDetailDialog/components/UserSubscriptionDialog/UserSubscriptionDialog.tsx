import {useState, useEffect, useCallback} from "react"
import {Button, ModalClose} from "@vervstack/chures"

import cls
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/UserSubscriptionDialog.module.css"
import {AdminSubscriptionsAPI, SubscriptionPlanEntry} from "@/app/api/artel/admin_subscriptions.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import SubscriptionForm
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/components/SubscriptionForm/SubscriptionForm.tsx"
import {FeatureOverrideState}
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/components/FeatureOverrideRow/FeatureOverrideRow.tsx"
import {FEATURES, bytesToMb, mbToBytes, initFeatureStates}
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/processes/subscriptionOverrides.ts"

interface UserSubscriptionDialogProps {
    userId: string
    onUpdated?: () => void
}

export default function UserSubscriptionDialog({userId, onUpdated}: UserSubscriptionDialogProps) {
    const {auth} = useUser()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

    const [loading, setLoading] = useState(true)
    const [saving, setSaving] = useState(false)
    const [plans, setPlans] = useState<SubscriptionPlanEntry[]>([])

    const [active, setActive] = useState(false)
    const [planKey, setPlanKey] = useState("")
    const [featureStates, setFeatureStates] = useState<Record<string, FeatureOverrideState>>({})
    const [couchOverrideOn, setCouchOverrideOn] = useState(false)
    const [couchOverrideMb, setCouchOverrideMb] = useState("0")
    const [s3OverrideOn, setS3OverrideOn] = useState(false)
    const [s3OverrideMb, setS3OverrideMb] = useState("0")

    const load = useCallback(() => {
        setLoading(true)
        Promise.all([
            AdminSubscriptionsAPI.ListSubscriptionPlans({}, auth.getInitReq()),
            AdminSubscriptionsAPI.GetUserSubscription({userId}, auth.getInitReq()),
        ])
            .then(([plansRes, subRes]) => {
                setPlans(plansRes.plans ?? [])
                setActive(subRes.active ?? false)
                setPlanKey(subRes.planKey ?? "")
                setFeatureStates(initFeatureStates(subRes.featureOverrides))
                setCouchOverrideOn(subRes.couchQuotaOverrideBytes !== undefined)
                setCouchOverrideMb(String(bytesToMb(subRes.couchQuotaOverrideBytes)))
                setS3OverrideOn(subRes.s3QuotaOverrideBytes !== undefined)
                setS3OverrideMb(String(bytesToMb(subRes.s3QuotaOverrideBytes)))
            })
            .catch(err => bakeError("Failed to load subscription", err))
            .finally(() => setLoading(false))
    }, [auth, userId, bakeError])

    useEffect(() => { load() }, [load])

    function handleSave() {
        const featureOverrides: Record<string, boolean> = {}
        for (const feature of FEATURES) {
            const state = featureStates[feature.key]
            if (state === "on") featureOverrides[feature.key] = true
            if (state === "off") featureOverrides[feature.key] = false
        }

        setSaving(true)
        AdminSubscriptionsAPI.UpdateUserSubscription({
            userId,
            active,
            planKey,
            featureOverrides,
            couchQuotaOverrideBytes: couchOverrideOn ? mbToBytes(couchOverrideMb) : undefined,
            s3QuotaOverrideBytes: s3OverrideOn ? mbToBytes(s3OverrideMb) : undefined,
        }, auth.getInitReq())
            .then(() => {
                onUpdated?.()
                CloseDialog()
            })
            .catch(err => bakeError("Failed to update subscription", err))
            .finally(() => setSaving(false))
    }

    function handleFeatureStateChange(key: string, state: FeatureOverrideState) {
        setFeatureStates(prev => ({...prev, [key]: state}))
    }

    return (
        <div className={cls.UserSubscriptionDialogContainer} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>Manage subscription</h2>
                <ModalClose onClick={CloseDialog} className={cls.ModalClose} />
            </div>
            {loading ? (
                <p className={cls.Empty}>Loading…</p>
            ) : (
                <SubscriptionForm
                    plans={plans}
                    active={{value: active, onChange: setActive}}
                    plan={{value: planKey, onChange: setPlanKey}}
                    featureOverrides={{states: featureStates, onChange: handleFeatureStateChange}}
                    couchQuota={{
                        overrideOn: couchOverrideOn,
                        overrideMb: couchOverrideMb,
                        onToggle: setCouchOverrideOn,
                        onChangeValue: setCouchOverrideMb,
                    }}
                    s3Quota={{
                        overrideOn: s3OverrideOn,
                        overrideMb: s3OverrideMb,
                        onToggle: setS3OverrideOn,
                        onChangeValue: setS3OverrideMb,
                    }}
                />
            )}
            <div className={cls.ModalActions}>
                <Button variant="secondary" onClick={CloseDialog} disabled={saving}>Cancel</Button>
                <Button variant="primary" onClick={handleSave} disabled={loading || saving}>
                    {saving ? "Saving…" : "Save"}
                </Button>
            </div>
        </div>
    )
}
