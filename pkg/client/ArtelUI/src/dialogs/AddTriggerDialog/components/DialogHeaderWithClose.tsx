import {Button} from "@vervstack/chures"

import cls from "@/dialogs/AddTriggerDialog/AddTriggerDialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import CloseIcon from "@/icons/common/CloseIcon.tsx"

export default function DialogHeaderWithClose({title}: { title: string }) {
    const {CloseDialog} = useDialog()

    return (
        <div className={cls.DialogHeader}>
            <h2 className={cls.DialogTitle}>{title}</h2>
            <Button variant="ghost" className={cls.CloseBtn} onClick={CloseDialog} aria-label="Close">
                <CloseIcon className={cls.CloseIcon}/>
            </Button>
        </div>
    )
}
