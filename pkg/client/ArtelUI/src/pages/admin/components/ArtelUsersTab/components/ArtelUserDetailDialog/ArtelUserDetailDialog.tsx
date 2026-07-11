import {useState, useEffect} from "react"
import {Button, ModalClose} from "@vervstack/chures"

import cls
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/ArtelUserDetailDialog.module.css"
import {AdminUsersAPI, ArtelUserDetails} from "@/app/api/artel/admin_users.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import UserSessionsDialog
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSessionsDialog/UserSessionsDialog.tsx"

interface ArtelUserDetailDialogProps {
    userId: string
}

export default function ArtelUserDetailDialog({userId}: ArtelUserDetailDialogProps) {
    const {auth} = useUser()
    const {CloseDialog, OpenDialog} = useDialog()
    const bakeError = useBakeError()
    const [details, setDetails] = useState<ArtelUserDetails | null>(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        async function load() {
            try {
                const res = await AdminUsersAPI.GetArtelUser({userId}, auth.getInitReq())
                setDetails(res.user ?? null)
            } catch (err) {
                bakeError("Failed to load user details", err)
            } finally {
                setLoading(false)
            }
        }
        void load()
    }, [auth, userId])

    function openSessions() {
        OpenDialog(<UserSessionsDialog userId={userId} />)
    }

    return (
        <div className={cls.ArtelUserDetailDialogContainer} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>User details</h2>
                <ModalClose onClick={CloseDialog} className={cls.ModalClose} />
            </div>
            {loading ? (
                <p className={cls.Empty}>Loading…</p>
            ) : details ? (
                <>
                    <p className={cls.ModalSubtitle}>
                        <b>{details.username || "—"}</b> · {details.email || "no email"}
                    </p>
                    <div className={cls.Field}>
                        <span className={cls.FieldLabel}>Permissions</span>
                        <span className={cls.RowMeta}>
                            {details.isAdministrator ? "Admin " : ""}
                            {details.hasEmails ? "Emails " : ""}
                            {details.hasTaskTrackers ? "TaskTrackers " : ""}
                            {details.hasNotes ? "Notes " : ""}
                            {!details.isAdministrator && !details.hasEmails && !details.hasTaskTrackers
                                && !details.hasNotes ? "none" : ""}
                        </span>
                    </div>
                    <div className={cls.Field}>
                        <span className={cls.FieldLabel}>Subscription</span>
                        <span className={cls.RowMeta}>{details.subscriptionActive ? "Active" : "Inactive"}</span>
                    </div>
                </>
            ) : (
                <p className={cls.Empty}>User not found.</p>
            )}
            <div className={cls.ModalActions}>
                <Button variant="secondary" onClick={openSessions}>Sessions</Button>
                <Button variant="primary" onClick={CloseDialog}>Close</Button>
            </div>
        </div>
    )
}
