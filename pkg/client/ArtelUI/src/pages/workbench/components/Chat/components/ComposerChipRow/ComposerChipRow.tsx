import {Button} from "@vervstack/chures"

import type {TweaksSection} from "@/pages/workbench/processes/workbenchContext.ts"
import cls from "@/pages/workbench/components/Chat/components/ComposerChipRow/ComposerChipRow.module.css"

interface Props {
    onOpenTweaks: (section?: TweaksSection) => void
}

// Quick-access chips below the composer — each opens the Tweaks panel focused on
// its section (System prompt / Max tokens / Context window / Connections).
const CHIPS: {label: string; section: TweaksSection}[] = [
    {label: "System prompt", section: "system"},
    {label: "Max tokens", section: "tokens"},
    {label: "Context", section: "context"},
    {label: "Connections", section: "connections"},
]

export default function ComposerChipRow({onOpenTweaks}: Props) {
    return (
        <div className={cls.ComposerChipRowContainer}>
            {CHIPS.map(chip => (
                <Button
                    key={chip.section}
                    variant="unstyled"
                    className={cls.Chip}
                    onClick={() => onOpenTweaks(chip.section)}
                >
                    {chip.label}
                </Button>
            ))}
        </div>
    )
}
