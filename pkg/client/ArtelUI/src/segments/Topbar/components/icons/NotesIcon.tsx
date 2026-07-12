import {base, NavIconProps} from "@/segments/Topbar/components/icons/iconTypes.ts"

export default function NotesIcon({className}: NavIconProps) {
    return (
        <svg className={className} {...base}>
            <path d="M14 3H7a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V8z"/>
            <path d="M14 3v5h5"/>
            <path d="M9 13h6"/>
            <path d="M9 17h6"/>
        </svg>
    )
}
