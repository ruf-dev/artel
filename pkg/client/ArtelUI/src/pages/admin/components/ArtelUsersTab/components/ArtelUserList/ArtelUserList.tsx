import {ArtelUserEntry} from "@/app/api/artel/admin_users.pb.ts"
import ArtelUserRow from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserRow/ArtelUserRow.tsx"
import cls from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserList/ArtelUserList.module.css"

interface ArtelUserListProps {
    users: ArtelUserEntry[]
    loading: boolean
    total: number
}

export default function ArtelUserList({users, loading, total}: ArtelUserListProps) {
    if (loading) return <p className={cls.Empty}>Loading…</p>
    if (users.length === 0) return <p className={cls.Empty}>No users found.</p>
    return (
        <>
            <p className={cls.Empty}>{total} user{total !== 1 ? "s" : ""} total</p>
            {users.map(u => (
                <ArtelUserRow key={u.userId} user={u} />
            ))}
        </>
    )
}
