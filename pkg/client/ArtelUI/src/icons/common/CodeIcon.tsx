export default function CodeIcon({className}: { className?: string }) {
    return (
        <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={1.6} strokeLinecap="round" strokeLinejoin="round">
            <polyline points="8 6 2 12 8 18"/>
            <polyline points="16 6 22 12 16 18"/>
        </svg>
    )
}
