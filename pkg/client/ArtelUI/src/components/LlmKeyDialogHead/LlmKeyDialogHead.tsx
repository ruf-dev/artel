import {ModalClose} from "@vervstack/chures"

import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"
import cls from "@/components/LlmKeyDialogHead/LlmKeyDialogHead.module.css"

interface LlmKeyDialogHeadProps {
    provider: ExternalProvider
    title: string
    onClose: () => void
}

export default function LlmKeyDialogHead({provider, title, onClose}: LlmKeyDialogHeadProps) {
    return (
        <div className={cls.LlmKeyDialogHeadContainer}>
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
