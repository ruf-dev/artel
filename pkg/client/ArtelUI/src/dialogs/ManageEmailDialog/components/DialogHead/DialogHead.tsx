import {ModalClose} from "@vervstack/chures"

import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {mailProviderIcon} from "@/app/utils/mailProviderIcon"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"
import cls from "@/dialogs/ManageEmailDialog/components/DialogHead/DialogHead.module.css"

interface DialogHeadProps {
    title: string
    onClose: () => void
    disabled?: boolean
    email?: string
}

export default function DialogHead({title, onClose, disabled, email}: DialogHeadProps) {
    const icon = email ? mailProviderIcon(email) : null

    return (
        <div className={cls.DialogHeadContainer}>
            <div className={cls.ModalHeadLeft}>
                <div className={cls.ModalIcon}>
                    {icon
                        ? <img src={icon} className={cls.ModalProviderIcon} alt=""/>
                        : <ProviderIcon provider={ExternalProvider.EXTERNAL_PROVIDER_EMAIL}/>
                    }
                </div>
                <span className={cls.ModalTitle}>{title}</span>
            </div>
            <ModalClose onClick={onClose} disabled={disabled} className={cls.ModalClose}/>
        </div>
    )
}
