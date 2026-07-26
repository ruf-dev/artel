import Textarea from "@/components/atoms/Textarea/Textarea.tsx"
import cls
from "@/pages/admin/components/DockerApiTab/components/DockerHostFormDialog/components/DockerHostTlsFields.module.css"

interface Props {
    caCert: string
    clientCert: string
    clientKey: string
    isEdit: boolean
    disabled: boolean
    onCaCertChange: (v: string) => void
    onClientCertChange: (v: string) => void
    onClientKeyChange: (v: string) => void
}

export default function DockerHostTlsFields(props: Props) {
    const suffix = props.isEdit ? " (leave blank to keep current)" : ""

    return (
        <div className={cls.DockerHostTlsFieldsContainer}>
            <div className={cls.Field}>
                <span className={cls.FieldLabel}>CA certificate{suffix}</span>
                <Textarea
                    placeholder="-----BEGIN CERTIFICATE-----"
                    value={props.caCert}
                    setValue={props.onCaCertChange}
                    disabled={props.disabled}
                />
            </div>
            <div className={cls.Field}>
                <span className={cls.FieldLabel}>Client certificate{suffix}</span>
                <Textarea
                    placeholder="-----BEGIN CERTIFICATE-----"
                    value={props.clientCert}
                    setValue={props.onClientCertChange}
                    disabled={props.disabled}
                />
            </div>
            <div className={cls.Field}>
                <span className={cls.FieldLabel}>Client key{suffix}</span>
                <Textarea
                    placeholder="-----BEGIN PRIVATE KEY-----"
                    value={props.clientKey}
                    setValue={props.onClientKeyChange}
                    disabled={props.disabled}
                />
            </div>
        </div>
    )
}
