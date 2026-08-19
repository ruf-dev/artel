import {Loader} from "@vervstack/chures"

import cls from "@/pages/workbench/components/WorkbenchAuthScreen/components/AuthSignInPanel/AuthSignInPanel.module.css"

interface Props {
    url?: string
}

export default function AuthSignInPanel({url}: Props) {
    if (!url) {
        return (
            <div className={cls.AuthSignInPanelContainer}>
                <Loader variant="arcs" size="sm" color="var(--coral)"/>
                <span>Waiting for sign-in link…</span>
            </div>
        )
    }

    return (
        <div className={cls.AuthSignInPanelContainer}>
            <a className={cls.LinkButton} href={url} target="_blank" rel="noreferrer">Sign in with Claude</a>
        </div>
    )
}
