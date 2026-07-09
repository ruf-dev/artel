import {Button} from "@vervstack/chures"

import cls from "@/dialogs/ConnectionDetailDialog/components/NotConnectedContent/NotConnectedContent.module.css"

interface NotConnectedContentProps {
    description: string
    canConnect: boolean
    name: string
    connecting: boolean
    onConnect: () => void
}

export default function NotConnectedContent({description, canConnect, name, connecting, onConnect}: NotConnectedContentProps) {
    return (
        <div className={cls.NotConnectedContentContainer}>
            <p className={cls.ModalDesc}>{description}</p>
            {canConnect ? (
                <div className={cls.ModalActions}>
                    <Button variant="primary" onClick={onConnect} disabled={connecting}>
                        {connecting ? "Redirecting…" : `Connect with ${name}`}
                    </Button>
                </div>
            ) : (
                <p className={cls.ComingSoon}>Coming soon — this integration is not yet available.</p>
            )}
        </div>
    )
}
