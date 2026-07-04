import {useState, useEffect} from "react"

import cls from "@/components/ManageVaultS3Section/ManageVaultS3Section.module.css"

import {VaultItem, VaultsAPI} from "@/app/api/artel/vaults.pb.ts"
import {S3InstancesAPI, GetS3InstanceResponse} from "@/app/api/artel/s3_instances.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"

import Button from "@/components/shared/Button/Button.tsx"
import ConfirmDialog from "@/components/ConfirmDialog/ConfirmDialog.tsx"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"

export interface S3LinkPatch {
    s3InstanceId?: string
    s3BucketName?: string
}

interface Props {
    vault: VaultItem
    onLinked: (patch: S3LinkPatch) => void
}

export default function ManageVaultS3Section({vault, onLinked}: Props) {
    const linked = !!vault.s3InstanceId

    return (
        <section className={cls.SectionContainer}>
            <div className={cls.SectionHead}>
                <span className={cls.SectionTitle}>S3 storage</span>
            </div>
            {linked
                ? <LinkedView vault={vault} onLinked={onLinked}/>
                : <LinkForm vault={vault} onLinked={onLinked}/>
            }
        </section>
    )
}

function LinkedView({vault, onLinked}: Props) {
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
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

function LinkForm({vault, onLinked}: Props) {
    const {auth} = useUser()
    const bakeError = useBakeError()
    const [instances, setInstances] = useState<GetS3InstanceResponse[]>([])
    const [loading, setLoading] = useState(true)
    const [selectedId, setSelectedId] = useState("")
    const [bucketName, setBucketName] = useState("")
    const [linking, setLinking] = useState(false)

    useEffect(() => {
        S3InstancesAPI.ListS3Instances({}, auth.getInitReq())
            .then(res => setInstances(res.instances ?? []))
            .catch(err => bakeError("Failed to load S3 instances", err))
            .finally(() => setLoading(false))
    }, [auth, bakeError])

    function handleLink() {
        const vaultId = vault.id
        if (!vaultId || !selectedId || !bucketName) return
        setLinking(true)
        VaultsAPI.LinkS3Bucket({vaultId, s3InstanceId: selectedId, bucketName}, auth.getInitReq())
            .then(() => onLinked({s3InstanceId: selectedId, s3BucketName: bucketName}))
            .catch(err => bakeError("Failed to link bucket", err))
            .finally(() => setLinking(false))
    }

    if (loading) {
        return <p className={cls.Empty}>Loading…</p>
    }

    if (instances.length === 0) {
        return <p className={cls.Empty}>No S3 instances configured. Add one in Admin → S3 storage first.</p>
    }

    return (
        <div className={cls.LinkFormContainer}>
            <div className={cls.OptionList}>
                {instances.map(inst => (
                    <SelectOption
                        key={inst.id}
                        label={inst.endpoint ?? ""}
                        selected={selectedId === inst.id}
                        onSelect={() => setSelectedId(inst.id ?? "")}
                    />
                ))}
            </div>
            <input
                className={cls.Input}
                placeholder="bucket name"
                value={bucketName}
                onChange={e => setBucketName(e.target.value)}
                disabled={linking}
            />
            <Button variant="primary" onClick={handleLink} disabled={linking || !selectedId || !bucketName}>
                {linking ? "Linking…" : "Link bucket"}
            </Button>
        </div>
    )
}
