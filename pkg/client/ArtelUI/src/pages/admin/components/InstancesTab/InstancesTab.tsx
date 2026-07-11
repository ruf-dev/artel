import {useState, useEffect, useCallback} from "react"

import cls from "@/pages/admin/components/InstancesTab/InstancesTab.module.css"
import {CouchInstancesAPI, GetCouchInstanceResponse} from "@/app/api/artel/couch_instances.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import useUser from "@/hooks/user/User.ts"
import InstancesActionBar
    from "@/pages/admin/components/InstancesTab/components/InstancesActionBar/InstancesActionBar.tsx"
import InstanceList from "@/pages/admin/components/InstancesTab/components/InstanceList/InstanceList.tsx"
import InstanceFormDialog
    from "@/pages/admin/components/InstancesTab/components/InstanceFormDialog/InstanceFormDialog.tsx"

export default function InstancesTab() {
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const [instances, setInstances] = useState<GetCouchInstanceResponse[]>([])
    const [loading, setLoading] = useState(true)

    const loadInstances = useCallback(async () => {
        setLoading(true)
        try {
            const res = await CouchInstancesAPI.ListCouchInstances({}, auth.getInitReq())
            setInstances(res.instances ?? [])
        } finally {
            setLoading(false)
        }
    }, [auth])

    useEffect(() => { void loadInstances() }, [loadInstances])

    async function handleDelete(id: string) {
        await CouchInstancesAPI.DeleteCouchInstance({id}, auth.getInitReq())
        await loadInstances()
    }

    function openAddDialog() {
        OpenDialog(
            <InstanceFormDialog
                onSave={async (url, username, password) => {
                    await CouchInstancesAPI.RegisterCouchInstance({url, username, password}, auth.getInitReq())
                    await loadInstances()
                }}
            />
        )
    }

    function openEditDialog(instance: GetCouchInstanceResponse) {
        OpenDialog(
            <InstanceFormDialog
                initial={instance}
                onSave={async (url, username, password) => {
                    await CouchInstancesAPI.UpdateCouchInstance(
                        {id: instance.id, url, username, password},
                        auth.getInitReq()
                    )
                    await loadInstances()
                }}
            />
        )
    }

    return (
        <div className={cls.InstancesTabContainer}>
            <InstancesActionBar count={instances.length} onAddClick={openAddDialog} />
            <InstanceList instances={instances} loading={loading} onEdit={openEditDialog} onDelete={handleDelete} />
        </div>
    )
}
