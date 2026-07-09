import {useState, useEffect, useCallback} from "react"
import {Button, ConfirmDialog} from "@vervstack/chures"

import cls from "@/widgets/S3InstancesTab/S3InstancesTab.module.css"
import {S3InstancesAPI, GetS3InstanceResponse} from "@/app/api/artel/s3_instances.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import S3InstanceFormDialog from "@/components/S3InstanceFormDialog/S3InstanceFormDialog.tsx"

type TestStatus = "idle" | "testing" | "ok" | "fail"

export default function S3InstancesTab() {
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const bakeError = useBakeError()
    const [instances, setInstances] = useState<GetS3InstanceResponse[]>([])
    const [loading, setLoading] = useState(true)

    const loadInstances = useCallback(() => {
        setLoading(true)
        S3InstancesAPI.ListS3Instances({}, auth.getInitReq())
            .then(res => setInstances(res.instances ?? []))
            .catch(err => bakeError("Failed to load S3 instances", err))
            .finally(() => setLoading(false))
    }, [auth, bakeError])

    useEffect(() => {
        loadInstances()
    }, [loadInstances])

    function handleDelete(id: string) {
        return S3InstancesAPI.DeleteS3Instance({id}, auth.getInitReq()).then(() => loadInstances())
    }

    function openAddDialog() {
        OpenDialog(<S3InstanceFormDialog onSaved={loadInstances}/>)
    }

    function openEditDialog(instance: GetS3InstanceResponse) {
        OpenDialog(<S3InstanceFormDialog initial={instance} onSaved={loadInstances}/>)
    }

    return (
        <div className={cls.S3InstancesTabContainer}>
            <S3InstancesActionBar count={instances.length} onAddClick={openAddDialog}/>
            <S3InstanceList
                instances={instances}
                loading={loading}
                onEdit={openEditDialog}
                onDelete={handleDelete}
            />
        </div>
    )
}

function S3InstancesActionBar({count, onAddClick}: {count: number; onAddClick: () => void}) {
    return (
        <div className={cls.ActionBar}>
            <p className={cls.HeroSub}>
                <b>{count} {count === 1 ? "instance" : "instances"}</b>
                {" · "}<span>manage S3-compatible storage backends</span>
            </p>
            <Button variant="primary" onClick={onAddClick}>Add instance</Button>
        </div>
    )
}

function S3InstanceList({instances, loading, onEdit, onDelete}: {
    instances: GetS3InstanceResponse[]
    loading: boolean
    onEdit: (instance: GetS3InstanceResponse) => void
    onDelete: (id: string) => Promise<void>
}) {
    if (loading) {
        return <p className={cls.Empty}>Loading…</p>
    }
    if (instances.length === 0) {
        return <p className={cls.Empty}>No S3 instances yet. Add one to get started.</p>
    }
    return (
        <div className={cls.List}>
            {instances.map(instance => (
                <S3InstanceRow key={instance.id} instance={instance} onEdit={onEdit} onDelete={onDelete}/>
            ))}
        </div>
    )
}

function S3InstanceRow({instance, onEdit, onDelete}: {
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
        <div className={cls.RowContainer}>
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

function S3InstanceRowInfo({instance, testStatus}: {instance: GetS3InstanceResponse; testStatus: TestStatus}) {
    return (
        <div className={cls.RowInfo}>
            <span className={cls.RowEndpoint}>{instance.endpoint}</span>
            <span className={cls.RowMeta}>
                {instance.region || "no region"} · added {instance.createdAt?.slice(0, 10)}
            </span>
            <div className={cls.BadgeRow}>
                {instance.useSsl && <span className={cls.Badge}>SSL</span>}
                {instance.pathStyle && <span className={cls.Badge}>Path-style</span>}
                {testStatus === "ok" && <span className={cls.BadgeOk}>Connected</span>}
                {testStatus === "fail" && <span className={cls.BadgeFail}>Failed</span>}
            </div>
        </div>
    )
}
