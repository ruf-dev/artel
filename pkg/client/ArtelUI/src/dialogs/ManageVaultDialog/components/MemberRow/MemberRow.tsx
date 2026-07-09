import {useState} from "react"
import {Button} from "@vervstack/chures"
import cls from "@/dialogs/ManageVaultDialog/components/MemberRow/MemberRow.module.css"

import {VaultMemberInfo} from "@/app/hooks/Vaults.ts"
import RoleBadge from "@/dialogs/ManageVaultDialog/components/RoleBadge/RoleBadge.tsx"

interface Props {
    member: VaultMemberInfo
    isSelf: boolean
    onRemove: () => Promise<void>
}

export default function MemberRow({member, isSelf, onRemove}: Props) {
    const [removing, setRemoving] = useState(false)
    const initial = (member.username || member.email || "?")[0].toUpperCase()

    function handleRemove() {
        setRemoving(true)
        onRemove().finally(() => setRemoving(false))
    }

    return (
        <div className={cls.MemberRowContainer}>
            <div className={cls.Avatar}>{initial}</div>
            <div className={cls.MemberInfo}>
                <span className={cls.MemberName}>{member.username || member.email}</span>
                {member.username && member.email && (
                    <span className={cls.MemberEmail}>{member.email}</span>
                )}
            </div>
            <RoleBadge role={member.role ?? ""}/>
            {!isSelf && member.role !== "owner" && (
                <Button variant="danger" onClick={handleRemove} disabled={removing}>
                    {removing ? "…" : "Remove"}
                </Button>
            )}
        </div>
    )
}
