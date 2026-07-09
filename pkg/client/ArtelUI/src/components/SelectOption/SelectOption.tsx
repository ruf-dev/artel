import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn"
import cls from "@/components/SelectOption/SelectOption.module.css"

export default function SelectOption({label, selected, onSelect}: {
    label: string
    selected: boolean
    onSelect: () => void
}) {
    return (
        <Button
            variant="ghost"
            className={cn(cls.OptionRow, selected && cls.OptionRowSelected)}
            onClick={onSelect}
        >
            <span className={cls.OptionRadio}>{selected ? "●" : "○"}</span>
            <span>{label}</span>
        </Button>
    )
}
