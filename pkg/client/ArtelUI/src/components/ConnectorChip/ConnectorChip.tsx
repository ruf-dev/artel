import cls from "@/components/ConnectorChip/ConnectorChip.module.css"

import {McpConnectorInfo} from "@/app/api/artel/mcp_keys.pb.ts"
import {ExternalConnectionInfo, ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"

import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"

export function connectionLabel(c: ExternalConnectionInfo): string {
    if (c.google) return c.google.email ?? "Google account"
    if (c.provider === ExternalProvider.EXTERNAL_PROVIDER_EMAIL) {
        return c.generic?.fields?.username ?? "Email account"
    }
    return c.provider?.replace("EXTERNAL_PROVIDER_", "").toLowerCase() ?? "Connection"
}

const KNOWN_MAIL_DOMAIN_COLORS: Record<string, string> = {
    "yandex.ru": "#FFCC00",
    "yandex.com": "#FFCC00",
    "ya.ru": "#FFCC00",
    "gmail.com": "#4285F4",
    "google.com": "#4285F4",
    "googlemail.com": "#4285F4",
    "outlook.com": "#0078D4",
    "hotmail.com": "#0078D4",
    "live.com": "#0078D4",
    "msn.com": "#0078D4",
    "mail.ru": "#005FF9",
    "inbox.ru": "#005FF9",
    "list.ru": "#005FF9",
    "bk.ru": "#005FF9",
    "icloud.com": "#A2AAAD",
    "me.com": "#A2AAAD",
    "mac.com": "#A2AAAD",
    "protonmail.com": "#6D4AFF",
    "proton.me": "#6D4AFF",
    "yahoo.com": "#6001D2",
}

function mailDomainColor(domain: string): string {
    const known = KNOWN_MAIL_DOMAIN_COLORS[domain.toLowerCase()]
    if (known) return known

    let hash = 0
    for (let i = 0; i < domain.length; i++) {
        hash = (hash * 31 + domain.charCodeAt(i)) >>> 0
    }
    return `hsl(${hash % 360}, 65%, 45%)`
}

const PROVIDER_CHIP_CLASS: Partial<Record<ExternalProvider, string>> = {
    [ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS]: cls.SheetsChip,
    [ExternalProvider.EXTERNAL_PROVIDER_TRELLO]: cls.TrelloChip,
    [ExternalProvider.EXTERNAL_PROVIDER_MIRO]: cls.MiroChip,
}

export default function ConnectorChip({connector, connections}: {
    connector: McpConnectorInfo
    connections: ExternalConnectionInfo[]
}) {
    const match = connections.find(c => c.id === connector.externalConnectionId)

    if (!match) {
        return <GenericChip mcpName={connector.mcpName}/>
    }

    if (match.provider === ExternalProvider.EXTERNAL_PROVIDER_EMAIL) {
        return <EmailChip label={connectionLabel(match)}/>
    }

    if (match.provider && PROVIDER_CHIP_CLASS[match.provider]) {
        return (
            <ProviderChip
                provider={match.provider}
                variantClass={PROVIDER_CHIP_CLASS[match.provider]!}
                label={connectionLabel(match)}
            />
        )
    }

    return <GenericChip mcpName={connector.mcpName}/>
}

function GenericChip({mcpName}: {mcpName?: string}) {
    return (
        <span
            className={cls.Chip}
            data-tooltip-id="root-tooltip"
            data-tooltip-content={`${mcpName} connection`}
        >
            {mcpName}
        </span>
    )
}

function EmailChip({label}: {label: string}) {
    const atIndex = label.indexOf("@")
    if (atIndex === -1) {
        return <GenericChip mcpName={label}/>
    }

    const accent = mailDomainColor(label.slice(atIndex + 1))

    return (
        <span
            className={cls.MailChip}
            style={{"--chip-accent": accent} as React.CSSProperties}
            data-tooltip-id="root-tooltip"
            data-tooltip-content={`Email connection: ${label}`}
        >
            <span className={cls.MailAtBadge}>@</span>
            {label}
        </span>
    )
}

function ProviderChip({provider, variantClass, label}: {
    provider: ExternalProvider
    variantClass: string
    label: string
}) {
    return (
        <span
            className={`${cls.ProviderChip} ${variantClass}`}
            data-tooltip-id="root-tooltip"
            data-tooltip-content={`${label} connection`}
        >
            <span className={cls.ProviderChipIcon}><ProviderIcon provider={provider}/></span>
            {label}
        </span>
    )
}
