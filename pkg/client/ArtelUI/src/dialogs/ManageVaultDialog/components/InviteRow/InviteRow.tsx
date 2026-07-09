import {useState} from "react"
import {Button} from "@vervstack/chures"
import cls from "@/dialogs/ManageVaultDialog/components/InviteRow/InviteRow.module.css"
import {cn} from "@/app/utils/cn.ts"

import {VaultInviteItem} from "@/app/hooks/Vaults.ts"
import RoleBadge from "@/dialogs/ManageVaultDialog/components/RoleBadge/RoleBadge.tsx"

interface Props {
    invite: VaultInviteItem
    onRevoke: () => Promise<void>
}

export default function InviteRow({invite, onRevoke}: Props) {
    const [copied, setCopied] = useState(false)
    const [revoking, setRevoking] = useState(false)
    const link = `${window.location.origin}/join/${invite.token}`

    function handleCopy() {
        navigator.clipboard.writeText(link).then(() => {
            setCopied(true)
            setTimeout(() => setCopied(false), 2000)
        })
    }

    function handleRevoke() {
        setRevoking(true)
        onRevoke().finally(() => setRevoking(false))
    }

    return (
        <div className={cn(cls.InviteRowContainer, invite.revoked && cls.Revoked)}>
            <RoleBadge role={invite.role ?? ""}/>
            <span className={cls.InviteToken}>{invite.token?.slice(0, 8)}…</span>
            {invite.revoked ? (
                <span className={cls.RevokedLabel}>Revoked</span>
            ) : (
                <>
                    <Button variant="secondary" onClick={handleCopy}>
                        {copied ? "Copied!" : "Copy link"}
                    </Button>
                    <Button variant="danger" onClick={handleRevoke} disabled={revoking}>
                        {revoking ? "…" : "Revoke"}
                    </Button>
                </>
            )}
        </div>
    )
}
