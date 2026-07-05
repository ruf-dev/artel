import cls from "@/widgets/GenericConnectionContent/GenericConnectionContent.module.css"

import {ExternalConnectionInfo} from "@/app/api/artel/external_connections.pb.ts"

import {Button} from "@vervstack/chures"

export default function GenericConnectionContent({connection, onDisconnect}: {
    connection: ExternalConnectionInfo
    onDisconnect: () => void
}) {
    const fields = connection.generic?.fields ?? {}
    const connectedAt = connection.createdAt
        ? new Date(connection.createdAt).toLocaleDateString(undefined, {
            year: "numeric",
            month: "short",
            day: "numeric",
        })
        : "—"

    return (
        <>
            <div className={cls.InfoRows}>
                {Object.entries(fields).map(([key, value]) => (
                    <div key={key} className={cls.InfoRow}>
                        <span className={cls.InfoLabel}>{key}</span>
                        <span className={cls.InfoValue}>{value}</span>
                    </div>
                ))}
                <div className={cls.InfoRow}>
                    <span className={cls.InfoLabel}>Connected</span>
                    <span className={cls.InfoValue}>{connectedAt}</span>
                </div>
            </div>
            <div className={cls.Actions}>
                <Button variant="danger" onClick={onDisconnect}>Disconnect</Button>
            </div>
        </>
    )
}
