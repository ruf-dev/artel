interface ButtonConfig {
    label: string
    onClick: () => void
    className: string
    disabled?: boolean
}

interface Props {
    buttons: ButtonConfig[]
    containerClassName?: string
}

export default function ModalActions({buttons, containerClassName}: Props) {
    return (
        <div className={containerClassName}>
            {buttons.map((btn, idx) => (
                <button
                    key={idx}
                    className={btn.className}
                    type="button"
                    onClick={btn.onClick}
                    disabled={btn.disabled}
                >
                    {btn.label}
                </button>
            ))}
        </div>
    )
}