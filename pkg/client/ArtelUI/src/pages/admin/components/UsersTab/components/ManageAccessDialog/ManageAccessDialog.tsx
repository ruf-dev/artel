import {useState, useEffect, useCallback} from "react"
import {Button, ModalClose} from "@vervstack/chures"

import cls from "@/pages/admin/components/UsersTab/components/ManageAccessDialog/ManageAccessDialog.module.css"
import {AdminCouchAPI} from "@/app/api/artel/admin_couch.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import DbAccessList from "@/pages/admin/components/UsersTab/components/ManageAccessDialog/components/DbAccessList/DbAccessList.tsx"

interface ManageAccessDialogProps {
    instanceId: string
    username: string
}

export default function ManageAccessDialog({instanceId, username}: ManageAccessDialogProps) {
    const {auth} = useUser()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const [databases, setDatabases] = useState<string[]>([])
    const [grantedSet, setGrantedSet] = useState<Set<string>>(new Set())
    const [loading, setLoading] = useState(true)
    const [busyDb, setBusyDb] = useState<string | null>(null)

    const loadAccess = useCallback(async () => {
        setLoading(true)
        try {
            const [dbRes, accessRes] = await Promise.all([
                AdminCouchAPI.ListCouchDatabases({instanceId}, auth.getInitReq()),
                AdminCouchAPI.GetUserDatabaseAccess({instanceId, username}, auth.getInitReq()),
            ])
            const dbs = (dbRes.databases ?? []).filter(db => !db.startsWith("_"))
            setDatabases(dbs)
            setGrantedSet(new Set(accessRes.databases ?? []))
        } finally {
            setLoading(false)
        }
    }, [auth, instanceId, username])

    useEffect(() => { void loadAccess() }, [loadAccess])

    async function handleToggle(db: string, grant: boolean) {
        setBusyDb(db)
        try {
            if (grant) {
                await AdminCouchAPI.GrantDatabaseAccess({instanceId, dbName: db, username}, auth.getInitReq())
            } else {
                await AdminCouchAPI.RevokeDatabaseAccess({instanceId, dbName: db, username}, auth.getInitReq())
            }
            setGrantedSet(prev => {
                const next = new Set(prev)
                if (grant) {
                    next.add(db)
                } else {
                    next.delete(db)
                }
                return next
            })
        } catch (err) {
            bakeError(grant ? "Failed to grant access" : "Failed to revoke access", err)
        } finally {
            setBusyDb(null)
        }
    }

    return (
        <div className={cls.ManageAccessDialogContainer} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>Manage DB access</h2>
                <ModalClose onClick={CloseDialog} className={cls.ModalClose}/>
            </div>
            <p className={cls.ModalSubtitle}>User: <b>{username}</b></p>
            {loading ? (
                <p className={cls.Empty}>Loading…</p>
            ) : databases.length === 0 ? (
                <p className={cls.Empty}>No databases found.</p>
            ) : (
                <DbAccessList
                    databases={databases}
                    grantedSet={grantedSet}
                    busyDb={busyDb}
                    onToggle={handleToggle}
                />
            )}
            <div className={cls.ModalActions}>
                <Button variant="primary" onClick={CloseDialog}>Done</Button>
            </div>
        </div>
    )
}
