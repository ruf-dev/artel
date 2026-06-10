import {useEffect, useState} from "react"

import cls from "@/components/ManageVaultDialog/ManageVaultDialog.module.css"

import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {useVaultMutations, VaultMemberInfo, VaultInviteItem} from "@/app/hooks/Vaults.ts"
import {useDialog} from "@/app/hooks/Dialog.ts"
import ModalClose from "@/components/ModalClose/ModalClose.tsx"

interface Props {
    vault: VaultItem
    currentUserId: string
    onClose: () => void
    onDeleted: () => void
}

export default function ManageVaultDialog({vault, currentUserId, onClose, onDeleted}: Props) {
    const {listMembers, removeMember, listInvites, revokeInvite, remove} = useVaultMutations()
    const {OpenDialog} = useDialog()

    const [members, setMembers] = useState<VaultMemberInfo[]>([])
    const [invites, setInvites] = useState<VaultInviteItem[]>([])
    const [deleting, setDeleting] = useState(false)

    useEffect(() => {
        if (!vault.id) return
        listMembers(vault.id).then(setMembers).catch(console.error)
        listInvites(vault.id).then(setInvites).catch(console.error)
    }, [vault.id])

    async function handleRemoveMember(userId: string) {
        if (!vault.id) return
        await removeMember(vault.id, userId)
        setMembers(prev => prev.filter(m => m.userId !== userId))
    }

    async function handleRevokeInvite(inviteId: string) {
        if (!vault.id) return
        await revokeInvite(vault.id, inviteId)
        setInvites(prev => prev.map(inv => inv.id === inviteId ? {...inv, revoked: true} : inv))
    }

    function openCreateInviteDialog() {
        if (!vault.id) return
        OpenDialog(
            <CreateInviteLinkDialog
                vaultId={vault.id}
                onCreated={(inv) => setInvites(prev => [inv, ...prev])}
            />
        )
    }

    async function handleDelete() {
        if (!vault.id) return
        setDeleting(true)
        try {
            await remove(vault.id)
            onDeleted()
        } finally {
            setDeleting(false)
        }
    }

    return (
        <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>Manage vault</h2>
                <ModalClose onClick={onClose} className={cls.ModalClose}/>
            </div>
            <p className={cls.VaultName}>{vault.name}</p>

            <section className={cls.Section}>
                <div className={cls.SectionHead}>
                    <span className={cls.SectionTitle}>Members</span>
                </div>
                <div className={cls.MemberList}>
                    {members.map(m => (
                        <MemberRow
                            key={m.userId}
                            member={m}
                            isSelf={m.userId === currentUserId}
                            onRemove={() => handleRemoveMember(m.userId ?? "")}
                        />
                    ))}
                    {members.length === 0 && <p className={cls.Empty}>No members yet.</p>}
                </div>
            </section>

            <section className={cls.Section}>
                <div className={cls.SectionHead}>
                    <span className={cls.SectionTitle}>Invite links</span>
                    <button className={cls.BtnSecondary} onClick={openCreateInviteDialog}>
                        + New link
                    </button>
                </div>
                <div className={cls.InviteList}>
                    {invites.map(inv => (
                        <InviteRow
                            key={inv.id}
                            invite={inv}
                            onRevoke={() => handleRevokeInvite(inv.id ?? "")}
                        />
                    ))}
                    {invites.length === 0 && <p className={cls.Empty}>No invite links yet.</p>}
                </div>
            </section>

            <section className={cls.DangerZone}>
                <div className={cls.DangerZoneRow}>
                    <div>
                        <div className={cls.DangerZoneTitle}>Delete this vault</div>
                        <div className={cls.DangerZoneSub}>Permanent. Connection string stops working immediately.</div>
                    </div>
                    <button className={cls.BtnDanger} onClick={handleDelete} disabled={deleting}>
                        {deleting ? "Deleting…" : "Delete"}
                    </button>
                </div>
            </section>

            <div className={cls.ModalFooter}>
                <button className={cls.BtnGhost} onClick={onClose}>Close</button>
            </div>
        </div>
    )
}

function MemberRow({member, isSelf, onRemove}: {
    member: VaultMemberInfo
    isSelf: boolean
    onRemove: () => void
}) {
    const [removing, setRemoving] = useState(false)
    const initial = (member.username || member.email || "?")[0].toUpperCase()

    async function handleRemove() {
        setRemoving(true)
        try { await onRemove() } finally { setRemoving(false) }
    }

    return (
        <div className={cls.MemberRow}>
            <div className={cls.Avatar}>{initial}</div>
            <div className={cls.MemberInfo}>
                <span className={cls.MemberName}>{member.username || member.email}</span>
                {member.username && member.email && (
                    <span className={cls.MemberEmail}>{member.email}</span>
                )}
            </div>
            <RoleBadge role={member.role ?? ""} />
            {!isSelf && member.role !== "owner" && (
                <button className={cls.BtnRemove} onClick={handleRemove} disabled={removing}>
                    {removing ? "…" : "Remove"}
                </button>
            )}
        </div>
    )
}

function InviteRow({invite, onRevoke}: {
    invite: VaultInviteItem
    onRevoke: () => void
}) {
    const [copied, setCopied] = useState(false)
    const [revoking, setRevoking] = useState(false)
    const link = `${window.location.origin}/join/${invite.token}`

    function handleCopy() {
        navigator.clipboard.writeText(link).then(() => {
            setCopied(true)
            setTimeout(() => setCopied(false), 2000)
        })
    }

    async function handleRevoke() {
        setRevoking(true)
        try { await onRevoke() } finally { setRevoking(false) }
    }

    return (
        <div className={`${cls.InviteRow} ${invite.revoked ? cls.Revoked : ""}`}>
            <RoleBadge role={invite.role ?? ""} />
            <span className={cls.InviteToken}>{invite.token?.slice(0, 8)}…</span>
            {invite.revoked ? (
                <span className={cls.RevokedLabel}>Revoked</span>
            ) : (
                <>
                    <button className={cls.BtnCopy} onClick={handleCopy}>
                        {copied ? "Copied!" : "Copy link"}
                    </button>
                    <button className={cls.BtnRemove} onClick={handleRevoke} disabled={revoking}>
                        {revoking ? "…" : "Revoke"}
                    </button>
                </>
            )}
        </div>
    )
}

function RoleBadge({role}: { role: string }) {
    return <span className={`${cls.RoleBadge} ${cls[`Role_${role}`] ?? ""}`}>{role}</span>
}

function CreateInviteLinkDialog({vaultId, onCreated}: {
    vaultId: string
    onCreated: (inv: VaultInviteItem) => void
}) {
    const {createInvite} = useVaultMutations()
    const {CloseDialog} = useDialog()
    const [role, setRole] = useState<"reader" | "maintainer">("reader")
    const [creating, setCreating] = useState(false)
    const [created, setCreated] = useState<VaultInviteItem | null>(null)
    const [copied, setCopied] = useState(false)

    async function handleCreate() {
        setCreating(true)
        try {
            const inv = await createInvite(vaultId, role)
            setCreated(inv)
            onCreated(inv)
        } catch (err) {
            console.error(err)
        } finally {
            setCreating(false)
        }
    }

    function handleCopy() {
        if (!created?.token) return
        navigator.clipboard.writeText(`${window.location.origin}/join/${created.token}`).then(() => {
            setCopied(true)
            setTimeout(() => setCopied(false), 2000)
        })
    }

    return (
        <div className={cls.SubModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle}>New invite link</h2>
                <ModalClose onClick={CloseDialog} className={cls.ModalClose}/>
            </div>

            {!created ? (
                <>
                    <p className={cls.ModalSub}>Choose what the invited user can do.</p>
                    <div className={cls.RoleSelector}>
                        <RoleOption
                            selected={role === "reader"}
                            label="Reader"
                            desc="Can view and sync the vault, cannot write."
                            onSelect={() => setRole("reader")}
                        />
                        <RoleOption
                            selected={role === "maintainer"}
                            label="Maintainer"
                            desc="Can read and write to the vault."
                            onSelect={() => setRole("maintainer")}
                        />
                    </div>
                    <div className={cls.ModalFooter}>
                        <button className={cls.BtnGhost} onClick={CloseDialog}>Cancel</button>
                        <button className={cls.BtnPrimary} onClick={handleCreate} disabled={creating}>
                            {creating ? "Creating…" : "Create link"}
                        </button>
                    </div>
                </>
            ) : (
                <>
                    <p className={cls.ModalSub}>Share this link. Anyone with it can join as <strong>{role}</strong>.</p>
                    <div className={cls.LinkBox}>
                        <span className={cls.LinkText}>{window.location.origin}/join/{created.token}</span>
                    </div>
                    <div className={cls.ModalFooter}>
                        <button className={cls.BtnGhost} onClick={CloseDialog}>Done</button>
                        <button className={cls.BtnPrimary} onClick={handleCopy}>
                            {copied ? "Copied!" : "Copy link"}
                        </button>
                    </div>
                </>
            )}
        </div>
    )
}

function RoleOption({selected, label, desc, onSelect}: {
    selected: boolean
    label: string
    desc: string
    onSelect: () => void
}) {
    return (
        <button
            className={`${cls.RoleOption} ${selected ? cls.RoleOptionSelected : ""}`}
            onClick={onSelect}
            type="button"
        >
            <span className={cls.RoleOptionLabel}>{label}</span>
            <span className={cls.RoleOptionDesc}>{desc}</span>
        </button>
    )
}
