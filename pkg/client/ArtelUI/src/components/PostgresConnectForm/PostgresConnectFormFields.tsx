import {Dropdown, type DropdownOption} from "@vervstack/chures"

import StorageFormField from "@/components/StorageFormField/StorageFormField.tsx"

const SSL_MODE_OPTIONS: DropdownOption[] = [
    {id: "disable", name: "Disable"},
    {id: "allow", name: "Allow"},
    {id: "prefer", name: "Prefer"},
    {id: "require", name: "Require"},
    {id: "verify-ca", name: "Verify CA"},
    {id: "verify-full", name: "Verify Full"},
]

interface PostgresConnectFormFieldsProps {
    host: string
    port: number
    database: string
    username: string
    password: string
    sslMode: string
    connecting: boolean
    onHostChange: (v: string) => void
    onPortChange: (v: number) => void
    onDatabaseChange: (v: string) => void
    onUsernameChange: (v: string) => void
    onPasswordChange: (v: string) => void
    onSslModeChange: (v: string) => void
}

export default function PostgresConnectFormFields(props: PostgresConnectFormFieldsProps) {
    const context = props
    return <>
        <StorageFormField
            label="Host"
            placeholder="postgres.example.com"
            value={context.host}
            onChange={context.onHostChange}
            disabled={context.connecting}
            autoComplete="off"
        />
        <StorageFormField
            label="Port"
            type="text"
            placeholder="5432"
            value={String(context.port)}
            onChange={v => context.onPortChange(Number(v) || 5432)}
            disabled={context.connecting}
            autoComplete="off"
        />
        <StorageFormField
            label="Database"
            placeholder="postgres"
            value={context.database}
            onChange={context.onDatabaseChange}
            disabled={context.connecting}
            autoComplete="off"
        />
        <StorageFormField
            label="Username"
            placeholder="admin"
            value={context.username}
            onChange={context.onUsernameChange}
            disabled={context.connecting}
            autoComplete="off"
        />
        <StorageFormField
            label="Password"
            type="password"
            placeholder="••••••••"
            value={context.password}
            onChange={context.onPasswordChange}
            disabled={context.connecting}
            autoComplete="new-password"
        />
        <Dropdown
            label="SSL Mode"
            options={SSL_MODE_OPTIONS}
            value={[context.sslMode]}
            onChange={v => context.onSslModeChange(v[0] ?? "prefer")}
        />
    </>
}

export {SSL_MODE_OPTIONS}
