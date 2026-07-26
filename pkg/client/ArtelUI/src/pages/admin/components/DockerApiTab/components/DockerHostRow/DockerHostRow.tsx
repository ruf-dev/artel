import {useState} from "react"
import {Button, ConfirmDialog} from "@vervstack/chures"

import cls from "@/pages/admin/components/DockerApiTab/components/DockerHostRow/DockerHostRow.module.css"
import {GetDockerHostResponse} from "@/app/api/artel/docker_hosts.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"

interface DockerHostRowProps {
    host: GetDockerHostResponse
    onEdit: (host: GetDockerHostResponse) => void
    onDelete: (id: string) => Promise<void>
}

export default function DockerHostRow({host, onEdit, onDelete}: DockerHostRowProps) {
    const [deleting, setDeleting] = useState(false)
    const {OpenDialog, CloseDialog} = useDialog()
    const bakeError = useBakeError()

    function handleDelete() {
        const id = host.id
        if (!id) return
        OpenDialog(
            <ConfirmDialog
                title="Remove Docker host"
                message={`Remove "${host.url}"? This cannot be undone.`}
                danger
                confirmLabel="Remove"
                onClose={CloseDialog}
                onConfirm={async () => {
                    setDeleting(true)
                    try {
                        await onDelete(id)
                    } catch (err) {
                        bakeError("Failed to remove Docker host", err)
                    } finally {
                        setDeleting(false)
                    }
                }}
            />
        )
    }

    return (
        <div className={cls.DockerHostRowContainer}>
            <div className={cls.RowInfo}>
                <span className={cls.RowUrl}>{host.url}</span>
                <span className={cls.RowMeta}>added {host.createdAt?.slice(0, 10)}</span>
            </div>
            <div className={cls.RowActions}>
                <Button variant="secondary" onClick={() => onEdit(host)}>Edit</Button>
                <Button variant="danger" onClick={handleDelete} disabled={deleting}>
                    {deleting ? "Removing…" : "Remove"}
                </Button>
            </div>
        </div>
    )
}
