import {useState} from "react"
import {Button, ModalClose} from "@vervstack/chures"

import cls from "@/pages/admin/components/DockerApiTab/components/DockerHostFormDialog/DockerHostFormDialog.module.css"
import {GetDockerHostResponse} from "@/app/api/artel/docker_hosts.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import FormField from "@/components/FormField/FormField.tsx"

interface DockerHostFormDialogProps {
    initial?: GetDockerHostResponse
    onSave: (url: string) => Promise<void>
}

export default function DockerHostFormDialog({initial, onSave}: DockerHostFormDialogProps) {
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const [saving, setSaving] = useState(false)
    const [url, setUrl] = useState(initial?.url ?? "")

    const isEdit = !!initial
    const title = isEdit ? "Edit Docker host" : "Add Docker host"

    async function handleSave() {
        setSaving(true)
        try {
            await onSave(url)
            CloseDialog()
        } catch (err) {
            bakeError(isEdit ? "Failed to update Docker host" : "Failed to add Docker host", err)
        } finally {
            setSaving(false)
        }
    }

    return (
        <div className={cls.DockerHostFormDialogContainer} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>{title}</h2>
                <ModalClose onClick={CloseDialog} disabled={saving} className={cls.ModalClose}/>
            </div>
            <FormField
                label="URL"
                placeholder="unix:///var/run/docker-workbenches.sock"
                defaultValue={url}
                onChange={setUrl}
                disabled={saving}
                fieldClassName={cls.Field}
                labelClassName={cls.FieldLabel}
            />
            <div className={cls.ModalActions}>
                <Button variant="primary" onClick={handleSave} disabled={saving}>
                    {saving ? "Saving…" : isEdit ? "Save changes" : "Add host"}
                </Button>
            </div>
        </div>
    )
}
