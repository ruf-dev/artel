import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function LayersIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <polygon points="12 2 22 8.5 12 15 2 8.5 12 2"/>
            <polyline points="2 15.5 12 22 22 15.5"/>
        </svg>
    )
}
