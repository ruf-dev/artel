import {Button} from "@vervstack/chures"

import {StepColor} from "@/pages/tract-canvas/icons/tractIconHelpers.ts"
import {IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"
import OptIcon
    from "@/pages/tract-canvas/components/TractBlockPicker/components/OptionCell/components/OptIcon/OptIcon.tsx"
import OptText
    from "@/pages/tract-canvas/components/TractBlockPicker/components/OptionCell/components/OptText/OptText.tsx"
import cls from "@/pages/tract-canvas/components/TractBlockPicker/components/OptionCell/OptionCell.module.css"

interface OptionCellProps {
    Icon: (props: IconProps) => React.JSX.Element
    color: StepColor
    name: string
    fn: string
    onSelect: () => void
}

export default function OptionCell({Icon, color, name, fn, onSelect}: OptionCellProps) {
    return (
        <div className={cls.OptionCellContainer}>
            <Button variant="ghost" className={cls.OptionBtn} onClick={onSelect}>
                <OptIcon Icon={Icon} color={color}/>
                <OptText name={name} fn={fn}/>
            </Button>
        </div>
    )
}
