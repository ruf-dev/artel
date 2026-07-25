import {useEffect} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/components/StepPickerDialog/StepPickerDialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"

interface Props {
    selectedConnectionId: string
    onSelect: (id: string) => void
    onBack: () => void
    onConfirm: () => void
}

export default function LlmConnectionStep({selectedConnectionId, onSelect, onBack, onConfirm}: Props) {
    const {CloseDialog} = useDialog()
    const {connections, fetch: fetchConnections} = useExternalConnections()

    useEffect(() => {
        void fetchConnections()
    }, [fetchConnections])

    const anthropicConnections = connections.filter(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_ANTHROPIC)

    useEffect(() => {
        if (!selectedConnectionId && anthropicConnections.length > 0) {
            onSelect(anthropicConnections[0].id ?? "")
        }
    }, [anthropicConnections, selectedConnectionId, onSelect])

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Choose a connection</h2>
            <div className={cls.List}>
                {anthropicConnections.map(c => (
                    <SelectOption
                        key={c.id}
                        label={connectionLabel(c)}
                        selected={c.id === selectedConnectionId}
                        onSelect={() => onSelect(c.id ?? "")}
                    />
                ))}
                {anthropicConnections.length === 0 && (
                    <p className={cls.Empty}>No Anthropic connections available.</p>
                )}
            </div>
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={onBack}>Back</Button>
                <div className={cls.ActionsRight}>
                    <Button variant="ghost" onClick={CloseDialog}>Cancel</Button>
                    <Button variant="primary" disabled={!selectedConnectionId} onClick={onConfirm}>Add</Button>
                </div>
            </div>
        </div>
    )
}
