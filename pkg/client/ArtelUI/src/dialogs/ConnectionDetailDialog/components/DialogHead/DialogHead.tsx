import {ModalClose} from "@vervstack/chures"

import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"
import cls from "@/dialogs/ConnectionDetailDialog/components/DialogHead/DialogHead.module.css"

interface DialogHeadProps {
    title: string
    provider: ExternalProvider
    onClose: () => void
}

export default function DialogHead({title, provider, onClose}: DialogHeadProps) {
    return (
        <div className={cls.DialogHeadContainer}>
            <div className={cls.ModalHeadLeft}>
                <div className={cls.ModalIcon}>
                    <ProviderIcon provider={provider}/>
                </div>
                <span className={cls.ModalTitle}>{title}</span>
            </div>
            <ModalClose onClick={onClose} className={cls.ModalClose}/>
        </div>
    )
}
