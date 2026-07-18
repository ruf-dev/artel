import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function CollapseIcon({className}: IconProps) {
    return (
        <svg className={className} {...base} strokeWidth={1.8}>
            <polyline points="9 15 9 21 3 21"/>
            <polyline points="15 9 15 3 21 3"/>
            <line x1="9" y1="15" x2="3" y2="21"/>
            <line x1="15" y1="9" x2="21" y2="3"/>
        </svg>
    )
}
