import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function ManualTriggerIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <polygon points="6 3 20 12 6 21 6 3"/>
        </svg>
    )
}
