import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn"
import cls from "@/components/OptionCard/OptionCard.module.css"

interface Props {
    selected: boolean
    label: string
    desc: string
    hint?: string
    disabled?: boolean
    onSelect: () => void
}

// Generalized selectable option card for picker screens.
// Reconciles WorkbenchModeOption and AuthModeOption into a single component
// with full styling support (hover, disabled, selected tint).
export default function OptionCard({selected, label, desc, hint, disabled, onSelect}: Props) {
    return (
        <Button
            className={cn(cls.OptionCardContainer, selected && cls.OptionCardSelected)}
            onClick={onSelect}
            disabled={disabled}
        >
            <span className={cls.OptionCardLabel}>{label}</span>
            <span className={cls.OptionCardDesc}>{desc}</span>
            {disabled && hint && <span className={cls.OptionCardHint}>{hint}</span>}
        </Button>
    )
}
