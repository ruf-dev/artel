import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function ChevronRightIcon({className}: IconProps) {
    return (
        <svg className={className} {...base} strokeWidth={2}>
            <polyline points="9 18 15 12 9 6"/>
        </svg>
    )
}
