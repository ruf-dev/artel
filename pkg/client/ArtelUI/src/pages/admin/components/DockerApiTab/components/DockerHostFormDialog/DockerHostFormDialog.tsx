import {Fragment, useState} from "react"
import {Button, ModalClose} from "@vervstack/chures"

import cls from "@/pages/admin/components/DockerApiTab/components/DockerHostFormDialog/DockerHostFormDialog.module.css"
import {GetDockerHostResponse} from "@/app/api/artel/docker_hosts.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import FormField from "@/components/FormField/FormField.tsx"
import DockerHostTlsFields
    from "@/pages/admin/components/DockerApiTab/components/DockerHostFormDialog/components/DockerHostTlsFields.tsx"
import DockerHostVpsGuide
    from "@/pages/admin/components/DockerApiTab/components/DockerHostFormDialog/components/DockerHostVpsGuide.tsx"

export interface DockerHostFormDialogSaveData {
    url: string
    caCert?: string
    clientCert?: string
    clientKey?: string
}

interface DockerHostFormDialogProps {
    initial?: GetDockerHostResponse
    onSave: (data: DockerHostFormDialogSaveData) => Promise<void>
}

export default function DockerHostFormDialog({initial, onSave}: DockerHostFormDialogProps) {
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const [saving, setSaving] = useState(false)
    const [url, setUrl] = useState(initial?.url ?? "")
    const [caCert, setCaCert] = useState("")
    const [clientCert, setClientCert] = useState("")
    const [clientKey, setClientKey] = useState("")

    const isEdit = !!initial
    const title = isEdit ? "Edit Docker host" : "Add Docker host"
    const showTlsFields = url.startsWith("tcp://")

    async function handleSave() {
        setSaving(true)
        try {
            await onSave({url, caCert, clientCert, clientKey})
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
            <p className={cls.ModalSub}>
                Point Artel at a Docker daemon it can create workbench containers on. Use a local socket path
                (e.g. unix:///var/run/docker-workbenches.sock) for a daemon on this same host, or tcp://host:2376
                with mutual TLS for a daemon on a remote host.
            </p>
            <FormField
                label="URL"
                placeholder="unix:///var/run/docker-workbenches.sock"
                defaultValue={url}
                onChange={setUrl}
                disabled={saving}
                fieldClassName={cls.Field}
                labelClassName={cls.FieldLabel}
            />
            {showTlsFields && (
                <Fragment>
                    <DockerHostVpsGuide/>
                    <DockerHostTlsFields
                        caCert={caCert}
                        clientCert={clientCert}
                        clientKey={clientKey}
                        isEdit={isEdit}
                        disabled={saving}
                        onCaCertChange={setCaCert}
                        onClientCertChange={setClientCert}
                        onClientKeyChange={setClientKey}
                    />
                </Fragment>
            )}
            <div className={cls.ModalActions}>
                <Button variant="primary" onClick={handleSave} disabled={saving}>
                    {saving ? "Saving…" : isEdit ? "Save changes" : "Add host"}
                </Button>
            </div>
        </div>
    )
}
