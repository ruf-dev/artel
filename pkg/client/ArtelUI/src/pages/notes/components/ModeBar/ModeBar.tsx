import cls from "./ModeBar.module.css"

type Mode = 'edit' | 'preview' | 'read'

interface ModeBarProps {
    active: Mode
    onModeChange: (mode: Mode) => void
}

const MODES: { key: Mode; label: string }[] = [
    { key: 'edit', label: 'Edit' },
    { key: 'preview', label: 'Preview' },
    { key: 'read', label: 'Read' },
]

export default function ModeBar({ active, onModeChange }: ModeBarProps) {
    return (
        <div className={cls.ModeBarContainer}>
            {MODES.map(({ key, label }) => (
                <button
                    key={key}
                    className={`${cls.Tab}${active === key ? ` ${cls.TabActive}` : ""}`}
                    onClick={() => onModeChange(key)}
                >
                    {label}
                </button>
            ))}
        </div>
    )
}
