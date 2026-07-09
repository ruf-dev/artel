import {useEffect} from "react"
import {InfoDialog} from "@vervstack/chures"

import cls from "@/pages/home/HomePage.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useVaults} from "@/app/hooks/Vaults.ts"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import ManageVaultDialog from "@/dialogs/ManageVaultDialog/ManageVaultDialog.tsx"
import type {GrpcStatusError} from "@/processes/grpcErrors.ts";
import {isMissingSubscription} from "@/processes/UserErrors.ts";
import HeroSegment from "@/pages/home/components/HeroSegment/HeroSegment.tsx"
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
            <HeroSegment onCreateClick={() => OpenDialog(<CreateVaultDialog/>)}/>
            <ContentSegment onEditClick={openEditDialog}/>
        </div>
    )
}


