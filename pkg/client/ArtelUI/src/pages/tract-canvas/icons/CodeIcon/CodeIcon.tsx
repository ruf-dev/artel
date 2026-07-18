import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function CodeIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <polyline points="16 18 22 12 16 6"/>
            <polyline points="8 6 2 12 8 18"/>
        </svg>
    )
}
