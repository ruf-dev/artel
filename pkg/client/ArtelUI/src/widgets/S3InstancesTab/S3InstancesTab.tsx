import {useState, useEffect, useCallback} from "react"

import cls from "@/widgets/S3InstancesTab/S3InstancesTab.module.css"
import {S3InstancesAPI, GetS3InstanceResponse} from "@/app/api/artel/s3_instances.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import S3InstanceFormDialog from "@/components/S3InstanceFormDialog/S3InstanceFormDialog.tsx"
import S3InstancesActionBar from "@/components/S3InstancesTab/S3InstancesActionBar.tsx"
import S3InstanceList from "@/components/S3InstancesTab/S3InstanceList.tsx"

export default function S3InstancesTab() {
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const bakeError = useBakeError()
    const [instances, setInstances] = useState<GetS3InstanceResponse[]>([])
    const [loading, setLoading] = useState(true)

    const loadInstances = useCallback(() => {
        setLoading(true)
        S3InstancesAPI.ListS3Instances({}, auth.getInitReq())
            .then(res => setInstances(res.instances ?? []))
            .catch(err => bakeError("Failed to load S3 instances", err))
            .finally(() => setLoading(false))
    }, [auth, bakeError])

    useEffect(() => {
        loadInstances()
    }, [loadInstances])

    function handleDelete(id: string) {
        return S3InstancesAPI.DeleteS3Instance({id}, auth.getInitReq()).then(() => loadInstances())
    }

    function openAddDialog() {
        OpenDialog(<S3InstanceFormDialog onSaved={loadInstances}/>)
    }

    function openEditDialog(instance: GetS3InstanceResponse) {
        OpenDialog(<S3InstanceFormDialog initial={instance} onSaved={loadInstances}/>)
    }

    return (
        <div className={cls.S3InstancesTabContainer}>
            <S3InstancesActionBar count={instances.length} onAddClick={openAddDialog}/>
            <S3InstanceList
                instances={instances}
                loading={loading}
                onEdit={openEditDialog}
                onDelete={handleDelete}
            />
        </div>
    )
}
