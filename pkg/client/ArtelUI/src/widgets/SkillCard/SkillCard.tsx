import {Button, ConfirmDialog, Toggle} from "@vervstack/chures"

import cls from "@/widgets/SkillCard/SkillCard.module.css"
import {SkillInfo} from "@/app/api/artel/skills.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useSkills} from "@/app/hooks/Skills.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {cn} from "@/app/utils/cn.ts"
import EditSkillDialog from "@/dialogs/EditSkillDialog/EditSkillDialog.tsx"

const STORAGE_MODE_LABEL: Record<string, string> = {
    none: "none",
    freeform_notes: "freeform notes",
    structured: "structured",
}

interface Props {
    vaultId: string
    skill: SkillInfo
}

export default function SkillCard({vaultId, skill}: Props) {
    const {setHotPlug, remove} = useSkills()
    const bakeError = useBakeError()
    const {OpenDialog, CloseDialog} = useDialog()

    const slug = skill.slug ?? ""
    const isSystem = skill.isSystem ?? false

    function handleToggle() {
        setHotPlug(vaultId, slug, !skill.isHotPlug)
            .catch(err => bakeError("Failed to update skill", err))
    }

    function handleEdit(e: React.MouseEvent) {
        e.stopPropagation()
        OpenDialog(<EditSkillDialog vaultId={vaultId} slug={slug}/>)
    }

    function handleDelete(e: React.MouseEvent) {
        e.stopPropagation()
        OpenDialog(
            <ConfirmDialog
                title="Delete skill"
                message={`Delete "${skill.name}"? This cannot be undone.`}
                confirmLabel="Delete"
                danger
                onClose={CloseDialog}
                onConfirm={() => remove(vaultId, slug).catch(err => bakeError("Failed to delete skill", err))}
            />
        )
    }

    return (
        <div className={cls.SkillCardContainer}>
            <div className={cls.Info}>
                <span className={cls.Name}>{skill.name}</span>
                {skill.description && <span className={cls.Desc}>{skill.description}</span>}
            </div>
            <div className={cls.Meta}>
                {isSystem ? (
                    <span className={cn(cls.Chip, cls.ChipBuiltin)}>Built-in</span>
                ) : (
                    <span className={cls.Chip}>{STORAGE_MODE_LABEL[skill.storageMode ?? ""] ?? skill.storageMode}</span>
                )}
                <div className={cls.ToggleWrap} onClick={e => e.stopPropagation()}>
                    <Toggle checked={skill.isHotPlug ?? false} onChange={handleToggle} disabled={isSystem}/>
                </div>
            </div>
            {!isSystem && (
                <div className={cls.Actions}>
                    <Button className={cls.IconBtn} variant="ghost" onClick={handleEdit} aria-label="Edit skill">
                        <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor"
                             strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M12 20h9"/>
                            <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4Z"/>
                        </svg>
                    </Button>
                    <Button
                        className={cls.IconBtn} variant="iconDanger" onClick={handleDelete}
                        aria-label="Delete skill"
                    >
                        <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor"
                             strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M3 6h18"/>
                            <path
                                d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2m3 0-1 14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2L4 6"
                            />
                        </svg>
                    </Button>
                </div>
            )}
        </div>
    )
}
