import {base, NavIconProps} from "@/segments/Topbar/components/icons/iconTypes.ts"

export default function ConnectionsIcon({className}: NavIconProps) {
    return (
        <svg className={className} {...base}>
            <path d="M15 7h3a5 5 0 0 1 0 10h-3"/>
            <path d="M9 17H6a5 5 0 0 1 0-10h3"/>
            <path d="M8 12h8"/>
        </svg>
    )
}
