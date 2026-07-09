import {ModalClose} from "@vervstack/chures"

import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"
import cls from "@/dialogs/ManageGitlabDialog/components/DialogHead/DialogHead.module.css"

export default function DialogHead({onClose}: { onClose: () => void }) {
    return (
        <div className={cls.DialogHeadContainer}>
            <div className={cls.ModalHeadLeft}>
                <div className={cls.ModalIcon}>
                    <ProviderIcon provider={ExternalProvider.EXTERNAL_PROVIDER_GITLAB}/>
                </div>
                <span className={cls.ModalTitle}>GitLab</span>
            </div>
            <ModalClose onClick={onClose} className={cls.ModalClose}/>
        </div>
    )
}
