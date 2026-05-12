import {useCallback, useEffect, useRef, useState} from "react"
import type {KeyboardEvent} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/home/HomePage.module.css"

import {AuthMiddleware} from "@/processes/AuthMiddleware.ts"
import {VaultService} from "@/processes/Vaults.ts"
import {Path} from "@/app/routing/Router.tsx"
import {VaultsAPI} from "@/app/api/artel/vaults.pb.ts"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import VaultCard from "@/widgets/VaultCard/VaultCard.tsx"
import Topbar from "@/components/Topbar/Topbar.tsx"
import EmptyState from "@/components/EmptyState/EmptyState.tsx"

interface Props {
    auth: AuthMiddleware
    onLogout: () => void
}

const vaultService = new VaultService()

export default function HomePage({auth, onLogout}: Props) {
    const navigate = useNavigate()
    const {OpenDialog} = useDialog()
    const [vaultName, setVaultName] = useState("")
    const [creating, setCreating] = useState(false)
    const [vaults, setVaults] = useState<VaultItem[]>([])
    const [loading, setLoading] = useState(true)
    const inputRef = useRef<HTMLInputElement>(null)

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])


    const fetchVaults = useCallback(async () => {
        setLoading(true)
        try {
            const list = await vaultService.ListVaults()
            setVaults(list)
        } finally {
            setLoading(false)
        }
    }, [])

    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetchVaults()
        }
    }, [auth, fetchVaults])


    function handleLogout() {
        onLogout()
        navigate(Path.InitPage)
    }

    function openDialog() {
        setVaultName("")
        const {CloseDialog} = useDialog.getState()
        OpenDialog(
            <CreateVaultDialog
                vaultName={vaultName}
                creating={creating}
                onClose={CloseDialog}
                onNameChange={e => setVaultName(e.target.value)}
                onKeyDown={handleKeyDown}
                onCreate={handleCreate}
                inputRef={inputRef}
            />
        )
    }

    async function handleCreate() {
        const name = vaultName.trim()
        if (!name) return
        setCreating(true)
        try {
            await VaultsAPI.CreateVault({name}, auth.getInitReq())
            const {CloseDialog} = useDialog.getState()
            CloseDialog()
            void fetchVaults()
        } finally {
            setCreating(false)
        }
    }

    function openEditDialog(vaultId: string) {
        const vault = vaults.find(v => v.id === vaultId)
        if (!vault) return
        const {CloseDialog} = useDialog.getState()
        OpenDialog(
            <EditVaultDialog
                vault={vault}
                auth={auth}
                onClose={CloseDialog}
                onDeleted={() => {
                    CloseDialog()
                    void fetchVaults()
                }}
            />
        )
    }

    function handleKeyDown(e: KeyboardEvent<HTMLInputElement>) {
        if (e.key === "Enter") void handleCreate()
        if (e.key === "Escape") {
            const {CloseDialog} = useDialog.getState()
            CloseDialog()
        }
    }

    return (
        <div className={cls.Root}>
            <Topbar onLogout={handleLogout}/>
            <HeroSegment loading={loading} vaults={vaults} onCreateClick={openDialog}/>
            <ContentSegment loading={loading} vaults={vaults} onEditClick={openEditDialog}/>
        </div>
    )
}

function HeroSegment({loading, vaults, onCreateClick}: {
    loading: boolean;
    vaults: VaultItem[];
    onCreateClick: () => void
}) {
    return (
        <div className={cls.Hero}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>Workspace</div>
                <h1 className={cls.HeroTitle}>Your vaults</h1>
                <p className={cls.HeroSub}>
                    <b>{loading ? "…" : `${vaults.length} ${vaults.length === 1 ? "vault" : "vaults"}`}</b>
                    {" · "}<span>all systems operational</span>
                </p>
            </div>
            <button className={cls.NewVaultBtn} onClick={onCreateClick}>
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2"
                     strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                New vault
            </button>
        </div>
    )
}

function ContentSegment({loading, vaults, onEditClick}: {
    loading: boolean;
    vaults: VaultItem[];
    onEditClick: (id: string) => void
}) {

    let loadingState = null

    if (loading) {
        loadingState = (
            <p className={cls.Empty}>Loading…</p>
        )
    } else if (vaults.length === 0) {
        loadingState = (
            <EmptyState
                onCreateClick={() => {/* handled by HeroSegment */
                }}/>
        )
    }


    return (
        <div className={cls.Content}>
            {loadingState}
            {
                !loading && <div className={cls.Grid}>
                    {vaults.map(v => (
                        <VaultCard key={v.id} vault={v} onEdit={onEditClick}/>
                    ))}
				</div>
            }
        </div>
    )
}

function CreateVaultDialog({vaultName, creating, onClose, onNameChange, onKeyDown, onCreate, inputRef}: {
    vaultName: string
    creating: boolean
    onClose: () => void
    onNameChange: (e: React.ChangeEvent<HTMLInputElement>) => void
    onKeyDown: (e: KeyboardEvent<HTMLInputElement>) => void
    onCreate: () => void
    inputRef: React.RefObject<HTMLInputElement | null>
}) {
    return (
        <div className={cls.Overlay} onClick={onClose}>
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                 aria-labelledby="createModalTitle">
                <div className={cls.ModalHead}>
                    <h2 className={cls.ModalTitle} id="createModalTitle">New vault</h2>
                    <button className={cls.ModalClose} type="button" onClick={onClose} aria-label="Close">
                        <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor"
                             strokeWidth="1.8" strokeLinecap="round">
                            <line x1="18" y1="6" x2="6" y2="18"/>
                            <line x1="6" y1="6" x2="18" y2="18"/>
                        </svg>
                    </button>
                </div>
                <p className={cls.ModalSub}>Give it a name — that's all we need to spin one up.</p>

                <label className={cls.Field}>
                    <span className={cls.FieldLabel}>Vault name</span>
                    <input
                        ref={inputRef}
                        className={cls.Input}
                        placeholder="e.g. Marketplace inventory"
                        value={vaultName}
                        onChange={onNameChange}
                        onKeyDown={onKeyDown}
                        disabled={creating}
                        maxLength={48}
                        autoComplete="off"
                    />
                </label>

                <div className={cls.ModalActions}>
                    <button className={cls.BtnGhost} type="button" onClick={onClose} disabled={creating}>
                        Cancel
                    </button>
                    <button
                        className={cls.BtnPrimary}
                        type="button"
                        onClick={onCreate}
                        disabled={creating || !vaultName.trim()}
                    >
                        {creating ? "Creating…" : "Create vault"}
                    </button>
                </div>
            </div>
        </div>
    )
}

function EditVaultDialog({vault, auth, onClose, onDeleted}: {
    vault: VaultItem
    auth: AuthMiddleware
    onClose: () => void
    onDeleted: () => void
}) {
    const [deleting, setDeleting] = useState(false)

    async function handleDelete() {
        setDeleting(true)
        try {
            await VaultsAPI.DeleteVault({id: vault.id}, auth.getInitReq())
            onDeleted()
        } finally {
            setDeleting(false)
        }
    }

    return (
        <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
             aria-labelledby="editModalTitle">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle} id="editModalTitle">Edit vault</h2>
                <button className={cls.ModalClose} type="button" onClick={onClose} aria-label="Close">
                    <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" strokeWidth="1.8"
                         strokeLinecap="round">
                        <line x1="18" y1="6" x2="6" y2="18"/>
                        <line x1="6" y1="6" x2="18" y2="18"/>
                    </svg>
                </button>
            </div>
            <p className={cls.ModalSub}>Rename or delete this vault.</p>

            <div className={cls.FieldLabel} style={{marginBottom: 4}}>Vault name</div>
            <div className={cls.VaultNameDisplay}>{vault.name}</div>

            <div className={cls.DangerZone}>
                <div className={cls.DangerZoneRow}>
                    <div>
                        <div className={cls.DangerZoneTitle}>Delete this vault</div>
                        <div className={cls.DangerZoneSub}>Permanent. Connection string stops working immediately.</div>
                    </div>
                    <button
                        className={cls.BtnDanger}
                        type="button"
                        onClick={handleDelete}
                        disabled={deleting}
                    >
                        {deleting ? "Deleting…" : "Delete"}
                    </button>
                </div>
            </div>

            <div className={cls.ModalActions}>
                <button className={cls.BtnGhost} type="button" onClick={onClose}>
                    Close
                </button>
            </div>
        </div>
    )
}
