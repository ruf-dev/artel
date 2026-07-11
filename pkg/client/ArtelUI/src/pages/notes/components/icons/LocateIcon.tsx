interface LocateIconProps {
    className?: string
}

export default function LocateIcon({className}: LocateIconProps) {
    return (
        <svg viewBox="0 0 16 16" width={17} height={17} fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" className={className}>
            <circle cx="8" cy="8" r="2.5" />
            <path d="M8 1v2.4M8 12.6V15M1 8h2.4M12.6 8H15" />
        </svg>
    )
}
