interface Props {
    onClick: () => void
    disabled?: boolean
    className?: string
}

export default function ModalClose({onClick, disabled, className}: Props) {
    return (
        <button
            className={className}
            type="button"
            onClick={onClick}
            disabled={disabled}
            aria-label="Close"
        >
            <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor"
                 strokeWidth="1.8" strokeLinecap="round">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
        </button>
    )
}