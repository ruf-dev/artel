import {Button} from "@vervstack/chures"

import cls from "@/pages/workbench/components/Chat/components/ComposerChipRow/ComposerChipRow.module.css"

// Disabled placeholder chips below the composer. Real controls (system prompt,
// token cap, context budget, connection picker) get wired in Stages 5/7.
const CHIPS = ["System prompt", "Max tokens", "Context", "Connections"]

export default function ComposerChipRow() {
    return (
        <div className={cls.ComposerChipRowContainer}>
            {CHIPS.map(label => (
                <Button
                    key={label}
                    variant="unstyled"
                    className={cls.Chip}
                    aria-disabled="true"
                    data-tooltip-id="root-tooltip"
                    data-tooltip-content="Coming soon"
                >
                    {label}
                </Button>
            ))}
        </div>
    )
}
