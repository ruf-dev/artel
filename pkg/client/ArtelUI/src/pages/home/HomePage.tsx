import {useState, useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/home/HomePage.module.css"

import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useVaults} from "@/app/hooks/Vaults.ts"
import useUser from "@/hooks/user/User.ts"

import VaultCard from "@/widgets/VaultCard/VaultCard.tsx"
import Topbar from "@/segments/Topbar/Topbar.tsx"
import ModalClose from "@/components/ModalClose/ModalClose.tsx"
import FormField from "@/components/FormField/FormField.tsx"
import ModalActions from "@/components/ModalActions/ModalActions.tsx"

export default function HomePage() {
    const navigate = useNavigate()

    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const {vaults, fetch: fetchVaults, forbidden} = useVaults()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetchVaults()
        }
    }, [auth, fetchVaults])

    useEffect(() => {
        if (forbidden) navigate(Path.ClosedAlpha)
    }, [forbidden, navigate])


    function openEditDialog(vaultId: string) {
        const vault = vaults.find(v => v.id === vaultId)
        if (!vault) return
        const {CloseDialog} = useDialog.getState()
        OpenDialog(
            <EditVaultDialog
                vault={vault}
                onClose={CloseDialog}
                onDeleted={CloseDialog}
            />
        )
    }

    return (
        <div className={cls.Root}>
            <Topbar/>
            <HeroSegment onCreateClick={() => OpenDialog(<CreateVaultDialog/>)}/>
            <ContentSegment onEditClick={openEditDialog}/>
        </div>
    )
}

function HeroSegment({onCreateClick}: { onCreateClick: () => void }) {
    const {vaults, loading} = useVaults()

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

function ContentSegment({onEditClick}: { onEditClick: (id: string) => void }) {
    const {vaults, loading} = useVaults()

    let loadingState = null

    if (loading) {
        loadingState = (
            <p className={cls.Empty}>Loading…</p>
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


function CreateVaultDialog() {
    const [isCreating, setIsCreating] = useState(false)
    const [vaultName, setVaultName] = useState("")

    const {create} = useVaults()
    const {CloseDialog} = useDialog()


    async function handleCreate() {
        setIsCreating(true)
        const name = vaultName.trim()
        create(name)
            .then(CloseDialog)
            .catch(err => {
                console.error(err)
            })
            .finally(() => {
                setIsCreating(false)
            })
    }

    return (
        <div className={cls.Overlay}>
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                 aria-labelledby="createModalTitle">
                <div className={cls.ModalHead}>
                    <h2 className={cls.ModalTitle} id="createModalTitle">New vault</h2>
                    <ModalClose onClick={CloseDialog} disabled={isCreating} className={cls.ModalClose}/>
                </div>
                <p className={cls.ModalSub}>Give it a name — that's all we need to spin one up.</p>

                <FormField
                    label="Vault name"
                    placeholder="e.g. Marketplace inventory"
                    onChange={setVaultName}
                    disabled={isCreating}
                    maxLength={48}
                    fieldClassName={cls.Field}
                    labelClassName={cls.FieldLabel}
                    inputClassName={cls.Input}
                />

                <ModalActions
                    containerClassName={cls.ModalActions}
                    buttons={[
                        {
                            label: isCreating ? "Creating…" : "Create vault",
                            onClick: handleCreate,
                            className: cls.BtnPrimary,
                            disabled: isCreating
                        }
                    ]}
                />
            </div>
        </div>
    )
}

function EditVaultDialog({vault, onClose, onDeleted}: {
    vault: VaultItem
    onClose: () => void
    onDeleted: () => void
}) {
    const {remove} = useVaults()
    const [deleting, setDeleting] = useState(false)

    async function handleDelete() {
        if (!vault.id) return
        setDeleting(true)
        try {
            await remove(vault.id)
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
                <ModalClose onClick={onClose} className={cls.ModalClose}/>
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

            <ModalActions
                containerClassName={cls.ModalActions}
                buttons={[
                    {
                        label: "Close",
                        onClick: onClose,
                        className: cls.BtnGhost
                    }
                ]}
            />
        </div>
    )
}
