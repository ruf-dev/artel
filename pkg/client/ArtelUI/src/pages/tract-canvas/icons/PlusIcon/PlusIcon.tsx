import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function PlusIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
    )
}
