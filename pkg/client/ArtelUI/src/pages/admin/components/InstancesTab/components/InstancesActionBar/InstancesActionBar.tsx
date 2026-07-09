import {Button} from "@vervstack/chures"

import cls from "@/pages/admin/components/InstancesTab/components/InstancesActionBar/InstancesActionBar.module.css"

interface InstancesActionBarProps {
    count: number
    onAddClick: () => void
}

export default function InstancesActionBar({count, onAddClick}: InstancesActionBarProps) {
    return (
        <div className={cls.InstancesActionBarContainer}>
            <p className={cls.HeroSub}>
                <b>{count} {count === 1 ? "instance" : "instances"}</b>
                {" · "}<span>manage CouchDB cluster nodes</span>
            </p>
            <Button variant="primary" onClick={onAddClick}>
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2"
                     strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Add instance
            </Button>
        </div>
    )
}
