import cls from "@/components/EmailChip/EmailChip.module.css"
import {cn} from "@/app/utils/cn"

import {useDialog} from "@/app/hooks/Dialog"

import ManageEmailDialog from "@/dialogs/ManageEmailDialog/ManageEmailDialog.tsx"
import GenericChip from "@/components/GenericChip/GenericChip.tsx"

import {mailProviderIcon} from "@/app/utils/mailProviderIcon"

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

export default function EmailChip({label}: {label: string}) {
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
