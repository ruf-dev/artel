import {base, NavIconProps} from "@/segments/Topbar/components/icons/iconTypes.ts"

export default function TractTemplatesIcon({className}: NavIconProps) {
    return (
        <svg className={className} {...base}>
            <rect x="3.5" y="3.5" width="7" height="7" rx="1.5"/>
            <rect x="13.5" y="3.5" width="7" height="7" rx="1.5"/>
            <rect x="3.5" y="13.5" width="7" height="7" rx="1.5"/>
            <rect x="13.5" y="13.5" width="7" height="7" rx="1.5"/>
        </svg>
    )
}
