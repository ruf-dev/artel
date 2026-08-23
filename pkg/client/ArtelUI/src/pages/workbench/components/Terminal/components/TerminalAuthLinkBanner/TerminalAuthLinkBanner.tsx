import {useState} from "react"
import {Button} from "@vervstack/chures"

import cls
    from "@/pages/workbench/components/Terminal/components/TerminalAuthLinkBanner/TerminalAuthLinkBanner.module.css"

interface Props {
    url?: string
}

export default function TerminalAuthLinkBanner({url}: Props) {
    const [dismissedUrl, setDismissedUrl] = useState<string | null>(null)

    if (!url || url === dismissedUrl) return null

    return (
        <div className={cls.TerminalAuthLinkBannerContainer}>
            <p className={cls.Text}>Claude needs you to authorize access to continue.</p>
            <a className={cls.LinkButton} href={url} target="_blank" rel="noreferrer">Authorize Claude</a>
            <Button
                variant="secondary"
                className={cls.DismissButton}
                onClick={() => setDismissedUrl(url)}
                aria-label="Dismiss"
            >
                <svg viewBox="0 0 24 24" width="10" height="10" fill="currentColor">
                    <line x1="6" y1="6" x2="18" y2="18" stroke="currentColor" strokeWidth="2"/>
                    <line x1="18" y1="6" x2="6" y2="18" stroke="currentColor" strokeWidth="2"/>
                </svg>
            </Button>
        </div>
    )
}
