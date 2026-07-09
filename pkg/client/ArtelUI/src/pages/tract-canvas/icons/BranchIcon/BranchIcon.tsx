import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function BranchIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <polyline points="22 7 13.5 15.5 8.5 10.5 2 17"/>
            <polyline points="16 7 22 7 22 13"/>
        </svg>
    )
}
