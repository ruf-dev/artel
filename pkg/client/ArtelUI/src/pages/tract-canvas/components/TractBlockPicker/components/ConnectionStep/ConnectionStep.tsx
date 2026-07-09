import {useEffect} from "react"
import {Button} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"
import cls from "@/pages/tract-canvas/components/TractBlockPicker/components/ConnectionStep/ConnectionStep.module.css"

interface ConnectionStepProps {
    mcp: string
    selectedConnectionId: string
    onSelect: (id: string) => void
    onBack: () => void
    onConfirm: () => void
}

export default function ConnectionStep({mcp, selectedConnectionId, onSelect, onBack, onConfirm}: ConnectionStepProps) {
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
        <div className={cls.ConnectionStepContainer} role="dialog" aria-modal="true">
            <div className={cls.Head}>
                <div className={cls.Title}>Choose a connection</div>
                <Button variant="ghost" onClick={CloseDialog} aria-label="Close">✕</Button>
            </div>
            <div className={cls.Scroll}>
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
            </div>
            <div className={cls.Actions}>
                <Button variant="ghost" onClick={onBack}>Back</Button>
                <Button variant="primary" disabled={!selectedConnectionId} onClick={onConfirm}>Add</Button>
            </div>
        </div>
    )
}
