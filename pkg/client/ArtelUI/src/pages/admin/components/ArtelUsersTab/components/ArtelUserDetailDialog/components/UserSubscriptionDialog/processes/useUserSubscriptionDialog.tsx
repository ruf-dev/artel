import {useState, useEffect, useCallback} from "react"
import {ConfirmDialog} from "@vervstack/chures"

import {AdminSubscriptionsAPI, SubscriptionPlanEntry} from "@/app/api/artel/admin_subscriptions.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import {FeatureOverrideState}
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/components/FeatureOverrideRow/FeatureOverrideRow.tsx"
import {FEATURES, bytesToMb, mbToBytes, initFeatureStates}
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/processes/subscriptionOverrides.ts"

export function useUserSubscriptionDialog(userId: string, onUpdated?: () => void) {
    const {auth} = useUser()
    const {OpenDialog, CloseDialog} = useDialog()
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

    function buildFeatureOverrides(): Record<string, boolean> {
        const featureOverrides: Record<string, boolean> = {}
        for (const feature of FEATURES) {
            const state = featureStates[feature.key]
            if (state === "on") featureOverrides[feature.key] = true
            if (state === "off") featureOverrides[feature.key] = false
        }
        return featureOverrides
    }

    function updateSubscription(nextActive: boolean) {
        return AdminSubscriptionsAPI.UpdateUserSubscription({
            userId,
            active: nextActive,
            planKey,
            featureOverrides: buildFeatureOverrides(),
            couchQuotaOverrideBytes: couchOverrideOn ? mbToBytes(couchOverrideMb) : undefined,
            s3QuotaOverrideBytes: s3OverrideOn ? mbToBytes(s3OverrideMb) : undefined,
        }, auth.getInitReq())
    }

    function handleSave() {
        setSaving(true)
        updateSubscription(active)
            .then(() => {
                onUpdated?.()
                CloseDialog()
            })
            .catch(err => bakeError("Failed to update subscription", err))
            .finally(() => setSaving(false))
    }

    function handleCancelSubscription() {
        CloseDialog()
        OpenDialog(
            <ConfirmDialog
                title="Cancel subscription"
                message="This immediately disables the subscription for this user."
                confirmLabel="Cancel subscription"
                danger
                onClose={CloseDialog}
                onConfirm={() => updateSubscription(false)
                    .then(() => onUpdated?.())
                    .catch(err => bakeError("Failed to cancel subscription", err))}
            />
        )
    }

    function handleFeatureStateChange(key: string, state: FeatureOverrideState) {
        setFeatureStates(prev => ({...prev, [key]: state}))
    }

    return {
        loading,
        saving,
        plans,
        active, setActive,
        planKey, setPlanKey,
        featureStates,
        couchOverrideOn, setCouchOverrideOn,
        couchOverrideMb, setCouchOverrideMb,
        s3OverrideOn, setS3OverrideOn,
        s3OverrideMb, setS3OverrideMb,
        handleSave,
        handleCancelSubscription,
        handleFeatureStateChange,
    }
}
