import {useEffect, useState} from "react"
import {Button, Loader, ModalClose} from "@vervstack/chures"

import cls from "@/dialogs/EditSkillDialog/EditSkillDialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useSkills} from "@/app/hooks/Skills.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {skillsService} from "@/processes/Skills.ts"
import FormField from "@/components/FormField/FormField.tsx"
import StorageModeSelect from "@/components/StorageModeSelect/StorageModeSelect.tsx"
import Textarea from "@/components/atoms/Textarea/Textarea.tsx"

interface Props {
    vaultId: string
    slug: string
}

export default function EditSkillDialog({vaultId, slug}: Props) {
    const [loading, setLoading] = useState(true)
    const [saving, setSaving] = useState(false)
    const [name, setName] = useState("")
    const [description, setDescription] = useState("")
    const [storageMode, setStorageMode] = useState("none")
    const [body, setBody] = useState("")

    const {update} = useSkills()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

    useEffect(() => {
        skillsService.get(vaultId, slug)
            .then(res => {
                setName(res.skill.name ?? "")
                setDescription(res.skill.description ?? "")
                setStorageMode(res.skill.storageMode ?? "none")
                setBody(res.body)
            })
            .catch(err => {
                bakeError("Failed to load skill", err)
                CloseDialog()
            })
            .finally(() => setLoading(false))
    }, [vaultId, slug])

    function handleSave() {
        if (!name) return
        setSaving(true)
        update(vaultId, slug, name, description, storageMode, body)
            .then(CloseDialog)
            .catch(err => bakeError("Failed to update skill", err))
            .finally(() => setSaving(false))
    }

    return (
        <div className={cls.EditSkillDialogContainer} onClick={e => e.stopPropagation()} role="dialog"
             aria-modal="true" aria-labelledby="editSkillTitle">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle} id="editSkillTitle">Edit skill</h2>
                <ModalClose onClick={CloseDialog} disabled={saving} className={cls.ModalClose}/>
            </div>

            {loading ? (
                <div className={cls.LoadingWrap}>
                    <Loader variant="arcs" size="sm" color="var(--coral)"/>
                </div>
            ) : (
                <>
                    <FormField
                        label="Name"
                        placeholder="e.g. Release notes"
                        defaultValue={name}
                        onChange={setName}
                        disabled={saving}
                        fieldClassName={cls.Field}
                        labelClassName={cls.FieldLabel}
                    />

                    <label className={cls.Field}>
                        <span className={cls.FieldLabel}>Description</span>
                        <Textarea
                            value={description}
                            setValue={setDescription}
                            disabled={saving}
                            rows={2}
                            placeholder="Short — this is what triggers the skill by keyword match."
                        />
                    </label>

                    <div className={cls.Field}>
                        <span className={cls.FieldLabel}>Storage mode</span>
                        <StorageModeSelect value={storageMode} onChange={setStorageMode} disabled={saving}/>
                    </div>

                    <label className={cls.Field}>
                        <span className={cls.FieldLabel}>Body</span>
                        <Textarea
                            value={body}
                            setValue={setBody}
                            disabled={saving}
                            rows={10}
                            className={cls.BodyTextarea}
                            placeholder="Instructional markdown the agent reads when the skill runs…"
                        />
                    </label>

                    <div className={cls.ModalActions}>
                        <Button variant="primary" onClick={handleSave} disabled={saving || !name}>
                            {saving ? "Saving…" : "Save changes"}
                        </Button>
                    </div>
                </>
            )}
        </div>
    )
}
