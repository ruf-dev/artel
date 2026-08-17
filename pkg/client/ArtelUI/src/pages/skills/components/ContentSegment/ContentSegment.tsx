import {useSkills} from "@/app/hooks/Skills.ts"
import SkillCard from "@/widgets/SkillCard/SkillCard.tsx"
import cls from "@/pages/skills/components/ContentSegment/ContentSegment.module.css"

interface Props {
    vaultId: string
}

export default function ContentSegment({vaultId}: Props) {
    const {skills, loading} = useSkills()

    if (loading) {
        return (
            <div className={cls.ContentSegmentContainer}>
                <p className={cls.Empty}>Loading…</p>
            </div>
        )
    }

    return (
        <div className={cls.ContentSegmentContainer}>
            <div className={cls.List}>
                {skills.map(skill => (
                    <SkillCard key={skill.slug} vaultId={vaultId} skill={skill}/>
                ))}
                {skills.length === 0 && (
                    <p className={cls.Empty}>No skills yet. Create one to get started.</p>
                )}
            </div>
        </div>
    )
}
