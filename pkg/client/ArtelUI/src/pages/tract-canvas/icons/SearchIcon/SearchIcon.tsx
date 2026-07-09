import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function SearchIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <circle cx="11" cy="11" r="8"/>
            <line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
    )
}
