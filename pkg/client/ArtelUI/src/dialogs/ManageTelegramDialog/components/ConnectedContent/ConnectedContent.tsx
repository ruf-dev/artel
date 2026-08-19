import {Button} from "@vervstack/chures"

import cls from "@/dialogs/ManageTelegramDialog/components/ConnectedContent/ConnectedContent.module.css"

interface ConnectedContentProps {
    botUsername: string
    onDisconnect: () => void
}

export default function ConnectedContent({botUsername, onDisconnect}: ConnectedContentProps) {
    return (
        <div className={cls.ConnectedContentContainer}>
            <p className={cls.ModalSub}>
                Connected as <b>@{botUsername}</b>. Messages you send this bot also relay into your workbench chat.
            </p>
            <div className={cls.ModalActions}>
                <Button variant="danger" onClick={onDisconnect}>Disconnect</Button>
            </div>
        </div>
    )
}
