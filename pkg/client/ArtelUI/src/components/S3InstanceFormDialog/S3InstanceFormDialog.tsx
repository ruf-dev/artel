import {useState} from "react"

import {Button, ModalClose} from "@vervstack/chures"
import cls from "@/components/S3InstanceFormDialog/S3InstanceFormDialog.module.css"

import {S3InstancesAPI, GetS3InstanceResponse} from "@/app/api/artel/s3_instances.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"

import FormField from "@/components/FormField/FormField.tsx"

interface Props {
    initial?: GetS3InstanceResponse
    onSaved: () => void
}

export default function S3InstanceFormDialog({initial, onSaved}: Props) {
    const {auth} = useUser()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

    const isEdit = !!initial
    const [saving, setSaving] = useState(false)
    const [endpoint, setEndpoint] = useState(initial?.endpoint ?? "")
    const [region, setRegion] = useState(initial?.region ?? "")
    const [accessKey, setAccessKey] = useState("")
    const [secretKey, setSecretKey] = useState("")
    const [useSsl, setUseSsl] = useState(initial?.useSsl ?? true)
    const [pathStyle, setPathStyle] = useState(initial?.pathStyle ?? false)

    const title = isEdit ? "Edit S3 instance" : "Add S3 instance"

    function handleSave() {
        setSaving(true)
        const call = isEdit
            ? S3InstancesAPI.UpdateS3Instance(
                {id: initial?.id, endpoint, region, accessKey, secretKey, useSsl, pathStyle},
                auth.getInitReq(),
            )
            : S3InstancesAPI.RegisterS3Instance(
                {endpoint, region, accessKey, secretKey, useSsl, pathStyle},
                auth.getInitReq(),
            )

        call
            .then(() => {
                onSaved()
                CloseDialog()
            })
            .catch(err => bakeError(isEdit ? "Failed to update instance" : "Failed to add instance", err))
            .finally(() => setSaving(false))
    }

    return (
        <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>{title}</h2>
                <ModalClose onClick={CloseDialog} disabled={saving} className={cls.ModalClose}/>
            </div>
            <FormField
                label="Endpoint"
                placeholder="https://s3.example.com or http://garage.local:3900"
                defaultValue={endpoint}
                onChange={setEndpoint}
                disabled={saving}
                fieldClassName={cls.Field}
                labelClassName={cls.FieldLabel}
                inputClassName={cls.Input}
            />
            <FormField
                label="Region (optional)"
                placeholder="us-east-1"
                defaultValue={region}
                onChange={setRegion}
                disabled={saving}
                fieldClassName={cls.Field}
                labelClassName={cls.FieldLabel}
                inputClassName={cls.Input}
            />
            <FormField
                label="Access key"
                placeholder="access key"
                defaultValue={accessKey}
                onChange={setAccessKey}
                disabled={saving}
                fieldClassName={cls.Field}
                labelClassName={cls.FieldLabel}
                inputClassName={cls.Input}
            />
            <div className={cls.Field}>
                <span className={cls.FieldLabel}>Secret key{isEdit ? " (leave blank to keep current)" : ""}</span>
                <input
                    type="password"
                    placeholder={isEdit ? "••••••••" : "secret key"}
                    className={cls.Input}
                    onChange={e => setSecretKey(e.target.value)}
                    disabled={saving}
                    autoComplete="new-password"
                />
            </div>
            <S3ToggleFields
                useSsl={useSsl}
                pathStyle={pathStyle}
                disabled={saving}
                onUseSslChange={setUseSsl}
                onPathStyleChange={setPathStyle}
            />
            <div className={cls.ModalFooter}>
                <Button variant="ghost" onClick={CloseDialog} disabled={saving}>Cancel</Button>
                <Button variant="primary" onClick={handleSave} disabled={saving || !endpoint}>
                    {saving ? "Saving…" : isEdit ? "Save changes" : "Add instance"}
                </Button>
            </div>
        </div>
    )
}

function S3ToggleFields({useSsl, pathStyle, disabled, onUseSslChange, onPathStyleChange}: {
    useSsl: boolean
    pathStyle: boolean
    disabled: boolean
    onUseSslChange: (v: boolean) => void
    onPathStyleChange: (v: boolean) => void
}) {
    return (
        <div className={cls.ToggleGroup}>
            <label className={cls.ToggleRow}>
                <input
                    type="checkbox"
                    checked={useSsl}
                    disabled={disabled}
                    onChange={e => onUseSslChange(e.target.checked)}
                />
                <span className={cls.ToggleLabel}>Use SSL</span>
            </label>
            <label className={cls.ToggleRow}>
                <input
                    type="checkbox"
                    checked={pathStyle}
                    disabled={disabled}
                    onChange={e => onPathStyleChange(e.target.checked)}
                />
                <span className={cls.ToggleLabel}>Path-style addressing</span>
            </label>
            <p className={cls.ToggleCaption}>
                Required for Garage and most self-hosted S3-compatible stores.
            </p>
        </div>
    )
}
