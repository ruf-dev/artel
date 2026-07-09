import {type ReactNode} from "react"
import {ModalClose} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import cls from "@/dialogs/ManageKeyDialog/components/DialogHead/DialogHead.module.css"

interface DialogHeadProps {
    titleId: string
    title: ReactNode
    disabled?: boolean
}

export default function DialogHead({titleId, title, disabled}: DialogHeadProps) {
    const {CloseDialog} = useDialog()
    return (
        <div className={cls.DialogHeadContainer}>
            <h2 className={cls.ModalTitle} id={titleId}>{title}</h2>
            <ModalClose onClick={CloseDialog} disabled={disabled} className={cls.ModalClose}/>
        </div>
    )
}
