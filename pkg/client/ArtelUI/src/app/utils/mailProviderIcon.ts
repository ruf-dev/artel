import yandexMailIcon from "@/icons/mails/yandex.png";
import gMailIcon from "@/icons/mails/gmail.png";

const MAIL_DOMAIN_ICONS: Record<string, string | null> = {
    "yandex.ru": yandexMailIcon,
    "yandex.com": yandexMailIcon,
    "ya.ru": yandexMailIcon,
    "gmail.com": gMailIcon,
    "google.com": gMailIcon,
    "googlemail.com": gMailIcon,

    "outlook.com": null,
    "hotmail.com": null,
    "live.com": null,
    "msn.com": null,
    "mail.ru": null,
    "inbox.ru": null,
    "list.ru": null,
    "bk.ru": null,
    "icloud.com": null,
    "me.com": null,
    "mac.com": null,
    "protonmail.com": null,
    "proton.me": null,
    "yahoo.com": null,
}

export function mailProviderIcon(email: string): string | null {
    const at = email.indexOf("@")
    if (at === -1) return null
    return MAIL_DOMAIN_ICONS[email.slice(at + 1).toLowerCase()] ?? null
}
