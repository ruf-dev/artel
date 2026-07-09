import {useState, useEffect, useCallback} from "react"

import cls from "@/pages/admin/components/UsersTab/UsersTab.module.css"
import {AdminCouchAPI, CouchUserEntry} from "@/app/api/artel/admin_couch.pb.ts"
import {CouchInstancesAPI, GetCouchInstanceResponse} from "@/app/api/artel/couch_instances.pb.ts"
import useUser from "@/hooks/user/User.ts"
import InstanceSelector from "@/pages/admin/components/UsersTab/components/InstanceSelector/InstanceSelector.tsx"
import UserList from "@/pages/admin/components/UsersTab/components/UserList/UserList.tsx"

export default function UsersTab() {
    const {auth} = useUser()
    const [instances, setInstances] = useState<GetCouchInstanceResponse[]>([])
    const [selectedInstanceId, setSelectedInstanceId] = useState("")
    const [users, setUsers] = useState<CouchUserEntry[]>([])
    const [usersLoading, setUsersLoading] = useState(false)

    useEffect(() => {
        async function loadInstances() {
            const res = await CouchInstancesAPI.ListCouchInstances({}, auth.getInitReq())
            setInstances(res.instances ?? [])
        }
        void loadInstances()
    }, [auth])

    const loadUsers = useCallback(async () => {
        if (!selectedInstanceId) {
            setUsers([])
            return
        }
        setUsersLoading(true)
        try {
            const res = await AdminCouchAPI.ListCouchUsers({instanceId: selectedInstanceId}, auth.getInitReq())
            setUsers(res.users ?? [])
        } finally {
            setUsersLoading(false)
        }
    }, [auth, selectedInstanceId])

    useEffect(() => { void loadUsers() }, [loadUsers])

    return (
        <div className={cls.UsersTabContainer}>
            <InstanceSelector instances={instances} value={selectedInstanceId} onChange={setSelectedInstanceId} />
            <UserList
                instanceId={selectedInstanceId}
                users={users}
                loading={usersLoading}
                onRefresh={loadUsers}
            />
        </div>
    )
}
