export default function ArtelLogoIcon() {
    return (
        <svg viewBox="0 0 100 100" width={22} height={22} style={{ flexShrink: 0 }}>
            <defs>
                <mask id="mobile-artel-mask">
                    <rect width="100" height="100" fill="white" />
                    <path d="M 50 22 L 28 78 L 38 78 L 43.5 64 L 56.5 64 L 62 78 L 72 78 Z M 46.5 56 L 50 38 L 53.5 56 Z" fill="black" />
                </mask>
            </defs>
            <circle cx="50" cy="50" r="50" fill="#FF4B3E" mask="url(#mobile-artel-mask)" />
        </svg>
    )
}
