interface ArtelMarkProps {
    size?: number
    maskId: string
}

export default function ArtelMark({ size = 11, maskId }: ArtelMarkProps) {
    return (
        <svg viewBox="0 0 100 100" width={size} height={size} style={{ flexShrink: 0, display: "block" }}>
            <defs>
                <mask id={maskId}>
                    <rect width="100" height="100" fill="white" />
                    <path d="M50 22L28 78h10l5.5-14h13L62 78h10ZM46.5 56L50 38l3.5 18Z" fill="black" />
                </mask>
            </defs>
            <circle cx="50" cy="50" r="50" fill="#FF4B3E" mask={`url(#${maskId})`} />
        </svg>
    )
}
