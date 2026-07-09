import {useEffect, useState} from "react"

import cls from "@/dialogs/ManageVaultDialog/widgets/MembersSection/MembersSection.module.css"

import {useVaultMutations, VaultMemberInfo} from "@/app/hooks/Vaults.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

import MemberRow from "@/dialogs/ManageVaultDialog/components/MemberRow/MemberRow.tsx"

interface Props {
    vaultId: string
    currentUserId: string
}

export default function MembersSection({vaultId, currentUserId}: Props) {
    const {listMembers, removeMember} = useVaultMutations()
    const bakeError = useBakeError()
    const [members, setMembers] = useState<VaultMemberInfo[]>([])

    useEffect(() => {
        listMembers(vaultId)
            .then(setMembers)
            .catch(e => bakeError('Error during vault members loading', e))
    }, [vaultId])

    function handleRemoveMember(userId: string) {
        return removeMember(vaultId, userId)
            .then(() => setMembers(prev => prev.filter(m => m.userId !== userId)))
            .catch(e => bakeError('Error removing member', e))
    }

    return (
        <section className={cls.MembersSectionContainer}>
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
    )
}
