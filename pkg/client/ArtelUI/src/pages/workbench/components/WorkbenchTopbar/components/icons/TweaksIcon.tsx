// Sliders/tweaks glyph for the (not-yet-built) tweaks panel toggle — traced from
// the Claude Design mock's #tweaksToggle svg.
export default function TweaksIcon() {
    return (
        <svg
            width="17"
            height="17"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="1.6"
            strokeLinecap="round"
        >
            <path d="M4 6h16M4 12h16M4 18h16"/>
            <circle cx="9" cy="6" r="2" fill="var(--color-bg-card)"/>
            <circle cx="16" cy="12" r="2" fill="var(--color-bg-card)"/>
            <circle cx="10" cy="18" r="2" fill="var(--color-bg-card)"/>
        </svg>
    )
}
