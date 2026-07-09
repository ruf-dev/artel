import {Button} from "@vervstack/chures"

import cls from "@/dialogs/AddTriggerDialog/AddTriggerDialog.module.css"
import {TriggerSource} from "@/processes/Tracts.ts"
import {providerLabel} from "@/components/TriggerPanel/triggerLabels.ts"

// Read-only preview of a resolved preset: description + its payload schema fields (no inputs),
// plus a provider-connection nudge for provider-linked presets (e.g. gitlab_push).
export default function PresetDetails({preset, hasProviderConnection, onConnectProvider}: {
    preset: TriggerSource
    hasProviderConnection: boolean
    onConnectProvider: () => void
}) {
    const propertyEntries = Object.entries(preset.payloadSchema.properties)

    return (
        <div className={cls.Field}>
            <span className={cls.FieldLabel}>{preset.label}</span>
            <p className={cls.Notice}>{preset.description}</p>
            {propertyEntries.length > 0 && (
                <div className={cls.SchemaPreview}>
                    {propertyEntries.map(([propName, prop]) => (
                        <div className={cls.SchemaFieldRow} key={propName}>
                            <span className={cls.SchemaFieldName}>{propName}</span>
                            <span className={cls.SchemaFieldTypeSelect}>{prop.type}</span>
                            {prop.description && <span className={cls.Notice}>{prop.description}</span>}
                        </div>
                    ))}
                </div>
            )}
            {preset.provider && (
                hasProviderConnection ? (
                    <p className={cls.Notice}>Uses your connected {providerLabel(preset.provider)}.</p>
                ) : (
                    <div className={cls.ProviderNotice}>
                        <p className={cls.Notice}>Connect {providerLabel(preset.provider)} first to use this trigger.</p>
                        <Button variant="secondary" onClick={onConnectProvider}>Connect {providerLabel(preset.provider)}</Button>
                    </div>
                )
            )}
        </div>
    )
}
