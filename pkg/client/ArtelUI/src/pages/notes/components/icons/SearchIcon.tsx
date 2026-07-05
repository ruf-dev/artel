interface SearchIconProps {
    className?: string
}

export default function SearchIcon({className}: SearchIconProps) {
    return (
        <svg viewBox="0 0 16 16" width={17} height={17} fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" className={className}>
            <circle cx="6.5" cy="6.5" r="4.5" />
            <path d="M10.5 10.5l3 3" />
        </svg>
    )
}
