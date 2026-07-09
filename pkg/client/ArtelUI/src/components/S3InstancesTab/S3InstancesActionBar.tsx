import {Button} from "@vervstack/chures"

import cls from "@/components/S3InstancesTab/S3InstancesActionBar.module.css"

export default function S3InstancesActionBar({count, onAddClick}: {count: number; onAddClick: () => void}) {
    return (
        <div className={cls.S3InstancesActionBarContainer}>
            <p className={cls.HeroSub}>
                <b>{count} {count === 1 ? "instance" : "instances"}</b>
                {" · "}<span>manage S3-compatible storage backends</span>
            </p>
            <Button variant="primary" onClick={onAddClick}>Add instance</Button>
        </div>
    )
}
