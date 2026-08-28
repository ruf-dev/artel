import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import cls
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksConnectionsSection/components/ConnectionRow/ConnectionRow.module.css"

interface Props {
    name: string
    detail?: string
    online?: boolean
    actionLabel?: string
    onAction?: () => void
}

// One line in the Tweaks panel's Connections section — status dot, a name, an
// optional monospace detail (masked key / container state), and an optional
// secondary action button. Mirrors the mock's `.conn-row`.
export default function ConnectionRow({name, detail, online, actionLabel, onAction}: Props) {
    return (
        <div className={cls.ConnectionRowContainer}>
            <span className={cn(cls.Dot, !online && cls.DotOff)}/>
            <span className={cls.Name}>{name}</span>
            {detail && <span className={cls.Detail}>{detail}</span>}
            {actionLabel && (
                <Button className={cls.Action} variant="secondary" onClick={onAction}>
                    {actionLabel}
                </Button>
            )}
        </div>
    )
}
