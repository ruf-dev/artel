export default function HistoryIcon({className}: { className?: string }) {
    return (
        <svg className={className} viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
             strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
            <circle cx="12" cy="12" r="10"/>
            <polyline points="12 6 12 12 16 14"/>
        </svg>
    )
}
