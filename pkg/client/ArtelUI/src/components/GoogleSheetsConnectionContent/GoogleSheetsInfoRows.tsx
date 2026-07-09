import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import cls from "@/components/GoogleSheetsConnectionContent/GoogleSheetsInfoRows.module.css"
import {SCOPE_INFO, trimScope} from "@/app/utils/googleScopes.ts"

function buildScopeTooltipHtml(scopeList: string[]): string {
    return scopeList
        .map(s => {
            const name = trimScope(s)
            const desc = SCOPE_INFO[name]
            return desc ? `<b>${name}</b>: ${desc}` : `<b>${name}</b>`
        })
        .join("<br/>")
}

export default function GoogleSheetsInfoRows({email, connectedAt, scopeList}: {
    email: string | undefined
    connectedAt: string
    scopeList: string[]
}) {
    const scopeTooltipHtml = buildScopeTooltipHtml(scopeList)

    return (
        <div className={cls.GoogleSheetsInfoRowsContainer}>
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
    )
}
