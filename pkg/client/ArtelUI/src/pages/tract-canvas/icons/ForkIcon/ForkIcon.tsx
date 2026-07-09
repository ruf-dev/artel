import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function ForkIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <circle cx="12" cy="18" r="2.5"/>
            <circle cx="6" cy="6" r="2.5"/>
            <circle cx="18" cy="6" r="2.5"/>
            <path d="M6 8.5v1A2.5 2.5 0 0 0 8.5 12h7A2.5 2.5 0 0 0 18 9.5v-1"/>
            <path d="M12 12v3.5"/>
        </svg>
    )
}
