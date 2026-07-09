import {useEffect} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/components/StepPickerDialog/StepPickerDialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"

interface Props {
    mcp: string
    selectedConnectionId: string
    onSelect: (id: string) => void
    onBack: () => void
    onConfirm: () => void
}

export default function ConnectionStep({mcp, selectedConnectionId, onSelect, onBack, onConfirm}: Props) {
    const {CloseDialog} = useDialog()
    const {momCandidates, fetchMomCandidates} = useMcpKeys()

    useEffect(() => {
        void fetchMomCandidates()
    }, [fetchMomCandidates])

    const candidate = momCandidates.find(c => c.name === mcp)
    const connections = candidate?.connections ?? []

    useEffect(() => {
        if (!selectedConnectionId && connections.length > 0) {
            onSelect(connections[0].id ?? "")
        }
    }, [connections, selectedConnectionId, onSelect])

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Choose a connection</h2>
            <div className={cls.List}>
                {connections.map(c => (
                    <SelectOption
                        key={c.id}
                        label={connectionLabel(c)}
                        selected={c.id === selectedConnectionId}
                        onSelect={() => onSelect(c.id ?? "")}
                    />
                ))}
                {connections.length === 0 && <p className={cls.Empty}>No connections available for {mcp}.</p>}
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
