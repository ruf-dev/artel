import {Button} from "@vervstack/chures"

import cls
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/components/FeatureOverrideRow/FeatureOverrideRow.module.css"

export type FeatureOverrideState = "inherit" | "on" | "off"

interface FeatureOverrideRowProps {
    label: string
    state: FeatureOverrideState
    onChange: (state: FeatureOverrideState) => void
}

export default function FeatureOverrideRow({label, state, onChange}: FeatureOverrideRowProps) {
    return (
        <div className={cls.FeatureOverrideRowContainer}>
            <span className={cls.FeatureLabel}>{label}</span>
            <div className={cls.SegmentedControl}>
                <Button variant={state === "inherit" ? "primary" : "ghost"} onClick={() => onChange("inherit")}>
                    Inherit
                </Button>
                <Button variant={state === "on" ? "primary" : "ghost"} onClick={() => onChange("on")}>
                    On
                </Button>
                <Button variant={state === "off" ? "danger" : "ghost"} onClick={() => onChange("off")}>
                    Off
                </Button>
            </div>
        </div>
    )
}
