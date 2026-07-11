import {useState, useEffect} from "react"
import {ConfirmDialog} from "@vervstack/chures"

import {McpKeyInfo, McpConnectorInfo, MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"

export type ManageStep = "main" | "vault" | "addConnection" | "selectConnection"

export function useManageKeyDialog(mcpKey: McpKeyInfo) {
    const [step, setStep] = useState<ManageStep>("main")
    const [saving, setSaving] = useState(false)
    const [selectedVaultId, setSelectedVaultId] = useState(mcpKey.vaultId ?? "")
    const [selectedCandidate, setSelectedCandidate] = useState<MomCandidate | null>(null)
    const [selectedExternalConnectionId, setSelectedExternalConnectionId] = useState("")
    const [editingConnectorName, setEditingConnectorName] = useState<string | null>(null)

    const mcpKeysStore = useMcpKeys()
    const {OpenDialog, CloseDialog} = useDialog()

    const connectors = mcpKey.id ? mcpKeysStore.connectorsByKey[mcpKey.id] ?? [] : []

    useEffect(() => {
        if (mcpKey.id) {
            void mcpKeysStore.fetchConnectors(mcpKey.id)
        }
    }, [mcpKey.id, mcpKeysStore.fetchConnectors])

    useEffect(() => {
        if (step === "addConnection") {
            void mcpKeysStore.fetchMomCandidates()
        }
    }, [step, mcpKeysStore.fetchMomCandidates])

    async function handleSaveVault() {
        if (!mcpKey.id) return
        setSaving(true)
        try {
            await mcpKeysStore.setAccess(mcpKey.id, selectedVaultId)
            setStep("main")
        } finally {
            setSaving(false)
        }
    }

    function handleSelectCandidate(candidate: MomCandidate) {
        setSelectedCandidate(candidate)
        setSelectedExternalConnectionId("")
        setStep("selectConnection")
    }

    async function handleEditConnector(connector: McpConnectorInfo) {
        setEditingConnectorName(connector.mcpName ?? "")
        setSelectedExternalConnectionId(connector.externalConnectionId ?? "")
        await mcpKeysStore.fetchMomCandidates()
        const {momCandidates: fresh} = useMcpKeys.getState()
        const candidate = fresh.find(c => c.name === connector.mcpName)
            ?? {name: connector.mcpName ?? "", connections: []}
        setSelectedCandidate(candidate as MomCandidate)
        setStep("selectConnection")
    }

    async function handleAddConnector() {
        if (!mcpKey.id || !selectedCandidate?.name || !selectedExternalConnectionId) return
        setSaving(true)
        try {
            if (editingConnectorName) {
                await mcpKeysStore.removeConnector(mcpKey.id, editingConnectorName)
                setEditingConnectorName(null)
            }
            await mcpKeysStore.addConnector(mcpKey.id, selectedCandidate.name, selectedExternalConnectionId)
            setSelectedCandidate(null)
            setSelectedExternalConnectionId("")
            setStep("main")
        } finally {
            setSaving(false)
        }
    }

    function handleRevoke() {
        OpenDialog(
            <ConfirmDialog
                title="Revoke key"
                message={
                    `Revoke "${mcpKey.name}"? Any agents using this key will immediately lose access. `
                    + "This cannot be undone."
                }
                confirmLabel="Revoke"
                danger
                onClose={CloseDialog}
                onConfirm={() => {
                    if (!mcpKey.id || !mcpKey.vaultId) return
                    return mcpKeysStore.revoke(mcpKey.id, mcpKey.vaultId)
                }}
            />
        )
    }

    return {
        step, setStep,
        saving,
        selectedVaultId, setSelectedVaultId,
        selectedCandidate,
        selectedExternalConnectionId, setSelectedExternalConnectionId,
        editingConnectorName,
        connectors,
        handleSaveVault,
        handleSelectCandidate,
        handleEditConnector,
        handleAddConnector,
        handleRevoke,
    }
}
