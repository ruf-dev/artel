import {CouchUserEntry} from "@/app/api/artel/admin_couch.pb.ts"
import UserRow from "@/pages/admin/components/UsersTab/components/UserRow/UserRow.tsx"
import cls from "@/pages/admin/components/UsersTab/components/UserList/UserList.module.css"

interface UserListProps {
    instanceId: string
    users: CouchUserEntry[]
    loading: boolean
    onRefresh: () => Promise<void>
}

export default function UserList({instanceId, users, loading, onRefresh}: UserListProps) {
    if (!instanceId) {
        return <p className={cls.Empty}>Select an instance above to view its users.</p>
    }
    if (loading) {
        return <p className={cls.Empty}>Loading…</p>
    }
    if (users.length === 0) {
        return <p className={cls.Empty}>No users found in this instance.</p>
    }
    return (
        <>
            {users.map(user => (
                <UserRow key={user.name} instanceId={instanceId} user={user} onRefresh={onRefresh} />
            ))}
        </>
    )
}
