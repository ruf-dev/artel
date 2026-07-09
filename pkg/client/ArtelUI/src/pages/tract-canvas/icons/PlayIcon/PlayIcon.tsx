import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function PlayIcon({className}: IconProps) {
    return (
        <svg className={className} {...base} strokeWidth={2}>
            <polygon points="5 3 19 12 5 21 5 3"/>
        </svg>
    )
}
