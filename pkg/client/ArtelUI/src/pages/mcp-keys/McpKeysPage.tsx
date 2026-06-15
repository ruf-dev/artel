import {useState, useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/mcp-keys/McpKeysPage.module.css"

import {McpKeyInfo} from "@/app/api/artel/mcp_keys.pb.ts"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
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

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetchKeys()
        }
    }, [auth, fetchKeys])

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

    const vault = vaults.find(v => v.id === mcpKey.vaultId)

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
                        <span className={cls.Chip}>{vault.name}</span>
                    ) : (
                        <span className={`${cls.Chip} ${cls.ChipMuted}`}>No vault</span>
                    )}
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

type ManageStep = "choose" | "vault"

function ManageKeyDialog({mcpKey}: { mcpKey: McpKeyInfo }) {
    const [step, setStep] = useState<ManageStep>("choose")
    const [saving, setSaving] = useState(false)
    const [selectedVaultId, setSelectedVaultId] = useState(mcpKey.vaultId ?? "")

    const {setAccess} = useMcpKeys()
    const {CloseDialog} = useDialog()
    const {vaults} = useVaults()

    async function handleSave() {
        if (!mcpKey.id) return
        setSaving(true)
        try {
            await setAccess(mcpKey.id, selectedVaultId)
            CloseDialog()
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
                                onClick: () => setStep("choose"),
                                className: cls.BtnGhost,
                                disabled: saving
                            },
                            {
                                label: saving ? "Saving…" : "Save",
                                onClick: handleSave,
                                className: cls.BtnPrimary,
                                disabled: saving || !selectedVaultId
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
                    <h2 className={cls.ModalTitle} id="manageKeyTitle">Manage access</h2>
                    <ModalClose onClick={CloseDialog} className={cls.ModalClose}/>
                </div>
                <p className={cls.ModalSub}>Choose what you want to configure for <b>"{mcpKey.name}"</b>.</p>
                <div className={cls.EntityCards}>
                    <EntityCard
                        title="Vault"
                        description="Change which vault this key connects to"
                        onClick={() => setStep("vault")}
                    />
                </div>
                <ModalActions
                    containerClassName={cls.ModalActions}
                    buttons={[
                        {label: "Cancel", onClick: CloseDialog, className: cls.BtnGhost, disabled: false},
                    ]}
                />
            </div>
        </div>
    )
}

function EntityCard({title, description, onClick}: { title: string; description: string; onClick: () => void }) {
    return (
        <button className={cls.EntityCard} onClick={onClick} type="button">
            <span className={cls.EntityCardTitle}>{title}</span>
            <span className={cls.EntityCardDesc}>{description}</span>
        </button>
    )
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

