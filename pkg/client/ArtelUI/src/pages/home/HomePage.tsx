import {useEffect} from "react"
import {Button, InfoDialog} from "@vervstack/chures"

import cls from "@/pages/home/HomePage.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useVaults} from "@/app/hooks/Vaults.ts"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import ManageVaultDialog from "@/dialogs/ManageVaultDialog/ManageVaultDialog.tsx"
import type {GrpcStatusError} from "@/processes/grpcErrors.ts";
import {isMissingSubscription} from "@/processes/UserErrors.ts";
import HeroSegment from "@/components/HeroSegment/HeroSegment.tsx"
import ContentSegment from "@/pages/home/components/ContentSegment/ContentSegment.tsx"
import CreateVaultDialog from "@/pages/home/components/CreateVaultDialog/CreateVaultDialog.tsx"

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
            <HeroSegment
                eyebrow="Workspace"
                title="Your vaults"
                subtitle={
                    <>
                        <b>{isLoading ? "…" : `${vaults.length} ${vaults.length === 1 ? "vault" : "vaults"}`}</b>
                        {" · "}<span>all systems operational</span>
                    </>
                }
                action={
                    <Button variant="primary" onClick={() => OpenDialog(<CreateVaultDialog/>)}>
                        <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor"
                             strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <line x1="12" y1="5" x2="12" y2="19"/>
                            <line x1="5" y1="12" x2="19" y2="12"/>
                        </svg>
                        New vault
                    </Button>
                }
            />
            <ContentSegment onEditClick={openEditDialog}/>
        </div>
    )
}


