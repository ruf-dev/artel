import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function VaultIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <path d="M3 7l9-4 9 4-9 4-9-4z"/>
            <path d="M3 7v10l9 4 9-4V7"/>
        </svg>
    )
}
