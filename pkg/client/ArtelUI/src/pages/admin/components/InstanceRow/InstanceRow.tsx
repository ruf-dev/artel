import {useState} from "react"
import {Button, ConfirmDialog} from "@vervstack/chures"

import cls from "@/pages/admin/components/InstanceRow/InstanceRow.module.css"
import {GetCouchInstanceResponse} from "@/app/api/artel/couch_instances.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"

interface InstanceRowProps {
    instance: GetCouchInstanceResponse
    onEdit: (instance: GetCouchInstanceResponse) => void
    onDelete: (id: string) => Promise<void>
}

export default function InstanceRow({instance, onEdit, onDelete}: InstanceRowProps) {
    const [deleting, setDeleting] = useState(false)
    const {OpenDialog, CloseDialog} = useDialog()
    const bakeError = useBakeError()

    function handleDelete() {
        const id = instance.id
        if (!id) return
        OpenDialog(
            <ConfirmDialog
                title="Remove instance"
                message={`Remove "${instance.url}"? This cannot be undone.`}
                danger
                confirmLabel="Remove"
                onClose={CloseDialog}
                onConfirm={async () => {
                    setDeleting(true)
                    try {
                        await onDelete(id)
                    } catch (err) {
                        bakeError("Failed to remove instance", err)
                    } finally {
                        setDeleting(false)
                    }
                }}
            />
        )
    }

    return (
        <div className={cls.InstanceRowContainer}>
            <div className={cls.RowInfo}>
                <span className={cls.RowUrl}>{instance.url}</span>
                <span className={cls.RowMeta}>{instance.username} · added {instance.createdAt?.slice(0, 10)}</span>
            </div>
            <div className={cls.RowActions}>
                <Button variant="secondary" onClick={() => onEdit(instance)}>Edit</Button>
                <Button variant="danger" onClick={handleDelete} disabled={deleting}>
                    {deleting ? "Removing…" : "Remove"}
                </Button>
            </div>
        </div>
    )
}
