import {useState, useEffect} from "react"
import {Button, ModalClose} from "@vervstack/chures"

import cls
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSessionsDialog/UserSessionsDialog.module.css"
import {AdminUsersAPI, UserSession} from "@/app/api/artel/admin_users.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"

interface UserSessionsDialogProps {
    userId: string
}

export default function UserSessionsDialog({userId}: UserSessionsDialogProps) {
    const {auth} = useUser()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const [sessions, setSessions] = useState<UserSession[]>([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        async function load() {
            try {
                const res = await AdminUsersAPI.GetUserSessions({userId}, auth.getInitReq())
                setSessions(res.sessions ?? [])
            } catch (err) {
                bakeError("Failed to load sessions", err)
            } finally {
                setLoading(false)
            }
        }
        void load()
    }, [auth, userId])

    function formatDate(iso?: string) {
        if (!iso) return "—"
        return new Date(iso).toLocaleString()
    }

    return (
        <div className={cls.UserSessionsDialogContainer} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>Sessions</h2>
                <ModalClose onClick={CloseDialog} className={cls.ModalClose} />
            </div>
            <div className={cls.SessionsListContainer}>
                {loading ? (
                    <p className={cls.Empty}>Loading…</p>
                ) : sessions.length === 0 ? (
                    <p className={cls.Empty}>No active sessions.</p>
                ) : (
                    sessions.map(s => (
                        <div key={s.sessionId} className={cls.RowContainer}>
                            <div className={cls.RowInfo}>
                                <span className={cls.RowUrl}>
                                    Created: {formatDate(s.createdAt as unknown as string)}
                                </span>
                                <span className={cls.RowMeta}>
                                    Expires: {formatDate(s.expiresAt as unknown as string)}
                                </span>
                            </div>
                        </div>
                    ))
                )}
            </div>
            <div className={cls.ModalActions}>
                <Button variant="primary" onClick={CloseDialog}>Close</Button>
            </div>
        </div>
    )
}
