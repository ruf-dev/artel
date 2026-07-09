// TODO: placeholder glyph for providers without a dedicated brand icon yet - replace per-provider as icons are added.
export default function UnknownProviderIcon() {
    return (
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"
             stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
            <circle cx="12" cy="12" r="9"/>
            <path d="M9.5 9.5a2.5 2.5 0 1 1 3.5 2.3c-.7.35-1 .85-1 1.7v.3"/>
            <path d="M12 17h.01"/>
        </svg>
    )
}
