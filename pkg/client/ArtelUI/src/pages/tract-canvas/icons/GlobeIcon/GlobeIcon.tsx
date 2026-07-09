import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function GlobeIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <circle cx="12" cy="12" r="10"/>
            <line x1="2" y1="12" x2="22" y2="12"/>
            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
        </svg>
    )
}
