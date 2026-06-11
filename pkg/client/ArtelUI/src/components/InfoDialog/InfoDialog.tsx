import cls from "@/components/InfoDialog/InfoDialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"

interface Props {
    title: string
    message: string
}

export default function InfoDialog({title, message}: Props) {
    const {CloseDialog} = useDialog()

    return (
        <div className={cls.InfoContainer} role="dialog" aria-modal="true">
            <h2 className={cls.InfoTitle}>{title}</h2>
            <p className={cls.InfoMessage}>{message}</p>
            <div className={cls.InfoActions}>
                <button className={cls.BtnClose} type="button" onClick={CloseDialog}>
                    Got it
                </button>
            </div>
        </div>
    )
}
