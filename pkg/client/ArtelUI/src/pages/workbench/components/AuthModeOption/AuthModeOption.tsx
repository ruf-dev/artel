import {Button} from "@vervstack/chures"

import cls from "@/pages/workbench/components/AuthModeOption/AuthModeOption.module.css"
import {cn} from "@/app/utils/cn.ts"

interface Props {
    selected: boolean
    label: string
    desc: string
    disabled?: boolean
    hint?: string
    onSelect: () => void
}

export default function AuthModeOption({selected, label, desc, disabled, hint, onSelect}: Props) {
    return (
        <Button
            className={cn(cls.AuthModeOptionContainer, selected && cls.AuthModeOptionSelected)}
            onClick={onSelect}
            disabled={disabled}
        >
            <span className={cls.AuthModeOptionLabel}>{label}</span>
            <span className={cls.AuthModeOptionDesc}>{desc}</span>
            {disabled && hint && <span className={cls.AuthModeOptionHint}>{hint}</span>}
        </Button>
    )
}
