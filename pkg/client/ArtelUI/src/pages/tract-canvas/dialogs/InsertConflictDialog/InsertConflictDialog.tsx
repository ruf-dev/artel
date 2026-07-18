import {Button} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/dialogs/InsertConflictDialog/InsertConflictDialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {MoveConflictChoice} from "@/processes/tractStepsMove.ts"

interface Props {
    sourceName: string
    targetName: string
    onChoose: (choice: MoveConflictChoice) => void
}

export default function InsertConflictDialog({sourceName, targetName, onChoose}: Props) {
    const {CloseDialog} = useDialog()

    function choose(choice: MoveConflictChoice) {
        onChoose(choice)
        CloseDialog()
    }

    return (
        <div className={cls.InsertConflictDialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>{`"${targetName}" already has a next step`}</h2>
            <p className={cls.DialogMessage}>
                {`Choose how "${sourceName}" should join the flow after "${targetName}".`}
            </p>
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={CloseDialog}>Cancel</Button>
                <Button variant="secondary" onClick={() => choose("parallel")}>Run in parallel</Button>
                <Button variant="primary" onClick={() => choose("between")}>Insert between</Button>
            </div>
        </div>
    )
}
