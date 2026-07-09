import {useState, useEffect} from "react"
import {Button, ConfirmDialog} from "@vervstack/chures"

import cls from "@/components/ManageVaultS3Section/ManageVaultS3Section.module.css"
import {VaultItem, VaultsAPI} from "@/app/api/artel/vaults.pb.ts"
import {S3InstancesAPI, GetS3InstanceResponse} from "@/app/api/artel/s3_instances.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import type {S3LinkPatch} from "@/components/ManageVaultS3Section/ManageVaultS3Section"

interface Props {
    vault: VaultItem
    onLinked: (patch: S3LinkPatch) => void
}

export default function LinkedView({vault, onLinked}: Props) {
    const {auth} = useUser()
    const {OpenDialog, CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const [instance, setInstance] = useState<GetS3InstanceResponse | null>(null)
    const [unlinking, setUnlinking] = useState(false)

    useEffect(() => {
        if (!vault.s3InstanceId) return
        S3InstancesAPI.GetS3Instance({id: vault.s3InstanceId}, auth.getInitReq())
            .then(setInstance)
            .catch(err => bakeError("Failed to load S3 instance", err))
    }, [vault.s3InstanceId, auth, bakeError])

    function handleUnlink() {
        const vaultId = vault.id
        if (!vaultId) return
        OpenDialog(
            <ConfirmDialog
                title="Unlink S3 bucket"
                message="Unlink this vault's S3 bucket? Attachments will stop syncing to S3 until relinked."
                danger
                confirmLabel="Unlink"
                onClose={CloseDialog}
                onConfirm={() => {
                    setUnlinking(true)
                    return VaultsAPI.UnlinkS3Bucket({vaultId}, auth.getInitReq())
                        .then(() => onLinked({s3InstanceId: "", s3BucketName: ""}))
                        .catch(err => bakeError("Failed to unlink bucket", err))
                        .finally(() => setUnlinking(false))
                }}
            />
        )
    }

    return (
        <div className={cls.LinkedRow}>
            <div className={cls.LinkedInfo}>
                <span className={cls.LinkedEndpoint}>{instance?.endpoint ?? "…"}</span>
                <span className={cls.LinkedBucket}>{vault.s3BucketName}</span>
            </div>
            <Button variant="danger" onClick={handleUnlink} disabled={unlinking}>
                {unlinking ? "Unlinking…" : "Unlink"}
            </Button>
        </div>
    )
}
