import {Button} from "@vervstack/chures"

import cls from "@/widgets/GoogleConnectionContent/GoogleConnectionContent.module.css"
import {ExternalConnectionInfo} from "@/app/api/artel/external_connections.pb.ts"
import {parseScopeList, SCOPE_INFO, trimScope} from "@/app/utils/googleScopes.ts"
import {cn} from "@/app/utils/cn.ts"

export default function GoogleConnectionContent({connection, onDisconnect}: {
    connection: ExternalConnectionInfo
    onDisconnect: () => void
}) {
    const email = connection.google?.email
    const scopes = connection.google?.scopes
    const connectedAt = connection.createdAt
        ? new Date(connection.createdAt).toLocaleDateString(undefined, {
            year: "numeric",
            month: "short",
            day: "numeric",
        })
        : "—"

    const scopeList = scopes ? parseScopeList(scopes) : []
    const scopeTooltipHtml = scopeList
        .map(s => {
            const name = trimScope(s)
            const desc = SCOPE_INFO[name]
            return desc ? `<b>${name}</b>: ${desc}` : `<b>${name}</b>`
        })
        .join("<br/>")

    return (
        <>
            <div className={cls.InfoRows}>
                {email && (
                    <div className={cls.InfoRow}>
                        <span className={cls.InfoLabel}>Account</span>
                        <span className={cls.InfoValue}>{email}</span>
                    </div>
                )}
                {scopeList.length > 0 && (
                    <div className={cls.InfoRow}>
                        <span className={cn(cls.InfoLabel, cls.ScopesLabel)}>
                            Scopes
                            <Button
                                variant="ghost"
                                className={cls.ScopeHelp}
                                data-tooltip-id="root-tooltip"
                                data-tooltip-html={scopeTooltipHtml}
                                aria-label="What are scopes?"
                            >?</Button>
                        </span>
                        <span className={cls.InfoValue}>{scopeList.map(trimScope).join(", ")}</span>
                    </div>
                )}
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
