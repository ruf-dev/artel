import {colorForKind} from "@/pages/tract-canvas/icons/tractIconHelpers.ts"
import {LogicOption} from "@/pages/tract-canvas/components/TractBlockPicker/useTractBlockPickerData.ts"
import OptionCell from "@/pages/tract-canvas/components/TractBlockPicker/components/OptionCell/OptionCell.tsx"

interface LogicCellProps {
    option: LogicOption
    onSelect: () => void
}

export default function LogicCell({option, onSelect}: LogicCellProps) {
    const color = colorForKind(option.type)
    return <OptionCell Icon={option.Icon} color={color} name={option.name} fn={option.desc} onSelect={onSelect}/>
}
