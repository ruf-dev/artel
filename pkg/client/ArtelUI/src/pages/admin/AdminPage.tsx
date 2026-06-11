import {useState, useEffect, useCallback} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/admin/AdminPage.module.css"

import {AdminCouchAPI, CouchUserEntry} from "@/app/api/artel/admin_couch.pb.ts"
import {AdminUsersAPI, ArtelUserEntry, ArtelUserDetails, UserSession} from "@/app/api/artel/admin_users.pb.ts"
import {CouchInstancesAPI, GetCouchInstanceResponse, GetCouchInstanceStatusResponse} from "@/app/api/artel/couch_instances.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import useUser from "@/hooks/user/User.ts"

import FormField from "@/components/FormField/FormField.tsx"
import ModalClose from "@/components/ModalClose/ModalClose.tsx"
import ModalActions from "@/components/ModalActions/ModalActions.tsx"
import ConfirmDialog from "@/components/ConfirmDialog/ConfirmDialog.tsx"
import {useBakeError} from "@/app/hooks/useErrorToast"

type Tab = "instances" | "couch_users" | "users"

export default function AdminPage() {
    const navigate = useNavigate()
    const {auth, isAdmin} = useUser()
    const [tab, setTab] = useState<Tab>("instances")

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
            return
        }
        if (!isAdmin) {
            navigate(Path.HomePage)
        }
    }, [auth, isAdmin, navigate])

    return (
        <div className={cls.AdminPageContainer}>
            <AdminHero tab={tab} />
            <TabBar tab={tab} onTabChange={setTab} />
            {tab === "instances" && <InstancesTab />}
            {tab === "couch_users" && <UsersTab />}
            {tab === "users" && <ArtelUsersTab />}
        </div>
    )
}

// ── Shared header ────────────────────────────────────────────

function AdminHero({tab}: {tab: Tab}) {
    const titleMap: Record<Tab, string> = {
        instances: "CouchDB instances",
        couch_users: "CouchDB users",
        users: "Artel users",
    }
    const title = titleMap[tab]
    return (
        <div className={cls.HeroContainer}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>Admin</div>
                <h1 className={cls.HeroTitle}>{title}</h1>
            </div>
        </div>
    )
}

function TabBar({tab, onTabChange}: {tab: Tab; onTabChange: (t: Tab) => void}) {
    const instancesCls = tab === "instances" ? `${cls.Tab} ${cls.TabActive}` : cls.Tab
    const couchUsersCls = tab === "couch_users" ? `${cls.Tab} ${cls.TabActive}` : cls.Tab
    const usersCls = tab === "users" ? `${cls.Tab} ${cls.TabActive}` : cls.Tab
    return (
        <div className={cls.TabBar}>
            <button type="button" className={instancesCls} onClick={() => onTabChange("instances")}>
                Instances
            </button>
            <button type="button" className={couchUsersCls} onClick={() => onTabChange("couch_users")}>
                CouchUsers
            </button>
            <button type="button" className={usersCls} onClick={() => onTabChange("users")}>
                Users
            </button>
        </div>
    )
}

// ── Instances tab ─────────────────────────────────────────────

function InstancesTab() {
    const {auth} = useUser()
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

    useEffect(() => { void loadInstances() }, [loadInstances])

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
        <div className={cls.ContentContainer}>
            <InstancesActionBar count={instances.length} onAddClick={openAddDialog} />
            <InstanceList instances={instances} loading={loading} onEdit={openEditDialog} onDelete={handleDelete} />
        </div>
    )
}

function InstancesActionBar({count, onAddClick}: {count: number; onAddClick: () => void}) {
    return (
        <div className={cls.TabActionBar}>
            <p className={cls.HeroSub}>
                <b>{count} {count === 1 ? "instance" : "instances"}</b>
                {" · "}<span>manage CouchDB cluster nodes</span>
            </p>
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
        return <p className={cls.Empty}>Loading…</p>
    }
    if (instances.length === 0) {
        return <p className={cls.Empty}>No CouchDB instances yet. Add one to get started.</p>
    }
    return (
        <>
            {instances.map(instance => (
                <InstanceRow key={instance.id} instance={instance} onEdit={onEdit} onDelete={onDelete} />
            ))}
        </>
    )
}

function InstanceRow({instance, onEdit, onDelete}: {
    instance: GetCouchInstanceResponse
    onEdit: (instance: GetCouchInstanceResponse) => void
    onDelete: (id: string) => Promise<void>
}) {
    const [deleting, setDeleting] = useState(false)
    const {OpenDialog} = useDialog()
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
        <div className={cls.RowContainer}>
            <div className={cls.RowInfo}>
                <span className={cls.RowUrl}>{instance.url}</span>
                <span className={cls.RowMeta}>{instance.username} · added {instance.createdAt?.slice(0, 10)}</span>
            </div>
            <div className={cls.RowActions}>
                <button className={cls.BtnSecondary} type="button" onClick={() => onEdit(instance)}>Edit</button>
                <button className={cls.BtnDanger} type="button" onClick={handleDelete} disabled={deleting}>
                    {deleting ? "Removing…" : "Remove"}
                </button>
            </div>
        </div>
    )
}

function SetupStatusSection({instanceId, refreshKey}: {instanceId: string; refreshKey: number}) {
    const {auth} = useUser()
    const bakeError = useBakeError()
    const [status, setStatus] = useState<GetCouchInstanceStatusResponse | null>(null)

    useEffect(() => {
        async function loadStatus() {
            try {
                const res = await CouchInstancesAPI.GetCouchInstanceStatus({id: instanceId}, auth.getInitReq())
                setStatus(res)
            } catch (err) {
                bakeError("Failed to load instance status", err)
            }
        }
        void loadStatus()
    }, [instanceId, auth, refreshKey, bakeError])

    const checks = [
        {label: "Single-node cluster", ok: status?.clusterModeEnabled},
        {label: "_users database", ok: status?.usersDbInitialized},
        {label: "_replicator database", ok: status?.replicatorDbInitialized},
    ]

    return (
        <div className={cls.StatusSection}>
            <div className={cls.StatusSectionTitle}>Setup status</div>
            {checks.map(check => (
                <div key={check.label} className={cls.StatusRow}>
                    <span className={cls.StatusLabel}>{check.label}</span>
                    <span className={status === null ? cls.StatusLoading : check.ok ? cls.StatusOk : cls.StatusFail}>
                        {status === null ? "…" : check.ok ? "✓" : "✗"}
                    </span>
                </div>
            ))}
        </div>
    )
}

function InstanceFormDialog({initial, onSave}: {
    initial?: GetCouchInstanceResponse
    onSave: (url: string, username: string, password: string) => Promise<void>
}) {
    const {CloseDialog} = useDialog()
    const {auth} = useUser()
    const bakeError = useBakeError()
    const [saving, setSaving] = useState(false)
    const [settingUp, setSettingUp] = useState(false)
    const [setupRefresh, setSetupRefresh] = useState(0)
    const [url, setUrl] = useState(initial?.url ?? "")
    const [username, setUsername] = useState(initial?.username ?? "")
    const [password, setPassword] = useState("")

    const isEdit = !!initial
    const busy = saving || settingUp
    const title = isEdit ? "Edit instance" : "Add CouchDB instance"

    async function handleSave() {
        setSaving(true)
        try {
            await onSave(url, username, password)
            CloseDialog()
        } catch (err) {
            bakeError(isEdit ? "Failed to update instance" : "Failed to add instance", err)
        } finally {
            setSaving(false)
        }
    }

    async function handleSetup() {
        if (!initial?.id) return
        setSettingUp(true)
        try {
            await CouchInstancesAPI.SetupCouchInstance({id: initial.id}, auth.getInitReq())
            setSetupRefresh(r => r + 1)
        } catch (err) {
            bakeError("Failed to setup instance", err)
        } finally {
            setSettingUp(false)
        }
    }

    const buttons = isEdit
        ? [
            {label: settingUp ? "Setting up…" : "Setup", onClick: handleSetup, className: cls.BtnSecondary, disabled: busy},
            {label: saving ? "Saving…" : "Save changes", onClick: handleSave, className: cls.BtnPrimary, disabled: busy},
          ]
        : [
            {label: saving ? "Saving…" : "Add instance", onClick: handleSave, className: cls.BtnPrimary, disabled: busy},
          ]

    return (
        <div className={cls.ModalContainer} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>{title}</h2>
                <ModalClose onClick={CloseDialog} disabled={busy} className={cls.ModalClose}/>
            </div>
            {isEdit && initial?.id && <SetupStatusSection instanceId={initial.id} refreshKey={setupRefresh} />}
            <FormField
                label="URL"
                placeholder="https://couch.example.com:5984"
                defaultValue={url}
                onChange={setUrl}
                disabled={busy}
                fieldClassName={cls.Field}
                labelClassName={cls.FieldLabel}
                inputClassName={cls.Input}
            />
            <FormField
                label="Username"
                placeholder="admin"
                defaultValue={username}
                onChange={setUsername}
                disabled={busy}
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
                    disabled={busy}
                    autoComplete="new-password"
                />
            </div>
            <ModalActions containerClassName={cls.ModalActions} buttons={buttons} />
        </div>
    )
}

// ── Users tab ─────────────────────────────────────────────────

function UsersTab() {
    const {auth} = useUser()
    const [instances, setInstances] = useState<GetCouchInstanceResponse[]>([])
    const [selectedInstanceId, setSelectedInstanceId] = useState("")
    const [users, setUsers] = useState<CouchUserEntry[]>([])
    const [usersLoading, setUsersLoading] = useState(false)

    useEffect(() => {
        async function loadInstances() {
            const res = await CouchInstancesAPI.ListCouchInstances({}, auth.getInitReq())
            setInstances(res.instances ?? [])
        }
        void loadInstances()
    }, [auth])

    const loadUsers = useCallback(async () => {
        if (!selectedInstanceId) {
            setUsers([])
            return
        }
        setUsersLoading(true)
        try {
            const res = await AdminCouchAPI.ListCouchUsers({instanceId: selectedInstanceId}, auth.getInitReq())
            setUsers(res.users ?? [])
        } finally {
            setUsersLoading(false)
        }
    }, [auth, selectedInstanceId])

    useEffect(() => { void loadUsers() }, [loadUsers])

    return (
        <div className={cls.ContentContainer}>
            <InstanceSelector instances={instances} value={selectedInstanceId} onChange={setSelectedInstanceId} />
            <UserList
                instanceId={selectedInstanceId}
                users={users}
                loading={usersLoading}
                onRefresh={loadUsers}
            />
        </div>
    )
}

function InstanceSelector({instances, value, onChange}: {
    instances: GetCouchInstanceResponse[]
    value: string
    onChange: (id: string) => void
}) {
    return (
        <div className={cls.Field}>
            <span className={cls.FieldLabel}>CouchDB instance</span>
            <select
                className={cls.InstanceSelector}
                value={value}
                onChange={e => onChange(e.target.value)}
            >
                <option value="">— select an instance —</option>
                {instances.map(inst => (
                    <option key={inst.id} value={inst.id ?? ""}>{inst.url}</option>
                ))}
            </select>
        </div>
    )
}

function UserList({instanceId, users, loading, onRefresh}: {
    instanceId: string
    users: CouchUserEntry[]
    loading: boolean
    onRefresh: () => Promise<void>
}) {
    if (!instanceId) {
        return <p className={cls.Empty}>Select an instance above to view its users.</p>
    }
    if (loading) {
        return <p className={cls.Empty}>Loading…</p>
    }
    if (users.length === 0) {
        return <p className={cls.Empty}>No users found in this instance.</p>
    }
    return (
        <>
            {users.map(user => (
                <UserRow key={user.name} instanceId={instanceId} user={user} onRefresh={onRefresh} />
            ))}
        </>
    )
}

function UserRow({instanceId, user, onRefresh}: {
    instanceId: string
    user: CouchUserEntry
    onRefresh: () => Promise<void>
}) {
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const bakeError = useBakeError()
    const [deleting, setDeleting] = useState(false)

    function handleDelete() {
        const name = user.name
        if (!name) return
        OpenDialog(
            <ConfirmDialog
                title="Delete user"
                message={`Delete "${name}"? This cannot be undone.`}
                danger
                confirmLabel="Delete"
                onConfirm={async () => {
                    setDeleting(true)
                    try {
                        await AdminCouchAPI.DeleteCouchUser({instanceId, username: name}, auth.getInitReq())
                        await onRefresh()
                    } catch (err) {
                        bakeError("Failed to delete user", err)
                    } finally {
                        setDeleting(false)
                    }
                }}
            />
        )
    }

    function openChangePassword() {
        if (!user.name) return
        OpenDialog(
            <ChangePasswordDialog instanceId={instanceId} username={user.name} />
        )
    }

    function openManageAccess() {
        if (!user.name) return
        OpenDialog(
            <ManageAccessDialog instanceId={instanceId} username={user.name} />
        )
    }

    return (
        <div className={cls.RowContainer}>
            <div className={cls.RowInfo}>
                <span className={cls.RowUrl}>{user.name}</span>
                <span className={cls.RowMeta}>{(user.roles ?? []).join(", ") || "no roles"}</span>
            </div>
            <div className={cls.RowActions}>
                <button className={cls.BtnSecondary} type="button" onClick={openManageAccess}>
                    DB Access
                </button>
                <button className={cls.BtnSecondary} type="button" onClick={openChangePassword}>
                    Change Password
                </button>
                <button className={cls.BtnDanger} type="button" onClick={handleDelete} disabled={deleting}>
                    {deleting ? "Deleting…" : "Delete"}
                </button>
            </div>
        </div>
    )
}

function ChangePasswordDialog({instanceId, username}: {instanceId: string; username: string}) {
    const {auth} = useUser()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const [saving, setSaving] = useState(false)
    const [newPassword, setNewPassword] = useState("")

    async function handleSave() {
        if (!newPassword) return
        setSaving(true)
        try {
            await AdminCouchAPI.ChangeCouchUserPassword({instanceId, username, newPassword}, auth.getInitReq())
            CloseDialog()
        } catch (err) {
            bakeError("Failed to change password", err)
        } finally {
            setSaving(false)
        }
    }

    return (
        <div className={cls.ModalContainer} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>Change password</h2>
                <ModalClose onClick={CloseDialog} disabled={saving} className={cls.ModalClose}/>
            </div>
            <p className={cls.ModalSubtitle}>User: <b>{username}</b></p>
            <div className={cls.Field}>
                <span className={cls.FieldLabel}>New password</span>
                <input
                    type="password"
                    placeholder="new password"
                    className={cls.Input}
                    onChange={e => setNewPassword(e.target.value)}
                    disabled={saving}
                    autoComplete="new-password"
                />
            </div>
            <ModalActions
                containerClassName={cls.ModalActions}
                buttons={[
                    {
                        label: saving ? "Saving…" : "Save password",
                        onClick: handleSave,
                        className: cls.BtnPrimary,
                        disabled: saving || !newPassword,
                    }
                ]}
            />
        </div>
    )
}

function ManageAccessDialog({instanceId, username}: {instanceId: string; username: string}) {
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
                grant ? next.add(db) : next.delete(db)
                return next
            })
        } catch (err) {
            bakeError(grant ? "Failed to grant access" : "Failed to revoke access", err)
        } finally {
            setBusyDb(null)
        }
    }

    return (
        <div className={cls.ModalContainer} role="dialog" aria-modal="true">
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
            <ModalActions
                containerClassName={cls.ModalActions}
                buttons={[{label: "Done", onClick: CloseDialog, className: cls.BtnPrimary}]}
            />
        </div>
    )
}

function DbAccessList({databases, grantedSet, busyDb, onToggle}: {
    databases: string[]
    grantedSet: Set<string>
    busyDb: string | null
    onToggle: (db: string, grant: boolean) => Promise<void>
}) {
    return (
        <div className={cls.DbAccessList}>
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

function DbAccessRow({db, granted, busy, onToggle}: {
    db: string
    granted: boolean
    busy: boolean
    onToggle: (db: string, grant: boolean) => Promise<void>
}) {
    return (
        <div className={cls.DbAccessRow}>
            <span className={cls.DbAccessName}>{db}</span>
            <button
                type="button"
                className={granted ? cls.BtnDanger : cls.BtnPrimary}
                disabled={busy}
                onClick={() => onToggle(db, !granted)}
            >
                {busy ? "…" : granted ? "Revoke" : "Grant"}
            </button>
        </div>
    )
}

// ── Artel Users tab ───────────────────────────────────────────

function ArtelUsersTab() {
    const {auth} = useUser()
    const [search, setSearch] = useState("")
    const [users, setUsers] = useState<ArtelUserEntry[]>([])
    const [total, setTotal] = useState(0)
    const [loading, setLoading] = useState(true)
    const bakeError = useBakeError()
    const LIMIT = 50

    const loadUsers = useCallback(async () => {
        setLoading(true)
        try {
            const res = await AdminUsersAPI.ListArtelUsers({search, limit: LIMIT, offset: 0}, auth.getInitReq())
            setUsers(res.users ?? [])
            setTotal(res.total ? Number(res.total) : 0)
        } catch (err) {
            bakeError("Failed to load users", err)
        } finally {
            setLoading(false)
        }
    }, [auth, search])

    useEffect(() => { void loadUsers() }, [loadUsers])

    return (
        <div className={cls.ContentContainer}>
            <div className={cls.Field}>
                <input
                    className={cls.Input}
                    placeholder="Search by username or email…"
                    value={search}
                    onChange={e => setSearch(e.target.value)}
                />
            </div>
            <ArtelUserList users={users} loading={loading} total={total} />
        </div>
    )
}

function ArtelUserList({users, loading, total}: {users: ArtelUserEntry[]; loading: boolean; total: number}) {
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

function ArtelUserRow({user}: {user: ArtelUserEntry}) {
    const {OpenDialog} = useDialog()

    function openDetail() {
        if (!user.userId) return
        OpenDialog(<ArtelUserDetailDialog userId={user.userId} />)
    }

    return (
        <div className={cls.RowContainer}>
            <div className={cls.RowInfo}>
                <span className={cls.RowUrl}>{user.username || "—"}</span>
                <span className={cls.RowMeta}>{user.email || "no email"}</span>
            </div>
            <div className={cls.RowActions}>
                <button className={cls.BtnSecondary} type="button" onClick={openDetail}>
                    Details
                </button>
            </div>
        </div>
    )
}

function ArtelUserDetailDialog({userId}: {userId: string}) {
    const {auth} = useUser()
    const {CloseDialog, OpenDialog} = useDialog()
    const bakeError = useBakeError()
    const [details, setDetails] = useState<ArtelUserDetails | null>(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        async function load() {
            try {
                const res = await AdminUsersAPI.GetArtelUser({userId}, auth.getInitReq())
                setDetails(res.user ?? null)
            } catch (err) {
                bakeError("Failed to load user details", err)
            } finally {
                setLoading(false)
            }
        }
        void load()
    }, [auth, userId])

    function openSessions() {
        OpenDialog(<UserSessionsDialog userId={userId} />)
    }

    return (
        <div className={cls.ModalContainer} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>User details</h2>
                <ModalClose onClick={CloseDialog} className={cls.ModalClose} />
            </div>
            {loading ? (
                <p className={cls.Empty}>Loading…</p>
            ) : details ? (
                <>
                    <p className={cls.ModalSubtitle}><b>{details.username || "—"}</b> · {details.email || "no email"}</p>
                    <div className={cls.Field}>
                        <span className={cls.FieldLabel}>Permissions</span>
                        <span className={cls.RowMeta}>
                            {details.isAdministrator ? "Admin " : ""}
                            {details.hasEmails ? "Emails " : ""}
                            {details.hasTaskTrackers ? "TaskTrackers " : ""}
                            {details.hasNotes ? "Notes " : ""}
                            {!details.isAdministrator && !details.hasEmails && !details.hasTaskTrackers && !details.hasNotes ? "none" : ""}
                        </span>
                    </div>
                    <div className={cls.Field}>
                        <span className={cls.FieldLabel}>Subscription</span>
                        <span className={cls.RowMeta}>{details.subscriptionActive ? "Active" : "Inactive"}</span>
                    </div>
                </>
            ) : (
                <p className={cls.Empty}>User not found.</p>
            )}
            <ModalActions
                containerClassName={cls.ModalActions}
                buttons={[
                    {label: "Sessions", onClick: openSessions, className: cls.BtnSecondary},
                    {label: "Close", onClick: CloseDialog, className: cls.BtnPrimary},
                ]}
            />
        </div>
    )
}

function UserSessionsDialog({userId}: {userId: string}) {
    const {auth} = useUser()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const [sessions, setSessions] = useState<UserSession[]>([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        async function load() {
            try {
                const res = await AdminUsersAPI.GetUserSessions({userId}, auth.getInitReq())
                setSessions(res.sessions ?? [])
            } catch (err) {
                bakeError("Failed to load sessions", err)
            } finally {
                setLoading(false)
            }
        }
        void load()
    }, [auth, userId])

    function formatDate(iso?: string) {
        if (!iso) return "—"
        return new Date(iso).toLocaleString()
    }

    return (
        <div className={cls.ModalContainer} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>Sessions</h2>
                <ModalClose onClick={CloseDialog} className={cls.ModalClose} />
            </div>
            {loading ? (
                <p className={cls.Empty}>Loading…</p>
            ) : sessions.length === 0 ? (
                <p className={cls.Empty}>No active sessions.</p>
            ) : (
                sessions.map(s => (
                    <div key={s.sessionId} className={cls.RowContainer}>
                        <div className={cls.RowInfo}>
                            <span className={cls.RowUrl}>Created: {formatDate(s.createdAt as unknown as string)}</span>
                            <span className={cls.RowMeta}>Expires: {formatDate(s.expiresAt as unknown as string)}</span>
                        </div>
                    </div>
                ))
            )}
            <ModalActions
                containerClassName={cls.ModalActions}
                buttons={[{label: "Close", onClick: CloseDialog, className: cls.BtnPrimary}]}
            />
        </div>
    )
}
