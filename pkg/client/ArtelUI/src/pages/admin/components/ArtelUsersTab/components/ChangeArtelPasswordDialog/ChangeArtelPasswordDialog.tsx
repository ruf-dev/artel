import {useState} from "react"
import {Button, ModalClose} from "@vervstack/chures"

import Input from "@/components/atoms/Input/Input.tsx"
// eslint-disable-next-line max-len -- path too long to wrap
import cls from "@/pages/admin/components/ArtelUsersTab/components/ChangeArtelPasswordDialog/ChangeArtelPasswordDialog.module.css"
import {AdminUsersAPI} from "@/app/api/artel/admin_users.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"

interface ChangeArtelPasswordDialogProps {
    userId: string
    username: string
}

export default function ChangeArtelPasswordDialog({userId, username}: ChangeArtelPasswordDialogProps) {
    const {auth} = useUser()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const [saving, setSaving] = useState(false)
    const [newPassword, setNewPassword] = useState("")

    function handleSave() {
        if (!newPassword) return
        setSaving(true)
        AdminUsersAPI.ChangeArtelUserPassword({userId, newPassword}, auth.getInitReq())
            .then(() => {
                CloseDialog()
            })
            .catch((err) => {
                bakeError("Failed to change password", err)
            })
            .finally(() => {
                setSaving(false)
            })
    }

    return (
        <div className={cls.ChangeArtelPasswordDialogContainer} role="dialog" aria-modal="true">
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
