import {useEffect} from "react"
import {Button} from "@vervstack/chures"
import {generatePath, useNavigate, useParams} from "react-router-dom"

import cls from "@/pages/skills/SkillsPage.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useSkills} from "@/app/hooks/Skills.ts"
import {useVaults} from "@/app/hooks/Vaults.ts"
import useUser from "@/hooks/user/User.ts"
import HeroSegment from "@/components/HeroSegment/HeroSegment.tsx"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import ContentSegment from "@/pages/skills/components/ContentSegment/ContentSegment.tsx"
import CreateSkillDialog from "@/dialogs/CreateSkillDialog/CreateSkillDialog.tsx"

export default function SkillsPage() {
    const {vaultId} = useParams()
    const navigate = useNavigate()
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const {vaults} = useVaults()
    const {skills, loading, fetch: fetchSkills} = useSkills()

    // No vault named in the URL yet — if the user has exactly one, jump straight to it
    // (same fallback used by NotesPage), otherwise the vault-picker below handles it.
    useEffect(() => {
        if (vaultId || vaults.length !== 1) return
        navigate(generatePath(Path.SkillsPageVault, {vaultId: vaults[0].id ?? ""}), {replace: true})
    }, [vaultId, vaults, navigate])

    useEffect(() => {
        if (vaultId && auth.isAuthenticated()) {
            void fetchSkills(vaultId)
        }
    }, [vaultId, auth, fetchSkills])

    const hotPlugCount = skills.filter(s => s.isHotPlug).length

    return (
        <div className={cls.Root}>
            <HeroSegment
                eyebrow="Skills"
                title="Agent skills"
                subtitle={
                    vaultId ? (
                        <>
                            <b>{loading ? "…" : `${skills.length} total · ${hotPlugCount} hot-plug`}</b>
                            {" · "}<span>instructions your agent can pull in on demand</span>
                        </>
                    ) : (
                        <span>pick a vault to manage its skills</span>
                    )
                }
                action={vaultId ? (
                    <Button variant="primary" onClick={() => OpenDialog(<CreateSkillDialog vaultId={vaultId}/>)}>
                        <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor"
                             strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <line x1="12" y1="5" x2="12" y2="19"/>
                            <line x1="5" y1="12" x2="19" y2="12"/>
                        </svg>
                        New skill
                    </Button>
                ) : undefined}
            />
            {vaultId ? (
                <ContentSegment vaultId={vaultId}/>
            ) : (
                <div className={cls.VaultPicker}>
                    {vaults.map(v => (
                        <SelectOption
                            key={v.id}
                            label={v.name ?? ""}
                            selected={false}
                            onSelect={() => navigate(generatePath(Path.SkillsPageVault, {vaultId: v.id ?? ""}))}
                        />
                    ))}
                    {vaults.length === 0 && <p className={cls.Empty}>No vaults available. Create one first.</p>}
                </div>
            )}
        </div>
    )
}
