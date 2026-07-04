// Shared line-icon set for the tract canvas UI (list page, canvas nodes, block picker) — keeps
// icon choices consistent across all three instead of each widget inventing its own SVGs.

export interface IconProps {
    className?: string
}

const base = {
    viewBox: "0 0 24 24",
    fill: "none",
    stroke: "currentColor",
    strokeWidth: 1.6,
    strokeLinecap: "round" as const,
    strokeLinejoin: "round" as const,
}

export function VaultIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <path d="M3 7l9-4 9 4-9 4-9-4z"/>
            <path d="M3 7v10l9 4 9-4V7"/>
        </svg>
    )
}

export function ManualTriggerIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <polygon points="6 3 20 12 6 21 6 3"/>
        </svg>
    )
}

export function WebhookIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <path d="M18 16.98h-5.99c-1.1 0-1.95.94-2.48 1.9A4 4 0 0 1 2 17c0-1.1.94-2 2.06-2H8"/>
            <path d="M17 7.82a2 2 0 0 1 3 1.73V13"/>
            <path d="M13 6.51a2 2 0 0 0-3.46 0L6 12.5"/>
        </svg>
    )
}

export function ScheduleIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <circle cx="12" cy="12" r="10"/>
            <polyline points="12 6 12 12 16 14"/>
        </svg>
    )
}

export function SearchIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <circle cx="11" cy="11" r="8"/>
            <line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
    )
}

export function EditIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <path d="M12 20h9"/>
            <path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4z"/>
        </svg>
    )
}

export function PlusIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
    )
}

export function TrashIcon({className}: IconProps) {
    return (
        <svg className={className} {...base} strokeWidth={2}>
            <polyline points="3 6 5 6 21 6"/>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
        </svg>
    )
}

export function ListIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <line x1="3" y1="6" x2="21" y2="6"/>
            <path d="M3 12h9"/>
            <line x1="3" y1="18" x2="21" y2="18"/>
        </svg>
    )
}

export function EnvelopeIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
            <polyline points="22,6 12,13 2,6"/>
        </svg>
    )
}

export function PaperPlaneIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <line x1="22" y1="2" x2="11" y2="13"/>
            <polygon points="22 2 15 22 11 13 2 9 22 2"/>
        </svg>
    )
}

export function ChatIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
        </svg>
    )
}

export function GlobeIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <circle cx="12" cy="12" r="10"/>
            <line x1="2" y1="12" x2="22" y2="12"/>
            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
        </svg>
    )
}

export function BranchIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <polyline points="22 7 13.5 15.5 8.5 10.5 2 17"/>
            <polyline points="16 7 22 7 22 13"/>
        </svg>
    )
}

export function ForkIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <circle cx="12" cy="18" r="2.5"/>
            <circle cx="6" cy="6" r="2.5"/>
            <circle cx="18" cy="6" r="2.5"/>
            <path d="M6 8.5v1A2.5 2.5 0 0 0 8.5 12h7A2.5 2.5 0 0 0 18 9.5v-1"/>
            <path d="M12 12v3.5"/>
        </svg>
    )
}

export function LayersIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <polygon points="12 2 22 8.5 12 15 2 8.5 12 2"/>
            <polyline points="2 15.5 12 22 22 15.5"/>
        </svg>
    )
}

export function WrenchIcon({className}: IconProps) {
    return (
        <svg className={className} {...base}>
            <path d="M14.7 6.3a4 4 0 0 0-5.4 5.4l-6 6a2 2 0 0 0 2.8 2.8l6-6a4 4 0 0 0 5.4-5.4l-2.4 2.4-2-2z"/>
        </svg>
    )
}

export function ChevronRightIcon({className}: IconProps) {
    return (
        <svg className={className} {...base} strokeWidth={2}>
            <polyline points="9 18 15 12 9 6"/>
        </svg>
    )
}

export function CloseIcon({className}: IconProps) {
    return (
        <svg className={className} {...base} strokeWidth={1.8}>
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
        </svg>
    )
}

export function PlayIcon({className}: IconProps) {
    return (
        <svg className={className} {...base} strokeWidth={2}>
            <polygon points="5 3 19 12 5 21 5 3"/>
        </svg>
    )
}

