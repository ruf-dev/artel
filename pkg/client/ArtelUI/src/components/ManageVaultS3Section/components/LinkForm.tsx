import {useState, useEffect} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/components/ManageVaultS3Section/ManageVaultS3Section.module.css"
import Input from "@/components/atoms/Input/Input.tsx"
import {VaultItem, VaultsAPI} from "@/app/api/artel/vaults.pb.ts"
import {S3InstancesAPI, GetS3InstanceResponse} from "@/app/api/artel/s3_instances.pb.ts"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import type {S3LinkPatch} from "@/components/ManageVaultS3Section/ManageVaultS3Section"

interface Props {
    vault: VaultItem
    onLinked: (patch: S3LinkPatch) => void
}

export default function LinkForm({vault, onLinked}: Props) {
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
            <Input
                label="Bucket name"
                value={bucketName}
                setValue={setBucketName}
                disabled={linking}
            />
            <Button variant="primary" onClick={handleLink} disabled={linking || !selectedId || !bucketName}>
                {linking ? "Linking…" : "Link bucket"}
            </Button>
        </div>
    )
}
