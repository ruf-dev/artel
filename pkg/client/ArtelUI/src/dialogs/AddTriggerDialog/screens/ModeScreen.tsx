import {Button} from "@vervstack/chures"
import cls from "@/dialogs/AddTriggerDialog/AddTriggerDialog.module.css"

import {cn} from "@/app/utils/cn.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useAddTriggerDialog} from "@/dialogs/AddTriggerDialog/AddTriggerDialogContext.ts"

export default function ModeScreen() {
    const {CloseDialog} = useDialog()
    const context = useAddTriggerDialog()

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Add trigger</h2>
            <div className={cls.Body}>
                <Button variant="ghost" className={cls.OptionRow} onClick={() => context.setStep("create")}>
                    <span className={cls.Name}>Create new trigger</span>
                </Button>
                <Button variant="ghost" className={cls.OptionRow} onClick={() => context.setStep("link")}>
                    <span className={cls.Name}>Link existing trigger</span>
                </Button>
            </div>
            <div className={cn(cls.DialogActions, cls.DialogActionsEnd)}>
                <Button variant="ghost" onClick={CloseDialog}>Cancel</Button>
            </div>
        </div>
    )
}
