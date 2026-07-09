import {useState} from "react"
import {Button, ConfirmDialog} from "@vervstack/chures"

import cls from "@/components/S3InstancesTab/S3InstanceRow.module.css"
import {S3InstancesAPI, GetS3InstanceResponse} from "@/app/api/artel/s3_instances.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import S3InstanceRowInfo from "@/components/S3InstancesTab/S3InstanceRowInfo.tsx"

export type TestStatus = "idle" | "testing" | "ok" | "fail"

export default function S3InstanceRow({instance, onEdit, onDelete}: {
    instance: GetS3InstanceResponse
    onEdit: (instance: GetS3InstanceResponse) => void
    onDelete: (id: string) => Promise<void>
}) {
    const {auth} = useUser()
    const {OpenDialog, CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const [deleting, setDeleting] = useState(false)
    const [testStatus, setTestStatus] = useState<TestStatus>("idle")

    function handleTest() {
        if (!instance.id) return
        setTestStatus("testing")
        S3InstancesAPI.TestS3Instance({id: instance.id}, auth.getInitReq())
            .then(() => setTestStatus("ok"))
            .catch(err => {
                setTestStatus("fail")
                bakeError("Connection test failed", err)
            })
    }

    function handleDelete() {
        const id = instance.id
        if (!id) return
        OpenDialog(
            <ConfirmDialog
                title="Remove instance"
                message={`Remove "${instance.endpoint}"? This cannot be undone.`}
                danger
                confirmLabel="Remove"
                onClose={CloseDialog}
                onConfirm={() => {
                    setDeleting(true)
                    return onDelete(id)
                        .catch(err => bakeError("Failed to remove instance", err))
                        .finally(() => setDeleting(false))
                }}
            />
        )
    }

    return (
        <div className={cls.S3InstanceRowContainer}>
            <S3InstanceRowInfo instance={instance} testStatus={testStatus}/>
            <div className={cls.RowActions}>
                <Button variant="secondary" onClick={handleTest} disabled={testStatus === "testing"}>
                    {testStatus === "testing" ? "Testing…" : "Test"}
                </Button>
                <Button variant="secondary" onClick={() => onEdit(instance)}>Edit</Button>
                <Button variant="danger" onClick={handleDelete} disabled={deleting}>
                    {deleting ? "Removing…" : "Remove"}
                </Button>
            </div>
        </div>
    )
}
