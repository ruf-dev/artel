import {useState, useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/components/ManageKeyDialog/ManageKeyDialog.module.css"

import {McpKeyInfo, McpConnectorInfo, MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useVaults} from "@/app/hooks/Vaults.ts"

import ModalClose from "@/components/ModalClose/ModalClose.tsx"
import ModalActions from "@/components/ModalActions/ModalActions.tsx"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import ConnectorChip, {connectionLabel} from "@/components/ConnectorChip/ConnectorChip.tsx"
import Button from "@/components/shared/Button/Button.tsx"

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
        momCandidates,
        fetchMomCandidates,
    } = useMcpKeys()
    const {CloseDialog} = useDialog()
    const {vaults} = useVaults()
    const navigate = useNavigate()

    const connectors = mcpKey.id ? connectorsByKey[mcpKey.id] ?? [] : []
    const vault = vaults.find(v => v.id === selectedVaultId)

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

    async function handleRemoveConnector(mcpName: string) {
        if (!mcpKey.id) return
        await removeConnector(mcpKey.id, mcpName)
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

    if (step === "vault") {
        return (
            <div className={cls.Overlay}>
                <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                     aria-labelledby="manageVaultTitle">
                    <div className={cls.ModalHead}>
                        <h2 className={cls.ModalTitle} id="manageVaultTitle">Select vault</h2>
                        <ModalClose onClick={CloseDialog} disabled={saving} className={cls.ModalClose}/>
                    </div>
                    <p className={cls.ModalSub}>Choose which vault this key connects to.</p>
                    <div className={cls.OptionList}>
                        {vaults?.map(v => (
                            <SelectOption
                                key={v.id}
                                label={v.name ?? ""}
                                selected={selectedVaultId === v.id}
                                onSelect={() => setSelectedVaultId(v.id ?? "")}
                            />
                        ))}
                        {vaults.length === 0 && (
                            <p className={cls.Empty}>No vaults found.</p>
                        )}
                    </div>
                    <ModalActions
                        containerClassName={cls.ModalActions}
                        buttons={[
                            {
                                label: "Back",
                                onClick: () => setStep("main"),
                                className: cls.BtnGhost,
                                disabled: saving
                            },
                            {
                                label: saving ? "Saving…" : "Save",
                                onClick: handleSaveVault,
                                className: cls.BtnPrimary,
                                disabled: saving || !selectedVaultId
                            },
                        ]}
                    />
                </div>
            </div>
        )
    }

    if (step === "addConnection") {
        return (
            <div className={cls.Overlay}>
                <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                     aria-labelledby="addConnectionTitle">
                    <div className={cls.ModalHead}>
                        <h2 className={cls.ModalTitle} id="addConnectionTitle">Add connection</h2>
                        <ModalClose onClick={CloseDialog} className={cls.ModalClose}/>
                    </div>
                    <p className={cls.ModalSub}>Pick a service to connect to this key.</p>
                    <div className={cls.OptionList}>
                        {momCandidates.map(c => (
                            <MomCandidateCard
                                key={c.name}
                                candidate={c}
                                connected={connectors.some(con => con.mcpName === c.name)}
                                onSelect={() => handleSelectCandidate(c)}
                            />
                        ))}
                        {momCandidates.length === 0 && (
                            <p className={cls.Empty}>No services available yet.</p>
                        )}
                    </div>
                    <ModalActions
                        containerClassName={cls.ModalActions}
                        buttons={[
                            {label: "Back", onClick: () => setStep("main"), className: cls.BtnGhost, disabled: false},
                        ]}
                    />
                </div>
            </div>
        )
    }

    if (step === "selectConnection" && selectedCandidate) {
        const available = selectedCandidate.connections ?? []

        return (
            <div className={cls.Overlay}>
                <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                     aria-labelledby="selectConnectionTitle">
                    <div className={cls.ModalHead}>
                        <h2 className={cls.ModalTitle} id="selectConnectionTitle">{selectedCandidate.name}</h2>
                        <ModalClose onClick={CloseDialog} disabled={saving} className={cls.ModalClose}/>
                    </div>
                    <p className={cls.ModalSub}>Pick the connection this key should use for it.</p>
                    <div className={cls.OptionList}>
                        {available.map(c => (
                            <SelectOption
                                key={c.id}
                                label={connectionLabel(c)}
                                selected={selectedExternalConnectionId === c.id}
                                onSelect={() => setSelectedExternalConnectionId(c.id ?? "")}
                            />
                        ))}
                        {available.length === 0 && (
                            <p className={cls.Empty}>
                                No connections yet.{" "}
                                <button
                                    className={cls.LinkBtn}
                                    type="button"
                                    onClick={() => {
                                        CloseDialog()
                                        navigate(Path.ConnectionsPage)
                                    }}
                                >
                                    Set one up
                                </button>
                            </p>
                        )}
                    </div>
                    <ModalActions
                        containerClassName={cls.ModalActions}
                        buttons={[
                            {
                                label: "Back",
                                onClick: () => setStep(editingConnectorName ? "main" : "addConnection"),
                                className: cls.BtnGhost,
                                disabled: saving,
                            },
                            {
                                label: saving ? (editingConnectorName ? "Saving…" : "Adding…") : (editingConnectorName ? "Save" : "Add"),
                                onClick: handleAddConnector,
                                className: cls.BtnPrimary,
                                disabled: saving || !selectedExternalConnectionId,
                            },
                        ]}
                    />
                </div>
            </div>
        )
    }

    return (
        <div className={cls.Overlay}>
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                 aria-labelledby="manageKeyTitle">
                <div className={cls.ModalHead}>
                    <h2 className={cls.ModalTitle} id="manageKeyTitle">Manage <b>"{mcpKey.name}"</b></h2>
                    <ModalClose onClick={CloseDialog} className={cls.ModalClose}/>
                </div>

                <div className={cls.Field}>
                    <span className={cls.FieldLabel}>Vault</span>
                    <div className={cls.ConnectorRowWrapper}>
                        <span
                            className={vault ? cls.VaultChip : `${cls.Chip} ${cls.ChipMuted}`}
                            data-tooltip-id="root-tooltip"
                            data-tooltip-content={vault ? `Vault: ${vault.name}` : "No vault assigned"}
                        >
                            {vault && <span className={cls.VaultBadge}>A</span>}
                            {vault ? vault.name : "No vault"}
                        </span>
                        <button className={cls.BtnGhost} onClick={() => setStep("vault")} type="button">Change</button>
                    </div>
                </div>

                <div className={cls.Field}>
                    <span className={cls.FieldLabel}>Connections</span>
                    <ConnectorList connectors={connectors} onRemove={handleRemoveConnector} onEdit={handleEditConnector}/>
                </div>

                <ModalActions
                    containerClassName={cls.ModalActions}
                    buttons={[
                        {label: "Close", onClick: CloseDialog, className: cls.BtnGhost, disabled: false},
                        {
                            label: "Add connection",
                            onClick: () => setStep("addConnection"),
                            className: cls.BtnPrimary,
                            disabled: false,
                        },
                    ]}
                />
            </div>
        </div>
    )
}

function MomCandidateCard({candidate, connected, onSelect}: {
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

function ConnectorList({connectors, onRemove, onEdit}: {
    connectors: McpConnectorInfo[]
    onRemove: (mcpName: string) => void
    onEdit: (connector: McpConnectorInfo) => void
}) {
    if (connectors.length === 0) {
        return <p className={cls.Empty}>No connections linked to this key yet.</p>
    }
    return (
        <div className={cls.OptionList}>
            {connectors.map(c => (
                <ConnectorRow key={c.mcpName} connector={c} onRemove={() => onRemove(c.mcpName ?? "")} onEdit={() => onEdit(c)}/>
            ))}
        </div>
    )
}

function ConnectorRow({connector, onRemove, onEdit}: { connector: McpConnectorInfo; onRemove: () => void; onEdit: () => void }) {
    const {connections} = useExternalConnections()
    return (
        <div className={cls.ConnectorRowWrapper}>
            <ConnectorChip connector={connector} connections={connections}/>
            <div className={cls.ConnectorRowActions}>
                <button className={cls.IconBtn} onClick={onEdit} type="button" title="Edit connector" aria-label="Edit connector">
                    <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M12 20h9"/><path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4z"/>
                    </svg>
                </button>
                <Button variant="ghost" className={cls.BtnSmall} onClick={onRemove}>Remove</Button>
            </div>
        </div>
    )
}
