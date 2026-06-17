import cls from "@/components/ConnectorChip/ConnectorChip.module.css"
import {cn} from "@/app/utils/cn"

import {McpConnectorInfo} from "@/app/api/artel/mcp_keys.pb.ts"
import {ExternalConnectionInfo, ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"

import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"
import ManageEmailDialog from "@/dialogs/ManageEmailDialog/ManageEmailDialog.tsx"

import {mailProviderIcon} from "@/app/utils/mailProviderIcon"

export function connectionLabel(c: ExternalConnectionInfo): string {
    if (c.google) return c.google.email ?? "Google account"
    if (c.provider === ExternalProvider.EXTERNAL_PROVIDER_EMAIL) {
        return c.generic?.fields?.username ?? "Email account"
    }
    return c.provider?.replace("EXTERNAL_PROVIDER_", "").toLowerCase() ?? "Connection"
}

const KNOWN_MAIL_DOMAIN_CLASSES: Record<string, string> = {
    "yandex.ru": cls.AccentYandex,
    "yandex.com": cls.AccentYandex,
    "ya.ru": cls.AccentYandex,
    "gmail.com": cls.AccentGmail,
    "google.com": cls.AccentGmail,
    "googlemail.com": cls.AccentGmail,
    "outlook.com": cls.AccentOutlook,
    "hotmail.com": cls.AccentOutlook,
    "live.com": cls.AccentOutlook,
    "msn.com": cls.AccentOutlook,
    "mail.ru": cls.AccentMailRu,
    "inbox.ru": cls.AccentMailRu,
    "list.ru": cls.AccentMailRu,
    "bk.ru": cls.AccentMailRu,
    "icloud.com": cls.AccentICloud,
    "me.com": cls.AccentICloud,
    "mac.com": cls.AccentICloud,
    "protonmail.com": cls.AccentProton,
    "proton.me": cls.AccentProton,
    "yahoo.com": cls.AccentYahoo,
}


function mailDomainAccent(domain: string): { cls: string; style?: React.CSSProperties } {
    const known = KNOWN_MAIL_DOMAIN_CLASSES[domain.toLowerCase()]
    if (known) return {cls: known}

    let hash = 0
    for (let i = 0; i < domain.length; i++) {
        hash = (hash * 31 + domain.charCodeAt(i)) >>> 0
    }
    return {cls: "", style: {"--chip-accent": `hsl(${hash % 360}, 65%, 45%)`} as React.CSSProperties}
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
            className={cls.GenericChip}
            data-tooltip-id="root-tooltip"
            data-tooltip-content={`${mcpName} connection`}
        >
            {mcpName}
        </span>
    )
}

function EmailChip({label}: {label: string}) {
    const {OpenDialog} = useDialog()
    const atIndex = label.indexOf("@")

    if (atIndex === -1) {
        return <GenericChip mcpName={label}/>
    }

    const domain = label.slice(atIndex + 1)
    const accent = mailDomainAccent(domain)
    const icon = mailProviderIcon(label)

    return (
        <span
            className={cn(cls.MailChip, accent.cls)}
            style={accent.style}
            onClick={() => OpenDialog(<ManageEmailDialog/>)}
            data-tooltip-id="root-tooltip"
            data-tooltip-content={`Email connection: ${label}`}
        >
            {icon
                ? <img src={icon} className={cls.MailProviderIcon} alt={domain}/>
                : <span className={cls.MailAtBadge}>@</span>
            }
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
