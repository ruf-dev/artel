import {Button} from "@vervstack/chures"

import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"
import cls
    from "@/dialogs/ManageGitlabDialog/components/ConnectForm/components/UseDefaultButton/UseDefaultButton.module.css"

export default function UseDefaultButton({onClick, disabled}: { onClick: () => void, disabled: boolean }) {
    return (
        <Button variant="ghost" className={cls.UseDefaultButtonContainer} onClick={onClick} disabled={disabled}
            data-tooltip-id="root-tooltip" data-tooltip-content="Use default" aria-label="Use default instance URL">
            <span className={cls.Icon}>
                <ProviderIcon provider={ExternalProvider.EXTERNAL_PROVIDER_GITLAB}/>
            </span>
        </Button>
    )
}
