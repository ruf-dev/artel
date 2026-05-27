import {useState, useEffect, useCallback} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/admin/AdminPage.module.css"

import {CouchInstancesAPI, GetCouchInstanceResponse} from "@/app/api/artel/couch_instances.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import useUser from "@/hooks/user/User.ts"

import FormField from "@/components/FormField/FormField.tsx"
import ModalClose from "@/components/ModalClose/ModalClose.tsx"
import ModalActions from "@/components/ModalActions/ModalActions.tsx"

export default function AdminPage() {
    const navigate = useNavigate()
    const {auth, isAdmin} = useUser()
    const {OpenDialog} = useDialog()
    const [instances, setInstances] = useState<GetCouchInstanceResponse[]>([])
    const [loading, setLoading] = useState(true)

    const loadInstances = useCallback(async () => {
        setLoading(true)
        try {
            const res = await CouchInstancesAPI.ListCouchInstances({}, auth.getInitReq())
            setInstances(res.instances ?? [])
        } finally {
            setLoading(false)
        }
    }, [auth])

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
            return
        }
        if (!isAdmin) {
            navigate(Path.HomePage)
            return
        }
        void loadInstances()
    }, [auth, isAdmin, navigate, loadInstances])

    async function handleDelete(id: string) {
        await CouchInstancesAPI.DeleteCouchInstance({id}, auth.getInitReq())
        await loadInstances()
    }

    function openAddDialog() {
        OpenDialog(
            <InstanceFormDialog
                onSave={async (url, username, password) => {
                    await CouchInstancesAPI.RegisterCouchInstance({url, username, password}, auth.getInitReq())
                    await loadInstances()
                }}
            />
        )
    }

    function openEditDialog(instance: GetCouchInstanceResponse) {
        OpenDialog(
            <InstanceFormDialog
                initial={instance}
                onSave={async (url, username, password) => {
                    await CouchInstancesAPI.UpdateCouchInstance(
                        {id: instance.id, url, username, password},
                        auth.getInitReq()
                    )
                    await loadInstances()
                }}
            />
        )
    }

    return (
        <div className={cls.AdminPageContainer}>
            <HeroSegment count={instances.length} onAddClick={openAddDialog}/>
            <InstanceList
                instances={instances}
                loading={loading}
                onEdit={openEditDialog}
                onDelete={handleDelete}
            />
        </div>
    )
}

function HeroSegment({count, onAddClick}: { count: number; onAddClick: () => void }) {
    return (
        <div className={cls.HeroContainer}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>Admin</div>
                <h1 className={cls.HeroTitle}>CouchDB instances</h1>
                <p className={cls.HeroSub}>
                    <b>{count} {count === 1 ? "instance" : "instances"}</b>
                    {" · "}<span>manage CouchDB cluster nodes</span>
                </p>
            </div>
            <button className={cls.AddBtn} onClick={onAddClick} type="button">
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2"
                     strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Add instance
            </button>
        </div>
    )
}

function InstanceList({instances, loading, onEdit, onDelete}: {
    instances: GetCouchInstanceResponse[]
    loading: boolean
    onEdit: (instance: GetCouchInstanceResponse) => void
    onDelete: (id: string) => Promise<void>
}) {
    if (loading) {
        return (
            <div className={cls.ContentContainer}>
                <p className={cls.Empty}>Loading…</p>
            </div>
        )
    }

    return (
        <div className={cls.ContentContainer}>
            {instances.length === 0 && (
                <p className={cls.Empty}>No CouchDB instances yet. Add one to get started.</p>
            )}
            {instances.map(instance => (
                <InstanceRow
                    key={instance.id}
                    instance={instance}
                    onEdit={onEdit}
                    onDelete={onDelete}
                />
            ))}
        </div>
    )
}

function InstanceRow({instance, onEdit, onDelete}: {
    instance: GetCouchInstanceResponse
    onEdit: (instance: GetCouchInstanceResponse) => void
    onDelete: (id: string) => Promise<void>
}) {
    const [deleting, setDeleting] = useState(false)

    async function handleDelete() {
        if (!instance.id) return
        setDeleting(true)
        try {
            await onDelete(instance.id)
        } finally {
            setDeleting(false)
        }
    }

    return (
        <div className={cls.RowContainer}>
            <div className={cls.RowInfo}>
                <span className={cls.RowUrl}>{instance.url}</span>
                <span className={cls.RowMeta}>
                    {instance.username} · added {instance.createdAt?.slice(0, 10)}
                </span>
            </div>
            <div className={cls.RowActions}>
                <button className={cls.BtnSecondary} type="button" onClick={() => onEdit(instance)}>
                    Edit
                </button>
                <button className={cls.BtnDanger} type="button" onClick={handleDelete} disabled={deleting}>
                    {deleting ? "Removing…" : "Remove"}
                </button>
            </div>
        </div>
    )
}

function InstanceFormDialog({initial, onSave}: {
    initial?: GetCouchInstanceResponse
    onSave: (url: string, username: string, password: string) => Promise<void>
}) {
    const {CloseDialog} = useDialog()
    const [saving, setSaving] = useState(false)
    const [url, setUrl] = useState(initial?.url ?? "")
    const [username, setUsername] = useState(initial?.username ?? "")
    const [password, setPassword] = useState("")

    const isEdit = !!initial
    const title = isEdit ? "Edit instance" : "Add CouchDB instance"

    async function handleSave() {
        setSaving(true)
        try {
            await onSave(url, username, password)
            CloseDialog()
        } catch (err) {
            console.error(err)
        } finally {
            setSaving(false)
        }
    }

    return (
        <div className={cls.OverlayWrapper}>
            <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
                <div className={cls.ModalHead}>
                    <h2 className={cls.ModalTitle}>{title}</h2>
                    <ModalClose onClick={CloseDialog} disabled={saving} className={cls.ModalClose}/>
                </div>
                <FormField
                    label="URL"
                    placeholder="https://couch.example.com:5984"
                    defaultValue={url}
                    onChange={setUrl}
                    disabled={saving}
                    fieldClassName={cls.Field}
                    labelClassName={cls.FieldLabel}
                    inputClassName={cls.Input}
                />
                <FormField
                    label="Username"
                    placeholder="admin"
                    defaultValue={username}
                    onChange={setUsername}
                    disabled={saving}
                    fieldClassName={cls.Field}
                    labelClassName={cls.FieldLabel}
                    inputClassName={cls.Input}
                />
                <div className={cls.Field}>
                    <span className={cls.FieldLabel}>Password{isEdit ? " (leave blank to keep current)" : ""}</span>
                    <input
                        type="password"
                        placeholder={isEdit ? "••••••••" : "password"}
                        className={cls.Input}
                        onChange={e => setPassword(e.target.value)}
                        disabled={saving}
                        autoComplete="new-password"
                    />
                </div>
                <ModalActions
                    containerClassName={cls.ModalActions}
                    buttons={[
                        {
                            label: saving ? "Saving…" : (isEdit ? "Save changes" : "Add instance"),
                            onClick: handleSave,
                            className: cls.BtnPrimary,
                            disabled: saving,
                        }
                    ]}
                />
            </div>
        </div>
    )
}
