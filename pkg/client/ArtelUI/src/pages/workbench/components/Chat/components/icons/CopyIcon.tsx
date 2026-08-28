export default function CopyIcon({className}: {className?: string}) {
    return (
        <svg className={className} viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor"
             strokeWidth="1.7" strokeLinecap="round" strokeLinejoin="round">
            <rect x="9" y="9" width="11" height="11" rx="2"/>
            <path d="M5 15V5a2 2 0 0 1 2-2h8"/>
        </svg>
    )
}
