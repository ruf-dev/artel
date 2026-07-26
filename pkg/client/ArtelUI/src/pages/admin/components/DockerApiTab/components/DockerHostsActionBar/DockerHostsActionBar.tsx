import {Button} from "@vervstack/chures"

import cls from "@/pages/admin/components/DockerApiTab/components/DockerHostsActionBar/DockerHostsActionBar.module.css"

interface DockerHostsActionBarProps {
    count: number
    onAddClick: () => void
}

export default function DockerHostsActionBar({count, onAddClick}: DockerHostsActionBarProps) {
    return (
        <div className={cls.DockerHostsActionBarContainer}>
            <p className={cls.HeroSub}>
                <b>{count} {count === 1 ? "host" : "hosts"}</b>
                {" · "}<span>manage Docker hosts for Workbench</span>
            </p>
            <Button variant="primary" onClick={onAddClick}>
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2"
                     strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Add host
            </Button>
        </div>
    )
}
