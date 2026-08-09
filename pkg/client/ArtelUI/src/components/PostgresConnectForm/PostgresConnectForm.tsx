import {useState} from "react"
import {Button} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import type {
    AddPostgresConnectionRequest,
    CheckPostgresConnectionRequest,
    CheckPostgresConnectionResponse,
} from "@/app/api/artel/external_connections.pb.ts"
import StorageCheckButton from "@/components/StorageCheckButton/StorageCheckButton.tsx"
import PostgresConnectFormFields from "@/components/PostgresConnectForm/PostgresConnectFormFields.tsx"
import cls from "@/components/PostgresConnectForm/PostgresConnectForm.module.css"

interface PostgresConnectFormProps {
    addConnection: (req: AddPostgresConnectionRequest) => Promise<unknown>
    checkConnection: (req: CheckPostgresConnectionRequest) => Promise<CheckPostgresConnectionResponse>
}

export default function PostgresConnectForm({addConnection, checkConnection}: PostgresConnectFormProps) {
    const [connecting, setConnecting] = useState(false)
    const [host, setHost] = useState("")
    const [port, setPort] = useState(5432)
    const [database, setDatabase] = useState("")
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [sslMode, setSslMode] = useState("prefer")
    const [verified, setVerified] = useState(false)
    const [checkError, setCheckError] = useState<string | null>(null)

    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

    function withResetString(setter: (v: string) => void) {
        return (value: string) => {
            setter(value)
            setVerified(false)
        }
    }

    function withResetNumber(setter: (v: number) => void) {
        return (value: number) => {
            setter(value)
            setVerified(false)
        }
    }

    function handleConnect() {
        setConnecting(true)
        addConnection({host, port, database, username, password, sslMode})
            .then(CloseDialog)
            .catch(e => bakeError("Failed to connect Postgres", e))
            .finally(() => setConnecting(false))
    }

    return (
        <div className={cls.PostgresConnectFormContainer}>
            <p className={cls.ModalSub}>
                Connect your own PostgreSQL instance to let Artel provision databases on it. We&apos;ll verify
                the credentials against your server before saving them.
            </p>
            <PostgresConnectFormFields
                host={host}
                port={port}
                database={database}
                username={username}
                password={password}
                sslMode={sslMode}
                connecting={connecting}
                onHostChange={withResetString(setHost)}
                onPortChange={withResetNumber(setPort)}
                onDatabaseChange={withResetString(setDatabase)}
                onUsernameChange={withResetString(setUsername)}
                onPasswordChange={withResetString(setPassword)}
                onSslModeChange={withResetString(setSslMode)}
            />
            <div className={cls.ModalActions}>
                <StorageCheckButton
                    req={{host, port, database, username, password, sslMode}}
                    disabled={connecting || !host}
                    onResult={setVerified}
                    onError={setCheckError}
                    checkConnection={checkConnection}
                />
                <Button
                    variant="primary"
                    onClick={handleConnect}
                    disabled={connecting || !host || !verified}
                >
                    {connecting ? "Connecting…" : "Connect"}
                </Button>
            </div>
            {checkError && <p className={cls.CheckErrorText}>{checkError}</p>}
        </div>
    )
}
