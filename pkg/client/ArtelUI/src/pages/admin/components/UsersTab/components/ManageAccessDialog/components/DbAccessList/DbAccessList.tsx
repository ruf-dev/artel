import DbAccessRow from "@/pages/admin/components/UsersTab/components/ManageAccessDialog/components/DbAccessRow/DbAccessRow.tsx"
import cls from "@/pages/admin/components/UsersTab/components/ManageAccessDialog/components/DbAccessList/DbAccessList.module.css"

interface DbAccessListProps {
    databases: string[]
    grantedSet: Set<string>
    busyDb: string | null
    onToggle: (db: string, grant: boolean) => Promise<void>
}

export default function DbAccessList({databases, grantedSet, busyDb, onToggle}: DbAccessListProps) {
    return (
        <div className={cls.DbAccessListContainer}>
            {databases.map(db => (
                <DbAccessRow
                    key={db}
                    db={db}
                    granted={grantedSet.has(db)}
                    busy={busyDb === db}
                    onToggle={onToggle}
                />
            ))}
        </div>
    )
}
