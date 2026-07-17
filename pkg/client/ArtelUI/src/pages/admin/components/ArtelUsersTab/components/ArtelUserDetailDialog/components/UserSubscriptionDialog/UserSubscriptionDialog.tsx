import {Button, ModalClose} from "@vervstack/chures"

import cls
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/UserSubscriptionDialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import SubscriptionForm
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/components/SubscriptionForm/SubscriptionForm.tsx"
import {useUserSubscriptionDialog}
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/processes/useUserSubscriptionDialog.tsx"

interface UserSubscriptionDialogProps {
    userId: string
    onUpdated?: () => void
}

export default function UserSubscriptionDialog({userId, onUpdated}: UserSubscriptionDialogProps) {
    const {CloseDialog} = useDialog()
    const form = useUserSubscriptionDialog(userId, onUpdated)

    return (
        <div className={cls.UserSubscriptionDialogContainer} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>Manage subscription</h2>
                <ModalClose onClick={CloseDialog} className={cls.ModalClose} />
            </div>
            <div className={cls.ModalBody}>
                {form.loading ? (
                    <p className={cls.Empty}>Loading…</p>
                ) : (
                    <SubscriptionForm
                        plans={form.plans}
                        active={{value: form.active, onChange: form.setActive}}
                        plan={{value: form.planKey, onChange: form.setPlanKey}}
                        featureOverrides={{states: form.featureStates, onChange: form.handleFeatureStateChange}}
                        couchQuota={{
                            overrideOn: form.couchOverrideOn,
                            overrideMb: form.couchOverrideMb,
                            onToggle: form.setCouchOverrideOn,
                            onChangeValue: form.setCouchOverrideMb,
                        }}
                        s3Quota={{
                            overrideOn: form.s3OverrideOn,
                            overrideMb: form.s3OverrideMb,
                            onToggle: form.setS3OverrideOn,
                            onChangeValue: form.setS3OverrideMb,
                        }}
                    />
                )}
            </div>
            <div className={cls.ModalActions}>
                <Button
                    variant="danger"
                    onClick={form.handleCancelSubscription}
                    disabled={form.loading || form.saving || !form.active}
                >
                    Cancel subscription
                </Button>
                <Button variant="primary" onClick={form.handleSave} disabled={form.loading || form.saving}>
                    {form.saving ? "Saving…" : "Save"}
                </Button>
            </div>
        </div>
    )
}
