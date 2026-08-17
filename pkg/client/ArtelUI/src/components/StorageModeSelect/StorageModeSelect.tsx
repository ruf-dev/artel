import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import cls from "@/components/StorageModeSelect/StorageModeSelect.module.css"

const OPTIONS: { value: string, label: string, help: string }[] = [
    {value: "none", label: "None", help: "Ephemeral — nothing is saved between runs."},
    {value: "freeform_notes", label: "Freeform notes", help: "Jot notes in your vault as the skill runs."},
    {
        value: "structured",
        label: "Structured",
        help: "Track structured records in Postgres — table set up separately.",
    },
]

interface Props {
    value: string
    onChange: (v: string) => void
    disabled?: boolean
}

export default function StorageModeSelect({value, onChange, disabled}: Props) {
    return (
        <div className={cls.StorageModeSelectContainer}>
            {OPTIONS.map(opt => (
                <Button
                    key={opt.value}
                    variant="ghost"
                    className={cn(cls.Option, value === opt.value && cls.OptionSelected)}
                    onClick={() => onChange(opt.value)}
                    disabled={disabled}
                >
                    <span className={cls.OptionRadio}>{value === opt.value ? "●" : "○"}</span>
                    <span className={cls.OptionBody}>
                        <span className={cls.OptionLabel}>{opt.label}</span>
                        <span className={cls.OptionHelp}>{opt.help}</span>
                    </span>
                </Button>
            ))}
        </div>
    )
}
