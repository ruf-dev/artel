import {useState} from "react"
import {Button, ModalClose} from "@vervstack/chures"

import cls from "@/pages/admin/components/InstancesTab/components/InstanceFormDialog/InstanceFormDialog.module.css"
import {CouchInstancesAPI, GetCouchInstanceResponse} from "@/app/api/artel/couch_instances.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import FormField from "@/components/FormField/FormField.tsx"
import Input from "@/components/atoms/Input/Input.tsx"
import SetupStatusSection from "@/pages/admin/components/SetupStatusSection/SetupStatusSection.tsx"

interface InstanceFormDialogProps {
    initial?: GetCouchInstanceResponse
    onSave: (url: string, username: string, password: string) => Promise<void>
}

export default function InstanceFormDialog({initial, onSave}: InstanceFormDialogProps) {
    const {CloseDialog} = useDialog()
    const {auth} = useUser()
    const bakeError = useBakeError()
    const [saving, setSaving] = useState(false)
    const [settingUp, setSettingUp] = useState(false)
    const [setupRefresh, setSetupRefresh] = useState(0)
    const [url, setUrl] = useState(initial?.url ?? "")
    const [username, setUsername] = useState(initial?.username ?? "")
    const [password, setPassword] = useState("")

    const isEdit = !!initial
    const busy = saving || settingUp
    const title = isEdit ? "Edit instance" : "Add CouchDB instance"

    async function handleSave() {
        setSaving(true)
        try {
            await onSave(url, username, password)
            CloseDialog()
        } catch (err) {
            bakeError(isEdit ? "Failed to update instance" : "Failed to add instance", err)
        } finally {
            setSaving(false)
        }
    }

    async function handleSetup() {
        if (!initial?.id) return
        setSettingUp(true)
        try {
            await CouchInstancesAPI.SetupCouchInstance({id: initial.id}, auth.getInitReq())
            setSetupRefresh(r => r + 1)
        } catch (err) {
            bakeError("Failed to setup instance", err)
        } finally {
            setSettingUp(false)
        }
    }

    return (
        <div className={cls.InstanceFormDialogContainer} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>{title}</h2>
                <ModalClose onClick={CloseDialog} disabled={busy} className={cls.ModalClose}/>
            </div>
            {isEdit && initial?.id && <SetupStatusSection instanceId={initial.id} refreshKey={setupRefresh} />}
            <FormField
                label="URL"
                placeholder="https://couch.example.com:5984"
                defaultValue={url}
                onChange={setUrl}
                disabled={busy}
                fieldClassName={cls.Field}
                labelClassName={cls.FieldLabel}
            />
            <FormField
                label="Username"
                placeholder="admin"
                defaultValue={username}
                onChange={setUsername}
                disabled={busy}
                fieldClassName={cls.Field}
                labelClassName={cls.FieldLabel}
            />
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Password{isEdit ? " (leave blank to keep current)" : ""}</span>
                <Input
                    type="password"
                    placeholder={isEdit ? "••••••••" : "password"}
                    value={password}
                    setValue={setPassword}
                    disabled={busy}
                    autoComplete="new-password"
                />
            </label>
            <div className={cls.ModalActions}>
                {isEdit ? (
                    <>
                        <Button variant="secondary" onClick={handleSetup} disabled={busy}>
                            {settingUp ? "Setting up…" : "Setup"}
                        </Button>
                        <Button variant="primary" onClick={handleSave} disabled={busy}>
                            {saving ? "Saving…" : "Save changes"}
                        </Button>
                    </>
                ) : (
                    <Button variant="primary" onClick={handleSave} disabled={busy}>
                        {saving ? "Saving…" : "Add instance"}
                    </Button>
                )}
            </div>
        </div>
    )
}
