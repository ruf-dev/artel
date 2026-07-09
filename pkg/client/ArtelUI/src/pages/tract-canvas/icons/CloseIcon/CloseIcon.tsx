import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function CloseIcon({className}: IconProps) {
    return (
        <svg className={className} {...base} strokeWidth={1.8}>
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
        </svg>
    )
}
