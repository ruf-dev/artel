import {useState} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/components/WebhookPicker/WebhookPicker.module.css"
import {categoryLabel} from "@/components/TriggerPanel/triggerLabels.ts"
import {useTriggerSources} from "@/app/hooks/Tracts.ts"

interface WebhookPickerProps {
    source: string
    onSourceChange: (source: string) => void
}

// Two-level webhook source picker: picking a category shows its preset list, and the
// schema for a preset is only resolved (by the caller, via `source`) once a preset is picked.
// Exception: a category whose only preset shares its key (e.g. "generic") IS the webhook —
// there's no sub-choice to make, so picking that category selects it immediately.
export default function WebhookPicker({source, onSourceChange}:WebhookPickerProps)
{
    const {triggerSources} = useTriggerSources()
    const categories = Array.from(new Set(triggerSources.map(s => s.category)))

    const selected = triggerSources.find(s => s.key === source)
    const [pendingCategory, setPendingCategory] = useState<string | null>(null)
    const category = pendingCategory ?? selected?.category ?? ""
    const presetsInCategory = triggerSources.filter(s => s.category === category)
    const isSelfPreset = presetsInCategory.length === 1 && presetsInCategory[0].key === category

    function pickCategory(key: string) {
        const presets = triggerSources.filter(s => s.category === key)
        if (presets.length === 1 && presets[0].key === key) {
            setPendingCategory(null)
            onSourceChange(presets[0].key)
        } else {
            setPendingCategory(key)
            onSourceChange("")
        }
    }

    function pickPreset(key: string) {
        setPendingCategory(null)
        onSourceChange(key)
    }

    return (
        <div className={cls.WebhookPickerContainer}>
            <span className={cls.FieldLabel}>Type</span>
            <div className={cls.SourceCategoryRow}>
                {categories.map(c => (
                    <Button
                        key={c}
                        type="button"
                        variant={c === category ? "primary" : "secondary"}
                        onClick={() => pickCategory(c)}
                    >
                        {categoryLabel(c)}
                    </Button>
                ))}
            </div>
            {presetsInCategory.length > 0 && !isSelfPreset && (
                <div className={cls.SourcePresetList}>
                    {presetsInCategory.map(p => (
                        <Button
                            key={p.key}
                            type="button"
                            variant={p.key === source ? "primary" : "secondary"}
                            onClick={() => pickPreset(p.key)}
                        >
                            {p.label}
                        </Button>
                    ))}
                </div>
            )}
        </div>
    )
}
