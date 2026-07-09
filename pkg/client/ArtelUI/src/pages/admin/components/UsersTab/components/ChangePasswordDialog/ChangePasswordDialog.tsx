import {useState} from "react"
import {Button, ModalClose, Input} from "@vervstack/chures"

import cls from "@/pages/admin/components/UsersTab/components/ChangePasswordDialog/ChangePasswordDialog.module.css"
import {AdminCouchAPI} from "@/app/api/artel/admin_couch.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"

interface ChangePasswordDialogProps {
    instanceId: string
    username: string
}

export default function ChangePasswordDialog({instanceId, username}: ChangePasswordDialogProps) {
    const {auth} = useUser()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const [saving, setSaving] = useState(false)
    const [newPassword, setNewPassword] = useState("")

    async function handleSave() {
        if (!newPassword) return
        setSaving(true)
        try {
            await AdminCouchAPI.ChangeCouchUserPassword({instanceId, username, newPassword}, auth.getInitReq())
            CloseDialog()
        } catch (err) {
            bakeError("Failed to change password", err)
        } finally {
            setSaving(false)
        }
    }

    return (
        <div className={cls.ChangePasswordDialogContainer} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>Change password</h2>
                <ModalClose onClick={CloseDialog} disabled={saving} className={cls.ModalClose}/>
            </div>
            <p className={cls.ModalSubtitle}>User: <b>{username}</b></p>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>New password</span>
                <Input
                    type="password"
                    placeholder="new password"
                    className={cls.Input}
                    value={newPassword}
                    setValue={setNewPassword}
                    disabled={saving}
                    autoComplete="new-password"
                />
            </label>
            <div className={cls.ModalActions}>
                <Button variant="primary" onClick={handleSave} disabled={saving || !newPassword}>
                    {saving ? "Saving…" : "Save password"}
                </Button>
            </div>
        </div>
    )
}
