export default function SidebarGearIcon({className}: { className?: string }) {
    return (
        <svg className={className} width="16" height="16" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <circle cx="12" cy="12" r="3" stroke="currentColor" strokeWidth="1.6"/>
            <path
                d="M12 2v3M12 17v3M22 12h-3M5 12H2M19 5l-2.122 2.122M7.122 16.878L5 19M19 19l-2.122-2.122M7.122 7.122L5 5"
                stroke="currentColor"
                strokeWidth="1.6"
                strokeLinecap="round"
                strokeLinejoin="round"
            />
        </svg>
    )
}
