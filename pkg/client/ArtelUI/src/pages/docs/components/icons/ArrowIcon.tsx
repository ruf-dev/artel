interface ArrowIconProps {
    open: boolean
}

export default function ArrowIcon({ open }: ArrowIconProps) {
    return (
        <svg
            viewBox="0 0 8 8"
            width={11}
            height={11}
            style={{ flexShrink: 0, opacity: 0.35, transform: open ? "rotate(90deg)" : "none", transition: "transform 0.15s" }}
        >
            <path d="M2 1.5l4 2.5-4 2.5z" fill="currentColor" />
        </svg>
    )
}
