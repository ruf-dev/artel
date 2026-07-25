import {useEffect, useState} from "react"
import {Button, InfoDialog} from "@vervstack/chures"
import {Tooltip} from "react-tooltip"
import {generatePath, Link, useNavigate} from "react-router-dom"

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
import {AuthAPI} from "@/app/api/artel"
import {Path} from "@/app/routing/Router.tsx"

export default function HomePage() {
    const {OpenDialog, CloseDialog} = useDialog()
    const {error: vaultsErr, isLoading, vaults} = useVaults()
    const bakeError = useBakeError()
    const {auth, isAdmin} = useUser()
    const navigate = useNavigate()

    // Default to "available" while loading — a false-positive briefly-enabled button just
    // falls back on the existing NoCouchDbInstance submit-time error, whereas a
    // false-positive briefly-disabled button would flash wrong on every normal page load.
    const [isCouchAvailable, setIsCouchAvailable] = useState(true)
    // Same rationale as isCouchAvailable above.
    const [isWorkbenchAvailable, setIsWorkbenchAvailable] = useState(true)

    useEffect(() => {
        AuthAPI.GetConfig({}, auth.getInitReq())
            .then(res => {
                setIsCouchAvailable(!!res.isCouchAvailable)
                setIsWorkbenchAvailable(!!res.isWorkbenchAvailable)
            })
    }, [])

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

    function openWorkbench(vaultId: string) {
        navigate(generatePath(Path.WorkbenchPage, {vaultId}))
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
                    <span
                        className={cls.NewVaultButtonWrapper}
                        data-tooltip-id={isCouchAvailable ? undefined : "create-vault-tooltip"}
                    >
                        <Button
                            variant="primary"
                            disabled={!isCouchAvailable}
                            onClick={() => OpenDialog(<CreateVaultDialog/>)}
                        >
                            <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor"
                                 strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                <line x1="12" y1="5" x2="12" y2="19"/>
                                <line x1="5" y1="12" x2="19" y2="12"/>
                            </svg>
                            New vault
                        </Button>
                    </span>
                }
            />
            <ContentSegment
                onEditClick={openEditDialog}
                onWorkbenchClick={openWorkbench}
                isWorkbenchAvailable={isWorkbenchAvailable}
            />
            {!isCouchAvailable && (
                <Tooltip id="create-vault-tooltip" clickable noArrow>
                    {isAdmin
                        ? <>No CouchDB instances configured. <Link to={Path.Admin}>Add one in Admin</Link></>
                        : "New vaults require a CouchDB instance to be configured. Contact an administrator."}
                </Tooltip>
            )}
        </div>
    )
}


