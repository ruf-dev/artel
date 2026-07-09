import {Button} from "@vervstack/chures"

import {colorForKind} from "@/pages/tract-canvas/icons/tractIconHelpers.ts"
import {LogicOption} from "@/pages/tract-canvas/components/TractBlockPicker/useTractBlockPickerData.ts"
import cls from "@/pages/tract-canvas/components/TractBlockPicker/components/LogicCell/LogicCell.module.css"

interface LogicCellProps {
    option: LogicOption
    onSelect: () => void
}

export default function LogicCell({option, onSelect}: LogicCellProps) {
    const color = colorForKind(option.type)
    return (
        <Button variant="ghost" className={cls.Opt} onClick={onSelect}>
            <div className={cls.OptIcon} style={{background: color.bg, borderColor: color.border, color: color.fg}}>
                <option.Icon/>
            </div>
            <div className={cls.OptText}>
                <div className={cls.OptName}>{option.name}</div>
                <div className={cls.OptFn}>{option.desc}</div>
            </div>
        </Button>
    )
}
