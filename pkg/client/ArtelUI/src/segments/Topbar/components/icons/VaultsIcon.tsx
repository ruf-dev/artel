import {base, NavIconProps} from "@/segments/Topbar/components/icons/iconTypes.ts"

export default function VaultsIcon({className}: NavIconProps) {
    return (
        <svg className={className} {...base}>
            <ellipse cx="12" cy="5" rx="8" ry="3"/>
            <path d="M4 5v14a8 3 0 0 0 16 0V5"/>
            <path d="M4 12a8 3 0 0 0 16 0"/>
        </svg>
    )
}
