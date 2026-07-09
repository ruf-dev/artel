import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function ScheduleIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <circle cx="12" cy="12" r="10"/>
            <polyline points="12 6 12 12 16 14"/>
        </svg>
    )
}
