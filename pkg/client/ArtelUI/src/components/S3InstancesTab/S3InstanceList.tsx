import cls from "@/components/S3InstancesTab/S3InstanceList.module.css"
import {GetS3InstanceResponse} from "@/app/api/artel/s3_instances.pb.ts"
import S3InstanceRow from "@/components/S3InstancesTab/S3InstanceRow.tsx"

export default function S3InstanceList({instances, loading, onEdit, onDelete}: {
    instances: GetS3InstanceResponse[]
    loading: boolean
    onEdit: (instance: GetS3InstanceResponse) => void
    onDelete: (id: string) => Promise<void>
}) {
    if (loading) {
        return <p className={cls.Empty}>Loading…</p>
    }
    if (instances.length === 0) {
        return <p className={cls.Empty}>No S3 instances yet. Add one to get started.</p>
    }
    return (
        <div className={cls.S3InstanceListContainer}>
            {instances.map(instance => (
                <S3InstanceRow key={instance.id} instance={instance} onEdit={onEdit} onDelete={onDelete}/>
            ))}
        </div>
    )
}
