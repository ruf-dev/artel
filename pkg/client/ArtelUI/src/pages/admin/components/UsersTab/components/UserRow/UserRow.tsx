import {useState} from "react"
import {Button, ConfirmDialog} from "@vervstack/chures"

import cls from "@/pages/admin/components/UsersTab/components/UserRow/UserRow.module.css"
import {AdminCouchAPI, CouchUserEntry} from "@/app/api/artel/admin_couch.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import ChangePasswordDialog
    from "@/pages/admin/components/UsersTab/components/ChangePasswordDialog/ChangePasswordDialog.tsx"
import ManageAccessDialog from "@/pages/admin/components/UsersTab/components/ManageAccessDialog/ManageAccessDialog.tsx"

interface UserRowProps {
    instanceId: string
    user: CouchUserEntry
    onRefresh: () => Promise<void>
}

export default function UserRow({instanceId, user, onRefresh}: UserRowProps) {
    const {auth} = useUser()
    const {OpenDialog, CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const [deleting, setDeleting] = useState(false)

    function handleDelete() {
        const name = user.name
        if (!name) return
        OpenDialog(
            <ConfirmDialog
                title="Delete user"
                message={`Delete "${name}"? This cannot be undone.`}
                danger
                confirmLabel="Delete"
                onClose={CloseDialog}
                onConfirm={async () => {
                    setDeleting(true)
                    try {
                        await AdminCouchAPI.DeleteCouchUser({instanceId, username: name}, auth.getInitReq())
                        await onRefresh()
                    } catch (err) {
                        bakeError("Failed to delete user", err)
                    } finally {
                        setDeleting(false)
                    }
                }}
            />
        )
    }

    function openChangePassword() {
        if (!user.name) return
        OpenDialog(
            <ChangePasswordDialog instanceId={instanceId} username={user.name} />
        )
    }

    function openManageAccess() {
        if (!user.name) return
        OpenDialog(
            <ManageAccessDialog instanceId={instanceId} username={user.name} />
        )
    }

    return (
        <div className={cls.UserRowContainer}>
            <div className={cls.RowInfo}>
                <span className={cls.RowUrl}>{user.name}</span>
                <span className={cls.RowMeta}>{(user.roles ?? []).join(", ") || "no roles"}</span>
            </div>
            <div className={cls.RowActions}>
                <Button variant="secondary" onClick={openManageAccess}>
                    DB Access
                </Button>
                <Button variant="secondary" onClick={openChangePassword}>
                    Change Password
                </Button>
                <Button variant="danger" onClick={handleDelete} disabled={deleting}>
                    {deleting ? "Deleting…" : "Delete"}
                </Button>
            </div>
        </div>
    )
}
