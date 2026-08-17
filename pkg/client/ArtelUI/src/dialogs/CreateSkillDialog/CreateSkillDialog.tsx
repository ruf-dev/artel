import {useState} from "react"
import {Button, ModalClose, Toggle} from "@vervstack/chures"

import cls from "@/dialogs/CreateSkillDialog/CreateSkillDialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useSkills} from "@/app/hooks/Skills.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import FormField from "@/components/FormField/FormField.tsx"
import StorageModeSelect from "@/components/StorageModeSelect/StorageModeSelect.tsx"
import Textarea from "@/components/atoms/Textarea/Textarea.tsx"

interface Props {
    vaultId: string
}

export default function CreateSkillDialog({vaultId}: Props) {
    const [creating, setCreating] = useState(false)
    const [name, setName] = useState("")
    const [description, setDescription] = useState("")
    const [storageMode, setStorageMode] = useState("none")
    const [body, setBody] = useState("")
    const [hotPlug, setHotPlug] = useState(false)

    const {create} = useSkills()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

    function handleCreate() {
        if (!name) return
        setCreating(true)
        create(vaultId, name, description, storageMode, body, hotPlug)
            .then(CloseDialog)
            .catch(err => bakeError("Failed to create skill", err))
            .finally(() => setCreating(false))
    }

    return (
        <div className={cls.CreateSkillDialogContainer} onClick={e => e.stopPropagation()} role="dialog"
             aria-modal="true" aria-labelledby="createSkillTitle">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle} id="createSkillTitle">New skill</h2>
                <ModalClose onClick={CloseDialog} disabled={creating} className={cls.ModalClose}/>
            </div>
            <p className={cls.ModalSub}>Skills are instructions your agent can pull in on demand.</p>

            <FormField
                label="Name"
                placeholder="e.g. Release notes"
                onChange={setName}
                disabled={creating}
                fieldClassName={cls.Field}
                labelClassName={cls.FieldLabel}
            />

            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Description</span>
                <Textarea
                    value={description}
                    setValue={setDescription}
                    disabled={creating}
                    rows={2}
                    placeholder="Short — this is what triggers the skill by keyword match."
                />
            </label>

            <div className={cls.Field}>
                <span className={cls.FieldLabel}>Storage mode</span>
                <StorageModeSelect value={storageMode} onChange={setStorageMode} disabled={creating}/>
            </div>

            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Body</span>
                <Textarea
                    value={body}
                    setValue={setBody}
                    disabled={creating}
                    rows={10}
                    className={cls.BodyTextarea}
                    placeholder="Instructional markdown the agent reads when the skill runs…"
                />
            </label>

            <div className={cls.Field}>
                <Toggle checked={hotPlug} onChange={setHotPlug} disabled={creating} label="Hot-plug"/>
                <p className={cls.Caption}>Counts against your hot-plug slots — always available to the agent.</p>
            </div>

            <div className={cls.ModalActions}>
                <Button variant="primary" onClick={handleCreate} disabled={creating || !name}>
                    {creating ? "Creating…" : "Create skill"}
                </Button>
            </div>
        </div>
    )
}
