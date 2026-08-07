import {useState} from "react"
import {Button} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import type {
    AddCouchDBConnectionRequest,
    CheckCouchDBConnectionRequest,
    CheckCouchDBConnectionResponse,
} from "@/app/api/artel/external_connections.pb.ts"
import StorageCheckButton from "@/components/StorageCheckButton/StorageCheckButton.tsx"
import StorageFormField from "@/components/StorageFormField/StorageFormField.tsx"
import cls from "@/components/CouchDBConnectForm/CouchDBConnectForm.module.css"

interface CouchDBConnectFormProps {
    addConnection: (req: AddCouchDBConnectionRequest) => Promise<unknown>
    checkConnection: (req: CheckCouchDBConnectionRequest) => Promise<CheckCouchDBConnectionResponse>
}

export default function CouchDBConnectForm({addConnection, checkConnection}: CouchDBConnectFormProps) {
    const [connecting, setConnecting] = useState(false)
    const [url, setUrl] = useState("")
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [verified, setVerified] = useState(false)
    const [checkError, setCheckError] = useState<string | null>(null)

    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

    function withReset(setter: (v: string) => void) {
        return (value: string) => {
            setter(value)
            setVerified(false)
        }
    }

    function handleConnect() {
        setConnecting(true)
        addConnection({url, username, password})
            .then(CloseDialog)
            .catch(e => bakeError("Failed to connect CouchDB", e))
            .finally(() => setConnecting(false))
    }

    return (
        <div className={cls.CouchDBConnectFormContainer}>
            <p className={cls.ModalSub}>
                Connect your own CouchDB instance to let Artel provision vaults on it. We&apos;ll verify
                the credentials against your server before saving them.
            </p>
            <StorageFormField
                label="URL"
                placeholder="https://couchdb.example.com"
                value={url}
                onChange={withReset(setUrl)}
                disabled={connecting}
                autoComplete="off"
            />
            <StorageFormField
                label="Username"
                placeholder="admin"
                value={username}
                onChange={withReset(setUsername)}
                disabled={connecting}
                autoComplete="off"
            />
            <StorageFormField
                label="Password"
                type="password"
                placeholder="••••••••"
                value={password}
                onChange={withReset(setPassword)}
                disabled={connecting}
                autoComplete="new-password"
            />
            <div className={cls.ModalActions}>
                <StorageCheckButton
                    req={{url, username, password}}
                    disabled={connecting || !url}
                    onResult={setVerified}
                    onError={setCheckError}
                    checkConnection={checkConnection}
                />
                <Button
                    variant="primary"
                    onClick={handleConnect}
                    disabled={connecting || !url || !verified}
                >
                    {connecting ? "Connecting…" : "Connect"}
                </Button>
            </div>
            {checkError && <p className={cls.CheckErrorText}>{checkError}</p>}
        </div>
    )
}
