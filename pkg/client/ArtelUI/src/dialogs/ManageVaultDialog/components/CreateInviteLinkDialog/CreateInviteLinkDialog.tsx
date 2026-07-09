import {useState} from "react"
import {Button, ModalClose} from "@vervstack/chures"

import cls from "@/dialogs/ManageVaultDialog/components/CreateInviteLinkDialog/CreateInviteLinkDialog.module.css"
import {useVaultMutations, VaultInviteItem} from "@/app/hooks/Vaults.ts"
import {useDialog} from "@/app/hooks/Dialog.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import RoleOption from "@/dialogs/ManageVaultDialog/components/RoleOption/RoleOption.tsx"

interface Props {
    vaultId: string
    onCreated: (inv: VaultInviteItem) => void
}

export default function CreateInviteLinkDialog({vaultId, onCreated}: Props) {
    const {createInvite} = useVaultMutations()
    const {CloseDialog} = useDialog()
    const [role, setRole] = useState<"reader" | "maintainer">("reader")
    const [creating, setCreating] = useState(false)
    const [created, setCreated] = useState<VaultInviteItem | null>(null)
    const [copied, setCopied] = useState(false)
    const bakeError = useBakeError()

    function handleCreateInvite() {
        setCreating(true)
        createInvite(vaultId, role)
            .then(inv => {
                setCreated(inv)
                onCreated(inv)
            })
            .catch(e => bakeError('Error creating invite', e))
            .finally(() => setCreating(false))
    }

    function handleCopy() {
        if (!created?.token) return
        navigator.clipboard.writeText(`${window.location.origin}/join/${created.token}`).then(() => {
            setCopied(true)
            setTimeout(() => setCopied(false), 2000)
        })
    }

    return (
        <div className={cls.CreateInviteLinkDialogContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>New invite link</h2>
                <ModalClose onClick={CloseDialog} className={cls.ModalClose}/>
            </div>

            {!created ? (
                <>
                    <p className={cls.ModalSub}>Choose what the invited user can do.</p>
                    <div className={cls.RoleSelector}>
                        <RoleOption
                            selected={role === "reader"}
                            label="Reader"
                            desc="Can view and sync the vault, cannot write."
                            onSelect={() => setRole("reader")}
                        />
                        <RoleOption
                            selected={role === "maintainer"}
                            label="Maintainer"
                            desc="Can read and write to the vault."
                            onSelect={() => setRole("maintainer")}
                        />
                    </div>
                    <div className={cls.ModalFooter}>
                        <Button variant="ghost" onClick={CloseDialog}>Cancel</Button>
                        <Button variant="primary" onClick={handleCreateInvite} disabled={creating}>
                            {creating ? "Creating…" : "Create link"}
                        </Button>
                    </div>
                </>
            ) : (
                <>
                    <p className={cls.ModalSub}>Share this link. Anyone with it can join as <strong>{role}</strong>.</p>
                    <div className={cls.LinkBox}>
                        <span className={cls.LinkText}>{window.location.origin}/join/{created.token}</span>
                    </div>
                    <div className={cls.ModalFooter}>
                        <Button variant="ghost" onClick={CloseDialog}>Done</Button>
                        <Button variant="primary" onClick={handleCopy}>
                            {copied ? "Copied!" : "Copy link"}
                        </Button>
                    </div>
                </>
            )}
        </div>
    )
}
