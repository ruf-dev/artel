import {useState, useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/mcp-keys/McpKeysPage.module.css"

import {McpKeyInfo, McpConnectorInfo, MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {ExternalConnectionInfo, ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useVaults} from "@/app/hooks/Vaults.ts"
import useUser from "@/hooks/user/User.ts"

import ModalClose from "@/components/ModalClose/ModalClose.tsx"
import ModalActions from "@/components/ModalActions/ModalActions.tsx"
import FormField from "@/components/FormField/FormField.tsx"

export default function McpKeysPage() {
    const navigate = useNavigate()
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const {fetch: fetchKeys} = useMcpKeys()
    const {fetch: fetchExternalConnections} = useExternalConnections()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetchKeys()
            void fetchExternalConnections()
        }
    }, [auth, fetchKeys, fetchExternalConnections])

    return (
        <div className={cls.Root}>
            <HeroSegment onCreateClick={() => OpenDialog(<CreateKeyDialog/>)}/>
            <ContentSegment/>
        </div>
    )
}

function HeroSegment({onCreateClick}: { onCreateClick: () => void }) {
    const {keys, loading} = useMcpKeys()

    return (
        <div className={cls.Hero}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>MCP</div>
                <h1 className={cls.HeroTitle}>API Keys</h1>
                <p className={cls.HeroSub}>
                    <b>{loading ? "…" : `${keys.length} ${keys.length === 1 ? "key" : "keys"}`}</b>
                    {" · "}<span>bridge your MCP agents to Artel</span>
                </p>
            </div>
            <button className={cls.AddBtn} onClick={onCreateClick}>
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2"
                     strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                New key
            </button>
        </div>
    )
}

function ContentSegment() {
    const {keys, loading} = useMcpKeys()
    const {OpenDialog} = useDialog()

    if (loading) {
        return (
            <div className={cls.Content}>
                <p className={cls.Empty}>Loading…</p>
            </div>
        )
    }

    return (
        <div className={cls.Content}>
            <div className={cls.List}>
                {keys.map(key => (
                    <McpKeyCard
                        key={key.id}
                        mcpKey={key}
                        onRevoke={() => OpenDialog(<RevokeKeyDialog mcpKey={key}/>)}
                        onManage={() => OpenDialog(<ManageKeyDialog mcpKey={key}/>)}
                    />
                ))}
                {keys.length === 0 && (
                    <p className={cls.Empty}>No API keys yet. Create one to get started.</p>
                )}
            </div>
        </div>
    )
}

function McpKeyCard({mcpKey, onRevoke, onManage}: {
    mcpKey: McpKeyInfo
    onRevoke: () => void
    onManage: () => void
}) {
    const {vaults} = useVaults()
    const {connectorsByKey, fetchConnectors} = useMcpKeys()
    const {connections} = useExternalConnections()

    const vault = vaults.find(v => v.id === mcpKey.vaultId)
    const connectors = mcpKey.id ? connectorsByKey[mcpKey.id] ?? [] : []

    useEffect(() => {
        if (mcpKey.id) {
            void fetchConnectors(mcpKey.id)
        }
    }, [mcpKey.id, fetchConnectors])

    function formatDate(iso: string | undefined): string {
        if (!iso) return "Never"
        const d = new Date(iso)
        if (isNaN(d.getTime())) return "Never"
        return d.toLocaleDateString(undefined, {year: "numeric", month: "short", day: "numeric"})
    }

    return (
        <div className={cls.Card}>
            <div className={cls.CardMain}>
                <div className={cls.CardHeader}>
                    <span className={cls.CardName}>{mcpKey.name}</span>
                    <span className={cls.CardPreview}>{mcpKey.keyPreview}…</span>
                </div>
                <div className={cls.CardChips}>
                    {vault ? (
                        <span className={cls.VaultChip} title={`Vault: ${vault.name}`}>
                            <span className={cls.VaultBadge}>A</span>
                            {vault.name}
                        </span>
                    ) : (
                        <span className={`${cls.Chip} ${cls.ChipMuted}`} title="No vault assigned">No vault</span>
                    )}
                    {connectors.map(c => (
                        <ConnectorChip key={c.mcpName} connector={c} connections={connections}/>
                    ))}
                </div>
                <div className={cls.CardMeta}>
                    Last accessed: {formatDate(mcpKey.lastAccessedAt)}
                </div>
            </div>
            <div className={cls.CardActions}>
                <button className={cls.BtnGhost} onClick={onManage} type="button">Manage</button>
                <button className={cls.BtnDanger} onClick={onRevoke} type="button">Revoke</button>
            </div>
        </div>
    )
}

function RevokeKeyDialog({mcpKey}: { mcpKey: McpKeyInfo }) {
    const [revoking, setRevoking] = useState(false)
    const {revoke} = useMcpKeys()
    const {CloseDialog} = useDialog()

    async function handleRevoke() {
        if (!mcpKey.id || !mcpKey.vaultId) return
        setRevoking(true)
        try {
            await revoke(mcpKey.id, mcpKey.vaultId)
            CloseDialog()
        } finally {
            setRevoking(false)
        }
    }

    return (
        <div className={cls.Overlay}>
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                 aria-labelledby="revokeKeyTitle">
                <div className={cls.ModalHead}>
                    <h2 className={cls.ModalTitle} id="revokeKeyTitle">Revoke key</h2>
                    <ModalClose onClick={CloseDialog} disabled={revoking} className={cls.ModalClose}/>
                </div>
                <p className={cls.ModalSub}>
                    Revoke <b>"{mcpKey.name}"</b>? Any agents using this key will immediately lose access. This cannot
                    be undone.
                </p>
                <ModalActions
                    containerClassName={cls.ModalActions}
                    buttons={[
                        {
                            label: "Cancel",
                            onClick: CloseDialog,
                            className: cls.BtnGhost,
                            disabled: revoking,
                        },
                        {
                            label: revoking ? "Revoking…" : "Revoke",
                            onClick: handleRevoke,
                            className: cls.BtnDanger,
                            disabled: revoking,
                        },
                    ]}
                />
            </div>
        </div>
    )
}

type ManageStep = "main" | "vault" | "addConnection" | "selectConnection"

function ManageKeyDialog({mcpKey}: { mcpKey: McpKeyInfo }) {
    const [step, setStep] = useState<ManageStep>("main")
    const [saving, setSaving] = useState(false)
    const [selectedVaultId, setSelectedVaultId] = useState(mcpKey.vaultId ?? "")
    const [selectedCandidate, setSelectedCandidate] = useState<MomCandidate | null>(null)
    const [selectedExternalConnectionId, setSelectedExternalConnectionId] = useState("")

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

    async function handleAddConnector() {
        if (!mcpKey.id || !selectedCandidate?.name || !selectedExternalConnectionId) return
        setSaving(true)
        try {
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
                                onClick: () => setStep("addConnection"),
                                className: cls.BtnGhost,
                                disabled: saving,
                            },
                            {
                                label: saving ? "Adding…" : "Add",
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
                            title={vault ? `Vault: ${vault.name}` : "No vault assigned"}
                        >
                            {vault && <span className={cls.VaultBadge}>A</span>}
                            {vault ? vault.name : "No vault"}
                        </span>
                        <button className={cls.BtnGhost} onClick={() => setStep("vault")} type="button">Change</button>
                    </div>
                </div>

                <div className={cls.Field}>
                    <span className={cls.FieldLabel}>Connections</span>
                    <ConnectorList connectors={connectors} onRemove={handleRemoveConnector}/>
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

function ConnectorList({connectors, onRemove}: {
    connectors: McpConnectorInfo[]
    onRemove: (mcpName: string) => void
}) {
    if (connectors.length === 0) {
        return <p className={cls.Empty}>No connections linked to this key yet.</p>
    }
    return (
        <div className={cls.OptionList}>
            {connectors.map(c => (
                <ConnectorRow key={c.mcpName} connector={c} onRemove={() => onRemove(c.mcpName ?? "")}/>
            ))}
        </div>
    )
}

function ConnectorRow({connector, onRemove}: { connector: McpConnectorInfo; onRemove: () => void }) {
    const {connections} = useExternalConnections()
    return (
        <div className={cls.ConnectorRowWrapper}>
            <ConnectorChip connector={connector} connections={connections}/>
            <button className={cls.BtnGhostSmall} onClick={onRemove} type="button">Remove</button>
        </div>
    )
}

function connectionLabel(c: ExternalConnectionInfo): string {
    if (c.google) return c.google.email ?? "Google account"
    if (c.provider === ExternalProvider.EXTERNAL_PROVIDER_EMAIL) {
        return c.generic?.fields?.username ?? "Email account"
    }
    return c.provider?.replace("EXTERNAL_PROVIDER_", "").toLowerCase() ?? "Connection"
}

function ConnectorChip({connector, connections}: {
    connector: McpConnectorInfo
    connections: ExternalConnectionInfo[]
}) {
    const match = connections.find(c => c.id === connector.externalConnectionId)
    const label = match ? connectionLabel(match) : undefined
    const atIndex = label?.indexOf("@") ?? -1

    if (!label || atIndex === -1) {
        return (
            <span className={cls.Chip} title={`${connector.mcpName} connection`}>
                {connector.mcpName}
            </span>
        )
    }

    const color = mailDomainColor(label.slice(atIndex + 1))

    return (
        <span className={cls.MailChip} style={{borderColor: color}} title={`Email connection: ${label}`}>
            <span className={cls.MailAtBadge} style={{background: color}}>@</span>
            {label}
        </span>
    )
}

const KNOWN_MAIL_DOMAIN_COLORS: Record<string, string> = {
    "yandex.ru": "#FFCC00",
    "yandex.com": "#FFCC00",
    "ya.ru": "#FFCC00",
    "gmail.com": "#4285F4",
    "google.com": "#4285F4",
    "googlemail.com": "#4285F4",
    "outlook.com": "#0078D4",
    "hotmail.com": "#0078D4",
    "live.com": "#0078D4",
    "msn.com": "#0078D4",
    "mail.ru": "#005FF9",
    "inbox.ru": "#005FF9",
    "list.ru": "#005FF9",
    "bk.ru": "#005FF9",
    "icloud.com": "#A2AAAD",
    "me.com": "#A2AAAD",
    "mac.com": "#A2AAAD",
    "protonmail.com": "#6D4AFF",
    "proton.me": "#6D4AFF",
    "yahoo.com": "#6001D2",
}

function mailDomainColor(domain: string): string {
    const known = KNOWN_MAIL_DOMAIN_COLORS[domain.toLowerCase()]
    if (known) return known

    let hash = 0
    for (let i = 0; i < domain.length; i++) {
        hash = (hash * 31 + domain.charCodeAt(i)) >>> 0
    }
    return `hsl(${hash % 360}, 65%, 45%)`
}

function SelectOption({label, selected, onSelect}: { label: string; selected: boolean; onSelect: () => void }) {
    return (
        <button
            className={selected ? `${cls.OptionRow} ${cls.OptionRowSelected}` : cls.OptionRow}
            onClick={onSelect}
            type="button"
        >
            <span className={cls.OptionRadio}>{selected ? "●" : "○"}</span>
            <span>{label}</span>
        </button>
    )
}

function CreateKeyDialog() {
    const [creating, setCreating] = useState(false)
    const [name, setName] = useState("")
    const [selectedVaultId, setSelectedVaultId] = useState("")
    const [rawToken, setRawToken] = useState("")
    const [copied, setCopied] = useState(false)

    const {create} = useMcpKeys()
    const {CloseDialog} = useDialog()
    const {vaults} = useVaults()

    async function handleCreate() {
        if (!name || !selectedVaultId) return
        setCreating(true)
        try {
            const resp = await create(name, selectedVaultId)
            setRawToken(resp.rawToken ?? "")
        } finally {
            setCreating(false)
        }
    }

    async function handleCopy() {
        await navigator.clipboard.writeText(rawToken)
        setCopied(true)
        setTimeout(() => setCopied(false), 2000)
    }

    if (rawToken) {
        return (
            <div className={cls.Overlay}>
                <div
                    className={cls.Modal}
                    onClick={e =>
                        e.stopPropagation()}
                    role="dialog"
                    aria-modal="true"
                    aria-labelledby="createdKeyTitle">
                    <div className={cls.ModalHead}>
                        <h2 className={cls.ModalTitle} id="createdKeyTitle">Key created</h2>
                    </div>
                    <p className={cls.ModalSub}>
                        Copy this key now — it will not be shown again.
                    </p>
                    <div className={cls.TokenBox}>
                        <span className={cls.TokenText}>{rawToken}</span>
                        <button className={cls.BtnCopy} onClick={handleCopy} type="button">
                            {copied ? "Copied!" : "Copy"}
                        </button>
                    </div>
                    <ModalActions
                        containerClassName={cls.ModalActions}
                        buttons={[
                            {
                                label: "Done",
                                onClick: CloseDialog,
                                className: cls.BtnPrimary,
                                disabled: false,
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
                 aria-labelledby="createKeyTitle">
                <div className={cls.ModalHead}>
                    <h2 className={cls.ModalTitle} id="createKeyTitle">New API key</h2>
                    <ModalClose onClick={CloseDialog} disabled={creating} className={cls.ModalClose}/>
                </div>
                <p className={cls.ModalSub}>Name your key and pick the vault it will access.</p>

                <FormField
                    label="Key name"
                    placeholder="e.g. My Assistant"
                    onChange={setName}
                    disabled={creating}
                    fieldClassName={cls.Field}
                    labelClassName={cls.FieldLabel}
                    inputClassName={cls.Input}
                />

                <div className={cls.Field}>
                    <span className={cls.FieldLabel}>Vault</span>
                    <div className={cls.OptionList}>
                        {vaults.map((v: VaultItem) => (
                            <SelectOption
                                key={v.id}
                                label={v.name ?? ""}
                                selected={selectedVaultId === v.id}
                                onSelect={() => setSelectedVaultId(v.id ?? "")}
                            />
                        ))}
                        {vaults?.length === 0 && (
                            <p className={cls.Empty}>No vaults available. Create one first.</p>
                        )}
                    </div>
                </div>

                <ModalActions
                    containerClassName={cls.ModalActions}
                    buttons={[
                        {
                            label: creating ? "Creating…" : "Create key",
                            onClick: handleCreate,
                            className: cls.BtnPrimary,
                            disabled: creating || !name || !selectedVaultId,
                        }
                    ]}
                />
            </div>
        </div>
    )
}

