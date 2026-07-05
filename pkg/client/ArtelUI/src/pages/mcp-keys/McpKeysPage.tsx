import {useState, useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/mcp-keys/McpKeysPage.module.css"

import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useVaults} from "@/app/hooks/Vaults.ts"
import useUser from "@/hooks/user/User.ts"

import {Button, ModalClose} from "@vervstack/chures"
import FormField from "@/components/FormField/FormField.tsx"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import ManageKeyDialog from "@/dialogs/ManageKeyDialog/ManageKeyDialog.tsx"
import McpKeyCard from "@/widgets/McpKeyCard/McpKeyCard.tsx"

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
            <Button variant="primary" onClick={onCreateClick}>
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2"
                     strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                New key
            </Button>
        </div>
    )
}

function ContentSegment() {
    const {keys, loading} = useMcpKeys()
    const {OpenDialog} = useDialog()

    // TODO add loader from chures
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
                        <Button variant="secondary" onClick={handleCopy}>
                            {copied ? "Copied!" : "Copy"}
                        </Button>
                    </div>
                    <div className={cls.ModalActions}>
                        <Button variant="primary" onClick={CloseDialog}>Done</Button>
                    </div>
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

                <div className={cls.ModalActions}>
                    <Button variant="primary" onClick={handleCreate} disabled={creating || !name || !selectedVaultId}>
                        {creating ? "Creating…" : "Create key"}
                    </Button>
                </div>
            </div>
        </div>
    )
}
