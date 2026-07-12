import {base, NavIconProps} from "@/segments/Topbar/components/icons/iconTypes.ts"

export default function ToolboxIcon({className}: NavIconProps) {
    return (
        <svg className={className} {...base}>
            <rect x="3" y="9" width="18" height="10" rx="2"/>
            <path d="M8 9V6a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v3"/>
            <path d="M3 13h18"/>
            <path d="M11 13v2"/>
            <path d="M13 13v2"/>
        </svg>
    )
}
