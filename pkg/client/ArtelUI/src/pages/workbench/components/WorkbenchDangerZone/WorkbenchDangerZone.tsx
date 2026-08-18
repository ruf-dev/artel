import {Button, ConfirmDialog} from "@vervstack/chures"

import cls from "@/pages/workbench/components/WorkbenchDangerZone/WorkbenchDangerZone.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useWorkbenchMutations} from "@/app/hooks/Workbench.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import DangerZoneText from "@/pages/workbench/components/WorkbenchDangerZone/components/DangerZoneText.tsx"

interface Props {
    vaultId: string
}

export default function WorkbenchDangerZone({vaultId}: Props) {
    const {CloseDialog, OpenDialog} = useDialog()
    const {remove, removing} = useWorkbenchMutations(vaultId)
    const bakeError = useBakeError()

    function handleDelete() {
        OpenDialog(
            <ConfirmDialog
                title="Delete Workbench"
                message="Delete this Workbench? This destroys its container and workspace — all
                    files and state saved in it are lost. This cannot be undone."
                confirmLabel="Delete"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() => remove()
                    .then(() => undefined)
                    .catch(e => bakeError("Failed to delete workbench", e))}
            />
        )
    }

    return (
        <section className={cls.WorkbenchDangerZoneContainer}>
            <div className={cls.WorkbenchDangerZoneRow}>
                <DangerZoneText/>
                <Button variant="danger" onClick={handleDelete} disabled={removing}>
                    {removing ? "Deleting…" : "Delete"}
                </Button>
            </div>
        </section>
    )
}
