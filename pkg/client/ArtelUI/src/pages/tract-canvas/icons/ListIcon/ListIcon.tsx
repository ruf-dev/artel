import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function ListIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <line x1="3" y1="6" x2="21" y2="6"/>
            <path d="M3 12h9"/>
            <line x1="3" y1="18" x2="21" y2="18"/>
        </svg>
    )
}
