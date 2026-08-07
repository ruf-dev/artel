import {useState} from "react"
import {Button} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import type {
    AddS3ConnectionRequest,
    CheckS3ConnectionRequest,
    CheckS3ConnectionResponse,
} from "@/app/api/artel/external_connections.pb.ts"
import StorageCheckButton from "@/components/StorageCheckButton/StorageCheckButton.tsx"
import StorageFormField from "@/components/StorageFormField/StorageFormField.tsx"
import S3ToggleFields from "@/components/S3ToggleFields/S3ToggleFields.tsx"
import cls from "@/components/S3ConnectForm/S3ConnectForm.module.css"

interface S3ConnectFormProps {
    addConnection: (req: AddS3ConnectionRequest) => Promise<unknown>
    checkConnection: (req: CheckS3ConnectionRequest) => Promise<CheckS3ConnectionResponse>
}

export default function S3ConnectForm({addConnection, checkConnection}: S3ConnectFormProps) {
    const [connecting, setConnecting] = useState(false)
    const [endpoint, setEndpoint] = useState("")
    const [region, setRegion] = useState("")
    const [accessKey, setAccessKey] = useState("")
    const [secretKey, setSecretKey] = useState("")
    const [useSsl, setUseSsl] = useState(true)
    const [pathStyle, setPathStyle] = useState(false)
    const [verified, setVerified] = useState(false)
    const [checkError, setCheckError] = useState<string | null>(null)

    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

    function withReset<T>(setter: (v: T) => void) {
        return (value: T) => {
            setter(value)
            setVerified(false)
        }
    }

    function handleConnect() {
        setConnecting(true)
        addConnection({endpoint, region, accessKey, secretKey, useSsl, pathStyle})
            .then(CloseDialog)
            .catch(e => bakeError("Failed to connect S3", e))
            .finally(() => setConnecting(false))
    }

    return (
        <div className={cls.S3ConnectFormContainer}>
            <p className={cls.ModalSub}>
                Connect your own S3-compatible bucket to let Artel provision vault storage on it. We&apos;ll
                verify the credentials before saving them.
            </p>
            <StorageFormField
                label="Endpoint"
                placeholder="https://s3.example.com or http://garage.local:3900"
                value={endpoint}
                onChange={withReset(setEndpoint)}
                disabled={connecting}
                autoComplete="off"
            />
            <StorageFormField
                label="Region (optional)"
                placeholder="us-east-1"
                value={region}
                onChange={withReset(setRegion)}
                disabled={connecting}
                autoComplete="off"
            />
            <StorageFormField
                label="Access key"
                placeholder="access key"
                value={accessKey}
                onChange={withReset(setAccessKey)}
                disabled={connecting}
                autoComplete="off"
            />
            <StorageFormField
                label="Secret key"
                type="password"
                placeholder="secret key"
                value={secretKey}
                onChange={withReset(setSecretKey)}
                disabled={connecting}
                autoComplete="new-password"
            />
            <S3ToggleFields
                useSsl={useSsl}
                pathStyle={pathStyle}
                disabled={connecting}
                onUseSslChange={withReset(setUseSsl)}
                onPathStyleChange={withReset(setPathStyle)}
            />
            <div className={cls.ModalActions}>
                <StorageCheckButton
                    req={{endpoint, region, accessKey, secretKey, useSsl, pathStyle}}
                    disabled={connecting || !endpoint}
                    onResult={setVerified}
                    onError={setCheckError}
                    checkConnection={checkConnection}
                />
                <Button
                    variant="primary"
                    onClick={handleConnect}
                    disabled={connecting || !endpoint || !verified}
                >
                    {connecting ? "Connecting…" : "Connect"}
                </Button>
            </div>
            {checkError && <p className={cls.CheckErrorText}>{checkError}</p>}
        </div>
    )
}
