import {useEffect, useState} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/dialogs/ManageVaultDialog/widgets/InviteLinksSection/InviteLinksSection.module.css"
import {useVaultMutations, VaultInviteItem} from "@/app/hooks/Vaults.ts"
import {useDialog} from "@/app/hooks/Dialog.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import InviteRow from "@/dialogs/ManageVaultDialog/components/InviteRow/InviteRow.tsx"
import CreateInviteLinkDialog from "@/dialogs/ManageVaultDialog/components/CreateInviteLinkDialog/CreateInviteLinkDialog.tsx"

interface Props {
    vaultId: string
}

export default function InviteLinksSection({vaultId}: Props) {
    const {listInvites, revokeInvite} = useVaultMutations()
    const {OpenDialog} = useDialog()
    const bakeError = useBakeError()
    const [invites, setInvites] = useState<VaultInviteItem[]>([])

    useEffect(() => {
        listInvites(vaultId)
            .then(setInvites)
            .catch(e => bakeError('Error during vault invites loading', e))
    }, [vaultId])

    function handleRevokeInvite(inviteId: string) {
        return revokeInvite(vaultId, inviteId)
            .then(() => setInvites(prev => prev.map(inv => inv.id === inviteId ? {...inv, revoked: true} : inv)))
            .catch(e => bakeError('Error revoking invite', e))
    }

    function openCreateInviteDialog() {
        OpenDialog(
            <CreateInviteLinkDialog
                vaultId={vaultId}
                onCreated={(inv) => setInvites(prev => [inv, ...prev])}
            />
        )
    }

    return (
        <section className={cls.InviteLinksSectionContainer}>
            <div className={cls.SectionHead}>
                <span className={cls.SectionTitle}>Invite links</span>
                <Button variant="secondary" onClick={openCreateInviteDialog}>
                    + New link
                </Button>
            </div>
            <div className={cls.InviteList}>
                {invites.map(inv => (
                    <InviteRow
                        key={inv.id}
                        invite={inv}
                        onRevoke={() => handleRevokeInvite(inv.id ?? "")}
                    />
                ))}
                {invites.length === 0 && <p className={cls.Empty}>No invite links yet.</p>}
            </div>
        </section>
    )
}
