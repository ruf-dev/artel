import {Button} from "@vervstack/chures"

import cls from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserRow/ArtelUserRow.module.css"
import {ArtelUserEntry} from "@/app/api/artel/admin_users.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import ArtelUserDetailDialog from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/ArtelUserDetailDialog.tsx"

interface ArtelUserRowProps {
    user: ArtelUserEntry
}

export default function ArtelUserRow({user}: ArtelUserRowProps) {
    const {OpenDialog} = useDialog()

    function openDetail() {
        if (!user.userId) return
        OpenDialog(<ArtelUserDetailDialog userId={user.userId} />)
    }

    return (
        <div className={cls.ArtelUserRowContainer}>
            <div className={cls.RowInfo}>
                <span className={cls.RowUrl}>{user.username || "—"}</span>
                <span className={cls.RowMeta}>{user.email || "no email"}</span>
            </div>
            <div className={cls.RowActions}>
                <Button variant="secondary" onClick={openDetail}>
                    Details
                </Button>
            </div>
        </div>
    )
}
