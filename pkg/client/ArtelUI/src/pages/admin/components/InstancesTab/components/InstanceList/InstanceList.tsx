import {GetCouchInstanceResponse} from "@/app/api/artel/couch_instances.pb.ts"
import InstanceRow from "@/pages/admin/components/InstanceRow/InstanceRow.tsx"
import cls from "@/pages/admin/components/InstancesTab/components/InstanceList/InstanceList.module.css"

interface InstanceListProps {
    instances: GetCouchInstanceResponse[]
    loading: boolean
    onEdit: (instance: GetCouchInstanceResponse) => void
    onDelete: (id: string) => Promise<void>
}

export default function InstanceList({instances, loading, onEdit, onDelete}: InstanceListProps) {
    if (loading) {
        return <p className={cls.Empty}>Loading…</p>
    }
    if (instances.length === 0) {
        return <p className={cls.Empty}>No CouchDB instances yet. Add one to get started.</p>
    }
    return (
        <>
            {instances.map(instance => (
                <InstanceRow key={instance.id} instance={instance} onEdit={onEdit} onDelete={onDelete} />
            ))}
        </>
    )
}
