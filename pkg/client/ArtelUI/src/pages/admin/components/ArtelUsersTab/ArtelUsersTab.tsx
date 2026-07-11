import {useState, useEffect, useCallback} from "react"

import SearchInput from "@/components/SearchInput/SearchInput.tsx"
import cls from "@/pages/admin/components/ArtelUsersTab/ArtelUsersTab.module.css"
import {AdminUsersAPI, ArtelUserEntry} from "@/app/api/artel/admin_users.pb.ts"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import ArtelUserList from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserList/ArtelUserList.tsx"

const LIMIT = 50

export default function ArtelUsersTab() {
    const {auth} = useUser()
    const [search, setSearch] = useState("")
    const [users, setUsers] = useState<ArtelUserEntry[]>([])
    const [total, setTotal] = useState(0)
    const [loading, setLoading] = useState(true)
    const bakeError = useBakeError()

    const loadUsers = useCallback(async () => {
        setLoading(true)
        try {
            const res = await AdminUsersAPI.ListArtelUsers({search, limit: LIMIT, offset: 0}, auth.getInitReq())
            setUsers(res.users ?? [])
            setTotal(res.total ? Number(res.total) : 0)
        } catch (err) {
            bakeError("Failed to load users", err)
        } finally {
            setLoading(false)
        }
    }, [auth, search])

    useEffect(() => { void loadUsers() }, [loadUsers])

    return (
        <div className={cls.ArtelUsersTabContainer}>
            <SearchInput
                className={cls.Field}
                placeholder="Search by username or email…"
                value={search}
                setValue={setSearch}
            />
            <ArtelUserList users={users} loading={loading} total={total} />
        </div>
    )
}
