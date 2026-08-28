// TODO: chures has no segmented-control / button-group primitive yet, drop this wrapper once it does
import {Button} from "@vervstack/chures"

import cls from "@/components/atoms/SegmentedControl/SegmentedControl.module.css"
import {cn} from "@/app/utils/cn.ts"

interface SegmentedControlProps {
    options: {key: string; label: string; icon?: React.ReactNode; disabled?: boolean; tooltip?: string}[]
    active: string
    onChange: (key: string) => void
    collapsed?: boolean
}

export default function SegmentedControl({options, active, onChange, collapsed}: SegmentedControlProps) {
    return (
        <div className={cn(cls.SegmentedControlContainer, collapsed && cls.Collapsed)}>
            {options.map(({key, label, icon, disabled, tooltip}) => (
                <Button
                    key={key}
                    variant="unstyled"
                    aria-disabled={disabled || undefined}
                    data-tooltip-id={disabled && tooltip ? "root-tooltip" : undefined}
                    data-tooltip-content={disabled && tooltip ? tooltip : undefined}
                    className={cn(cls.Segment, active === key && cls.SegmentActive, disabled && cls.SegmentDisabled)}
                    onClick={() => {
                        if (!disabled) onChange(key)
                    }}
                >
                    {icon && <span className={cls.SegmentIcon}>{icon}</span>}
                    <span className={cls.SegmentLabel}>{label}</span>
                </Button>
            ))}
        </div>
    )
}
