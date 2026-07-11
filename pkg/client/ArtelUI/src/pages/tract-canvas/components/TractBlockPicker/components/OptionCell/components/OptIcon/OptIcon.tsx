import {StepColor} from "@/pages/tract-canvas/icons/tractIconHelpers.ts"
import {IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"
import cls from "@/pages/tract-canvas/components/TractBlockPicker/components/OptionCell/components/OptIcon/OptIcon.module.css"

interface OptIconProps {
    Icon: (props: IconProps) => React.JSX.Element
    color: StepColor
}

export default function OptIcon({Icon, color}: OptIconProps) {
    return (
        <div className={cls.OptIconContainer} style={{background: color.bg, borderColor: color.border, color: color.fg}}>
            <Icon/>
        </div>
    )
}
