import {base, IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"

export function WebhookIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <path d="M18 16.98h-5.99c-1.1 0-1.95.94-2.48 1.9A4 4 0 0 1 2 17c0-1.1.94-2 2.06-2H8"/>
            <path d="M17 7.82a2 2 0 0 1 3 1.73V13"/>
            <path d="M13 6.51a2 2 0 0 0-3.46 0L6 12.5"/>
        </svg>
    )
}
