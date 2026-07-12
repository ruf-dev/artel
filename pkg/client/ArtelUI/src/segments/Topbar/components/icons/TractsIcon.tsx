import {base, NavIconProps} from "@/segments/Topbar/components/icons/iconTypes.ts"

export default function TractsIcon({className}: NavIconProps) {
    return (
        <svg className={className} {...base}>
            <line x1="6" y1="3" x2="6" y2="15"/>
            <circle cx="18" cy="6" r="2.5"/>
            <circle cx="6" cy="18" r="2.5"/>
            <path d="M18 8.5a7.5 7.5 0 0 1-7.5 7.5"/>
        </svg>
    )
}
