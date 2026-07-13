import {ModalClose} from "@vervstack/chures"

import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"
import cls from "@/dialogs/ManageTrelloDialog/components/DialogHead/DialogHead.module.css"

interface DialogHeadProps {
    title?: string
    onClose: () => void
    disabled?: boolean
}

export default function DialogHead({title = "Trello", onClose, disabled}: DialogHeadProps) {
    return (
        <div className={cls.DialogHeadContainer}>
            <div className={cls.ModalHeadLeft}>
                <div className={cls.ModalIcon}>
                    <ProviderIcon provider={ExternalProvider.EXTERNAL_PROVIDER_TRELLO}/>
                </div>
                <span className={cls.ModalTitle}>{title}</span>
            </div>
            <ModalClose onClick={onClose} disabled={disabled} className={cls.ModalClose}/>
        </div>
    )
}
