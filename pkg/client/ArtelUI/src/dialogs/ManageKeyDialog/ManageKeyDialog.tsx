import {useState, useEffect, type ReactNode} from "react"
import {useNavigate} from "react-router-dom"

import {Button, ModalClose, ConfirmDialog} from "@vervstack/chures"
import cls from "@/dialogs/ManageKeyDialog/ManageKeyDialog.module.css"

import {McpKeyInfo, McpConnectorInfo, MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useVaults} from "@/app/hooks/Vaults.ts"

import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import ConnectorChip from "@/components/ConnectorChip/ConnectorChip.tsx"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"

type ManageStep = "main" | "vault" | "addConnection" | "selectConnection"

export default function ManageKeyDialog({mcpKey}: { mcpKey: McpKeyInfo }) {
    const [step, setStep] = useState<ManageStep>("main")
    const [saving, setSaving] = useState(false)
    const [selectedVaultId, setSelectedVaultId] = useState(mcpKey.vaultId ?? "")
    const [selectedCandidate, setSelectedCandidate] = useState<MomCandidate | null>(null)
    const [selectedExternalConnectionId, setSelectedExternalConnectionId] = useState("")
    const [editingConnectorName, setEditingConnectorName] = useState<string | null>(null)

    const {
        setAccess,
        connectorsByKey,
        fetchConnectors,
        addConnector,
        removeConnector,
        fetchMomCandidates,
        revoke,
    } = useMcpKeys()
    const {OpenDialog, CloseDialog} = useDialog()

    const connectors = mcpKey.id ? connectorsByKey[mcpKey.id] ?? [] : []

    useEffect(() => {
        if (mcpKey.id) {
            void fetchConnectors(mcpKey.id)
        }
    }, [mcpKey.id, fetchConnectors])

    useEffect(() => {
        if (step === "addConnection") {
            void fetchMomCandidates()
        }
    }, [step, fetchMomCandidates])

    async function handleSaveVault() {
        if (!mcpKey.id) return
        setSaving(true)
        try {
            await setAccess(mcpKey.id, selectedVaultId)
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
        await fetchMomCandidates()
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
                await removeConnector(mcpKey.id, editingConnectorName)
                setEditingConnectorName(null)
            }
            await addConnector(mcpKey.id, selectedCandidate.name, selectedExternalConnectionId)
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
                message={`Revoke "${mcpKey.name}"? Any agents using this key will immediately lose access. This cannot be undone.`}
                confirmLabel="Revoke"
                danger
                onClose={CloseDialog}
                onConfirm={() => {
                    if (!mcpKey.id || !mcpKey.vaultId) return
                    return revoke(mcpKey.id, mcpKey.vaultId)
                }}
            />
        )
    }

    if (step === "vault") {
        return (
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                 aria-labelledby="manageVaultTitle">
                <DialogHead titleId="manageVaultTitle" title="Select vault" disabled={saving}/>
                <p className={cls.ModalSub}>Choose which vault this key connects to.</p>
                <VaultOptionList selectedVaultId={selectedVaultId} onSelect={setSelectedVaultId}/>
                <div className={cls.ModalActions}>
                    <Button variant="ghost" onClick={() => setStep("main")} disabled={saving}>Back</Button>
                    <Button variant="primary" onClick={handleSaveVault} disabled={saving || !selectedVaultId}>
                        {saving ? "Saving…" : "Save"}
                    </Button>
                </div>
            </div>
        )
    }

    if (step === "addConnection") {
        return (
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                 aria-labelledby="addConnectionTitle">
                <DialogHead titleId="addConnectionTitle" title="Add connection"/>
                <p className={cls.ModalSub}>Pick a service to connect to this key.</p>
                <CandidateOptionList connectors={connectors} onSelect={handleSelectCandidate}/>
                <div className={cls.ModalActions}>
                    <Button variant="ghost" onClick={() => setStep("main")}>Back</Button>
                </div>
            </div>
        )
    }

    if (step === "selectConnection" && selectedCandidate) {
        const available = selectedCandidate.connections ?? []

        return (
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                 aria-labelledby="selectConnectionTitle">
                <DialogHead titleId="selectConnectionTitle" title={selectedCandidate.name} disabled={saving}/>
                <p className={cls.ModalSub}>Pick the connection this key should use for it.</p>
                <ConnectionOptionList
                    available={available}
                    selectedId={selectedExternalConnectionId}
                    onSelect={setSelectedExternalConnectionId}
                />
                <div className={cls.ModalActions}>
                    <Button
                        variant="ghost"
                        onClick={() => setStep(editingConnectorName ? "main" : "addConnection")}
                        disabled={saving}>
                        Back
                    </Button>
                    <Button variant="primary" onClick={handleAddConnector} disabled={saving || !selectedExternalConnectionId}>
                        {saving ? (editingConnectorName ? "Saving…" : "Adding…") : (editingConnectorName ? "Save" : "Add")}
                    </Button>
                </div>
            </div>
        )
    }

    return (
        <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
             aria-labelledby="manageKeyTitle">
            <DialogHead
                titleId="manageKeyTitle"
                title={<>Manage mcp key: <b className={cls.KeyName}>{mcpKey.name}</b></>}
            />

            <VaultField
                selectedVaultId={selectedVaultId}
                onChangeVault={() => setStep("vault")}
            />

            <ConnectionsField
                connectors={connectors}
                mcpKeyId={mcpKey.id ?? ""}
                onEdit={handleEditConnector}
            />

            <div className={cls.ModalFooter}>
                <Button
                    variant="danger"
                    onClick={handleRevoke}>
                    Revoke key
                </Button>

                <Button variant="primary" onClick={() => setStep("addConnection")}>
                    Add connection
                </Button>
            </div>
        </div>
    )
}

function DialogHead({titleId, title, disabled}: {
    titleId: string
    title: ReactNode
    disabled?: boolean
}) {
    const {CloseDialog} = useDialog()
    return (
        <div className={cls.ModalHead}>
            <h2 className={cls.ModalTitle} id={titleId}>{title}</h2>
            <ModalClose onClick={CloseDialog} disabled={disabled} className={cls.ModalClose}/>
        </div>
    )
}

function VaultOptionList({selectedVaultId, onSelect}: {
    selectedVaultId: string
    onSelect: (id: string) => void
}) {
    const {vaults} = useVaults()
    return (
        <div className={cls.OptionList}>
            {vaults.map(v => (
                <SelectOption
                    key={v.id}
                    label={v.name ?? ""}
                    selected={selectedVaultId === v.id}
                    onSelect={() => onSelect(v.id ?? "")}
                />
            ))}
            {vaults.length === 0 && <p className={cls.Empty}>No vaults found.</p>}
        </div>
    )
}

function CandidateOptionList({connectors, onSelect}: {
    connectors: McpConnectorInfo[]
    onSelect: (c: MomCandidate) => void
}) {
    const {momCandidates: candidates} = useMcpKeys()
    return (
        <div className={cls.OptionList}>
            {candidates.map(c => (
                <MomCandidateCard
                    key={c.name}
                    candidate={c}
                    connected={connectors.some(con => con.mcpName === c.name)}
                    onSelect={() => onSelect(c)}
                />
            ))}
            {candidates.length === 0 && <p className={cls.Empty}>No services available yet.</p>}
        </div>
    )
}

function ConnectionOptionList({available, selectedId, onSelect}: {
    available: NonNullable<MomCandidate["connections"]>
    selectedId: string
    onSelect: (id: string) => void
}) {
    const {CloseDialog} = useDialog()
    const navigate = useNavigate()
    return (
        <div className={cls.OptionList}>
            {available.map(c => (
                <SelectOption
                    key={c.id}
                    label={connectionLabel(c)}
                    selected={selectedId === c.id}
                    onSelect={() => onSelect(c.id ?? "")}
                />
            ))}
            {available.length === 0 && (
                <p className={cls.Empty}>
                    No connections yet.{" "}
                    <button className={cls.LinkBtn} type="button"
                            onClick={() => {
                                CloseDialog();
                                navigate(Path.ConnectionsPage)
                            }}>Set one up
                    </button>
                </p>
            )}
        </div>
    )
}

function VaultChipDisplay({vault}: { vault: ReturnType<typeof useVaults>["vaults"][number] | undefined }) {
    return (
        <span
            className={vault ? cls.VaultChip : `${cls.Chip} ${cls.ChipMuted}`}
            data-tooltip-id="root-tooltip"
            data-tooltip-content={vault ? `Vault: ${vault.name}` : "No vault assigned"}
        >
            {vault && <span className={cls.VaultBadge}>A</span>}
            {vault ? vault.name : "No vault"}
        </span>
    )
}

function VaultField({selectedVaultId, onChangeVault}: {
    selectedVaultId: string
    onChangeVault: () => void
}) {
    const {vaults} = useVaults()
    const vault = vaults.find(v => v.id === selectedVaultId)
    return (
        <div className={cls.Field}>
            <span className={cls.FieldLabel}>Vault</span>
            <div className={cls.ConnectorRowWrapper}>
                <VaultChipDisplay vault={vault}/>
                <Button variant="ghost" onClick={onChangeVault}>Change</Button>
            </div>
        </div>
    )
}

function ConnectionsField({connectors, mcpKeyId, onEdit}: {
    connectors: McpConnectorInfo[]
    mcpKeyId: string
    onEdit: (connector: McpConnectorInfo) => void
}) {
    return (
        <div className={cls.Field}>
            <span className={cls.FieldLabel}>Connections</span>
            <div className={cls.OptionList}>
                {connectors.length == 0 ?
                    <p className={cls.Empty}>
                        No connections linked to this key yet.
                    </p>
                    :
                    connectors.map(c => (
                        <ConnectorRow
                            key={c.mcpName}
                            connector={c}
                            mcpKeyId={mcpKeyId}
                            onEdit={() => onEdit(c)}/>
                    ))}
            </div>
        </div>
    )
}

function MomCandidateCard(
    {candidate, connected, onSelect}: {
        candidate: MomCandidate
        connected: boolean
        onSelect: () => void
    }) {
    const count = candidate.connections?.length ?? 0
    return (
        <button
            className={cls.EntityCard}
            onClick={onSelect}
            disabled={connected}
            title={connected ? "Already connected" : undefined}
            type="button"
        >
            <div className={cls.MomCardHeader}>
                <span className={cls.EntityCardTitle}>{candidate.name}</span>
                <span className={count > 0 ? cls.Chip : `${cls.Chip} ${cls.ChipMuted}`}>
                    {connected ? "Already connected" : count > 0 ? `${count} connection${count > 1 ? "s" : ""}` : "No connection"}
                </span>
            </div>
            {candidate.description && <span className={cls.EntityCardDesc}>{candidate.description}</span>}
        </button>
    )
}

function ConnectorRow(
    {
        connector,
        mcpKeyId,
        onEdit,
    }: {
        connector: McpConnectorInfo;
        mcpKeyId: string;
        onEdit: () => void
    }) {
    const {connections} = useExternalConnections()
    const {removeConnector} = useMcpKeys()


    async function onRevokeClick() {
        await removeConnector(mcpKeyId, connector.mcpName ?? "")
    }

    return (
        <div
            className={cls.ConnectorRowWrapper}
            onClick={onEdit}
        >
            <ConnectorChip
                connector={connector}
                connections={connections}
            />
            <div
                className={cls.ConnectorRowActions}
            >
                <Button
                    variant="ghost"
                    className={cls.BtnSmall}
                    onClick={onRevokeClick}>
                    Remove
                </Button>
            </div>
        </div>
    )
}
