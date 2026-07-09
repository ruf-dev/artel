import {Button} from "@vervstack/chures"

import cls from "@/pages/admin/components/UsersTab/components/ManageAccessDialog/components/DbAccessRow/DbAccessRow.module.css"

interface DbAccessRowProps {
    db: string
    granted: boolean
    busy: boolean
    onToggle: (db: string, grant: boolean) => Promise<void>
}

export default function DbAccessRow({db, granted, busy, onToggle}: DbAccessRowProps) {
    return (
        <div className={cls.DbAccessRowContainer}>
            <span className={cls.DbAccessName}>{db}</span>
            <Button
                variant={granted ? "danger" : "primary"}
                disabled={busy}
                onClick={() => onToggle(db, !granted)}
            >
                {busy ? "…" : granted ? "Revoke" : "Grant"}
            </Button>
        </div>
    )
}
