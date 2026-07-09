import {useState, useEffect} from "react"

import cls from "@/pages/home/HomePage.module.css"

import {useDialog} from "@/app/hooks/Dialog"
import {useVaults, useVaultMutations} from "@/app/hooks/Vaults.ts"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"

import {Button, ModalClose, InfoDialog} from "@vervstack/chures"
import VaultCard from "@/widgets/VaultCard/VaultCard.tsx"
import FormField from "@/components/FormField/FormField.tsx"
import ManageVaultDialog from "@/dialogs/ManageVaultDialog/ManageVaultDialog.tsx"
import type {GrpcStatusError} from "@/processes/grpcErrors.ts";
import {isMissingSubscription} from "@/processes/UserErrors.ts";

export default function HomePage() {
    const {OpenDialog, CloseDialog} = useDialog()
    const {error: vaultsErr, isLoading, vaults} = useVaults()
    const bakeError = useBakeError()

    useEffect(() => {
        if (isLoading || !vaultsErr) return

        if (isMissingSubscription(vaultsErr as GrpcStatusError)) {
            OpenDialog(
                <InfoDialog
                    title="Subscriptions unavailable"
                    message="Subscriptions are currently not available. Subscribe to Telegram channel @artel_ai"
                    onClose={CloseDialog}
                />
            )
            return;
        }


        bakeError("Failed to load vaults", vaultsErr)
    }, [isLoading, vaultsErr])

    function openEditDialog(vaultId: string) {
        const vault = vaults.find(v => v.id === vaultId)
        if (!vault) return
        const {CloseDialog} = useDialog.getState()
        const {auth} = useUser.getState()
        OpenDialog(
            <ManageVaultDialog
                vault={vault}
                currentUserId={auth.userInfo?.id ?? ""}
                onClose={CloseDialog}
                onDeleted={CloseDialog}
            />
        )
    }

    return (
        <div className={cls.Root}>
            <HeroSegment onCreateClick={() => OpenDialog(<CreateVaultDialog/>)}/>
            <ContentSegment onEditClick={openEditDialog}/>
        </div>
    )
}

function HeroSegment({onCreateClick}: { onCreateClick: () => void }) {
    const {isLoading, vaults} = useVaults()

    return (
        <div className={cls.Hero}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>Workspace</div>
                <h1 className={cls.HeroTitle}>Your vaults</h1>
                <p className={cls.HeroSub}>
                    <b>{isLoading ? "…" : `${vaults.length} ${vaults.length === 1 ? "vault" : "vaults"}`}</b>
                    {" · "}<span>all systems operational</span>
                </p>
            </div>
            <Button variant="primary" onClick={onCreateClick}>
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2"
                     strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                New vault
            </Button>
        </div>
    )
}

function ContentSegment({onEditClick}: { onEditClick: (id: string) => void }) {
    const {isLoading, vaults} = useVaults()

    let loadingState = null

    if (isLoading) {
        loadingState = (
            <p className={cls.Empty}>Loading…</p>
        )
    }

    return (
        <div className={cls.Content}>
            {loadingState}
            {
                !isLoading && <div className={cls.Grid}>
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

    const {create} = useVaultMutations()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

    async function handleCreate() {
        setIsCreating(true)
        const name = vaultName.trim()
        create(name)
            .then(CloseDialog)
            .catch(e => bakeError('Error creating vault', e))
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

                <div className={cls.ModalActions}>
                    <Button variant="primary" onClick={handleCreate} disabled={isCreating}>
                        {isCreating ? "Creating…" : "Create vault"}
                    </Button>
                </div>
            </div>
        </div>
    )
}

