
import { Button } from "@vervstack/chures"

import { cn } from "@/app/utils/cn"
import cls from "@/pages/notes/components/ModeBar/ModeBar.module.css"

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
                <Button
                    key={key}
                    variant="unstyled"
                    className={cn(cls.Tab, active === key && cls.TabActive)}
                    onClick={() => onModeChange(key)}
                >
                    {label}
                </Button>
            ))}
        </div>
    )
}
