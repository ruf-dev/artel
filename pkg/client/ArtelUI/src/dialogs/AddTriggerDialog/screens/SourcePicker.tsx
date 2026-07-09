import {useState} from "react"

import {Dropdown} from "@vervstack/chures"
import cls from "@/dialogs/AddTriggerDialog/AddTriggerDialog.module.css"

import {TriggerSource} from "@/processes/Tracts.ts"
import {categoryLabel} from "@/components/TriggerPanel/triggerLabels.ts"

// Two-level webhook source picker: pick a category first, then (if the category has more
// than one preset) pick a preset within it. Categories/presets with a single option resolve
// immediately without a second pick.
export default function SourcePicker({triggerSources, source, onSourceChange}: {
    triggerSources: TriggerSource[]
    source: string
    onSourceChange: (source: string) => void
}) {
    const categories = Array.from(new Set(triggerSources.map(s => s.category)))
    const selected = triggerSources.find(s => s.key === source)
    const [pendingCategory, setPendingCategory] = useState<string | null>(null)
    const category = pendingCategory ?? selected?.category ?? ""
    const presetsInCategory = triggerSources.filter(s => s.category === category)

    function pickCategory(value: string[]) {
        const key = value[0] ?? ""
        const presets = triggerSources.filter(s => s.category === key)
        if (presets.length === 1) {
            setPendingCategory(null)
            onSourceChange(presets[0].key)
        } else {
            setPendingCategory(key)
            onSourceChange("")
        }
    }

    function pickPreset(value: string[]) {
        setPendingCategory(null)
        onSourceChange(value[0] ?? "")
    }

    return (
        <>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Category</span>
                <Dropdown
                    placeholder="Choose category…"
                    options={categories.map(c => ({id: c, name: categoryLabel(c)}))}
                    value={category ? [category] : []}
                    onChange={pickCategory}
                />
            </label>
            {category && presetsInCategory.length > 1 && (
                <label className={cls.Field}>
                    <span className={cls.FieldLabel}>Preset</span>
                    <Dropdown
                        placeholder="Choose preset…"
                        options={presetsInCategory.map(p => ({id: p.key, name: p.label}))}
                        value={source ? [source] : []}
                        onChange={pickPreset}
                    />
                </label>
            )}
        </>
    )
}
