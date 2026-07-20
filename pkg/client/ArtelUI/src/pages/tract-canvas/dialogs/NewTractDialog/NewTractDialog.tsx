import {useState} from "react"
import {useNavigate} from "react-router-dom"
import {Button, Input} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/dialogs/NewTractDialog/NewTractDialog.module.css"
import {useDialog, useDialogKeyboard} from "@/app/hooks/Dialog"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import BrowseTemplatesDialog from "@/dialogs/BrowseTemplatesDialog/BrowseTemplatesDialog.tsx"


export default function NewTractDialog() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {createTract} = useTracts()
    const bakeError = useBakeError()
    const navigate = useNavigate()
    const [name, setName] = useState("")
    const [creating, setCreating] = useState(false)

    function handleCreate() {
        const trimmed = name.trim()
        if (!trimmed) return
        setCreating(true)
        createTract(trimmed, "", {steps: []})
            .then(({tract}) => {
                CloseDialog()
                navigate(`/tracts/${tract.uuid}`)
            })
            .catch(err => bakeError("Failed to create tract", err))
            .finally(() => setCreating(false))
    }

    const onKeyDown = useDialogKeyboard(handleCreate)

    return (
        <div className={cls.NewTractDialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>New tract</h2>
            <Input
                inputClassName={cls.DialogInput}
                value={name}
                setValue={setName}
                onKeyDown={onKeyDown}
                placeholder="Tract name"
                autoFocus
                disabled={creating}
            />
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={CloseDialog} disabled={creating}>Cancel</Button>
                <Button variant="primary" onClick={handleCreate} disabled={creating || !name.trim()}>
                    {creating ? "…" : "Create"}
                </Button>
            </div>
            <div className={cls.ExploreLinkWrapper}>
                <Button
                    variant="ghost"
                    className={cls.ExploreLink}
                    onClick={() => OpenDialog(<BrowseTemplatesDialog/>)}
                    disabled={creating}
                >
                    or explore templates →
                </Button>
            </div>
        </div>
    )
}
