import {Button, ConfirmDialog} from "@vervstack/chures"

import cls from "@/widgets/CommunityConnectorCard/CommunityConnectorCard.module.css"
import {cn} from "@/app/utils/cn.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

interface CommunityConnectorCardProps {
    name: string
    author: string
    description: string
    toolCount: number
    isConnected: boolean
    viewerIsOwner: boolean
    onConnect: () => void
    onDelete: () => Promise<void>
}

export default function CommunityConnectorCard(props: CommunityConnectorCardProps) {
    const {OpenDialog, CloseDialog} = useDialog()
    const bakeError = useBakeError()

    function handleDelete() {
        OpenDialog(
            <ConfirmDialog
                title="Delete connector"
                message={`Delete "${props.name}"? Anyone connected to it will lose access. This cannot be undone.`}
                confirmLabel="Delete"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() => props.onDelete().catch(e => bakeError("Failed to delete connector", e))}
            />
        )
    }

    const statusCls = cn(cls.StatusDot, props.isConnected ? cls.StatusDotConnected : cls.StatusDotDisconnected)
    const labelCls = cn(cls.StatusLabel, props.isConnected ? cls.StatusLabelConnected : cls.StatusLabelDisconnected)

    return (
        <div className={cls.CommunityConnectorCardContainer}>
            <div className={cls.Header}>
                <div className={cls.Name}>{props.name}</div>
                <div className={cls.Author}>by {props.author || "unknown"}</div>
            </div>
            <p className={cls.Description}>{props.description || "No description provided."}</p>
            <div className={cls.Footer}>
                <span className={cls.ToolCount}>{props.toolCount} {props.toolCount === 1 ? "tool" : "tools"}</span>
                <span className={statusCls}/>
                <span className={labelCls}>{props.isConnected ? "Connected" : "Not connected"}</span>
            </div>
            <div className={cls.Actions}>
                <Button variant="primary" onClick={props.onConnect}>Connect</Button>
                {props.viewerIsOwner && <Button variant="danger" onClick={handleDelete}>Delete</Button>}
            </div>
        </div>
    )
}
