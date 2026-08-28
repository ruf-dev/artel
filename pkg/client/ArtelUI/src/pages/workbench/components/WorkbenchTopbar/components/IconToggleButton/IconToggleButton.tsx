import {ReactNode} from "react"
import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import cls from "@/pages/workbench/components/WorkbenchTopbar/components/IconToggleButton/IconToggleButton.module.css"

interface Props {
    icon: ReactNode
    label: string
    onClick?: () => void
    active?: boolean
    disabled?: boolean
    tooltip?: string
}

// Square transparent icon button styled after the Claude Design mock's `.icon-btn`
// — the nav toggle and the (disabled) tweaks toggle in WorkbenchTopbar. Disabled
// uses aria-disabled (not the real DOM prop, which would swallow the tooltip's
// pointer events) plus a guarded onClick.
export default function IconToggleButton({icon, label, onClick, active, disabled, tooltip}: Props) {
    return (
        <div className={cls.IconToggleButtonContainer}>
            <Button
                variant="unstyled"
                className={cn(cls.Btn, active && cls.BtnActive, disabled && cls.BtnDisabled)}
                aria-label={label}
                aria-disabled={disabled || undefined}
                data-tooltip-id={tooltip ? "root-tooltip" : undefined}
                data-tooltip-content={tooltip}
                onClick={() => {
                    if (!disabled) onClick?.()
                }}
            >
                {icon}
            </Button>
        </div>
    )
}
