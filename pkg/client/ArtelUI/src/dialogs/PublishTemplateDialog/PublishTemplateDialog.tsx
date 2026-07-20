import {useState} from "react"
import {Button, Input, useToaster} from "@vervstack/chures"

import cls from "@/dialogs/PublishTemplateDialog/PublishTemplateDialog.module.css"
import {useDialog, useDialogKeyboard} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {tractsService} from "@/processes/Tracts.ts"

interface Props {
    tractUuid: string
}

export default function PublishTemplateDialog({tractUuid}: Props) {
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const {bake} = useToaster()
    const [category, setCategory] = useState("")
    const [publishing, setPublishing] = useState(false)

    function handlePublish() {
        setPublishing(true)
        tractsService.publishTemplate(tractUuid, category.trim())
            .then(template => {
                CloseDialog()
                bake({
                    title: "Template published",
                    description: `"${template.name}" is now visible in the public template gallery.`,
                    level: "Info",
                })
            })
            .catch(err => bakeError("Failed to publish template", err))
            .finally(() => setPublishing(false))
    }

    const onKeyDown = useDialogKeyboard(handlePublish)

    return (
        <div className={cls.PublishTemplateDialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Publish as template</h2>
            <p className={cls.Warning}>
                This will be publicly visible to everyone, including any script code and literal
                values in your steps — review before publishing.
            </p>
            <Input
                inputClassName={cls.DialogInput}
                value={category}
                setValue={setCategory}
                onKeyDown={onKeyDown}
                placeholder="Category (optional)"
                autoFocus
                disabled={publishing}
            />
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={CloseDialog} disabled={publishing}>Cancel</Button>
                <Button variant="primary" onClick={handlePublish} disabled={publishing}>
                    {publishing ? "…" : "Publish"}
                </Button>
            </div>
        </div>
    )
}
