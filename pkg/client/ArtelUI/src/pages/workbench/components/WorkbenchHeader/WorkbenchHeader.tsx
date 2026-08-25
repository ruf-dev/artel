import {Link} from "react-router-dom"

import cls from "@/pages/workbench/WorkbenchPage.module.css"
import {Path} from "@/app/routing/Router.tsx"
import SimpleChatTopBar from "@/pages/workbench/components/SimpleChatTopBar/SimpleChatTopBar.tsx"
import WorkbenchToolbar, {type WorkbenchView} from "@/pages/workbench/components/WorkbenchToolbar/WorkbenchToolbar.tsx"
import type {WorkbenchMode} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"

interface SimpleChatTopBarProps {
    models: string[]
    currentModel: string
    modelsLoading?: boolean
    onChangeModel: (model: string) => void
    onNewChat: () => void
    onToggleHistory: () => void
}

interface Props {
    isLoading: boolean
    effectiveMode: WorkbenchMode | "picking"
    exists: boolean
    vaultId?: string
    vaultName: string
    status: string
    onStart: () => void
    onStop: () => void
    stopping: boolean
    starting: boolean
    view: WorkbenchView
    onViewChange: (view: WorkbenchView) => void
    chatLocked: boolean
    onToggleHistory: () => void
    simpleChatTopBar: SimpleChatTopBarProps
}

// The docker toolbar / Simple Chat top bar / bare back-link 3-way switch, split
// out purely to keep WorkbenchPage.tsx's render function under the
// max-lines-per-function lint limit — same rationale as WorkbenchPanels.tsx.
//
// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function WorkbenchHeader(props: Props) {
    if (!props.isLoading && props.effectiveMode === "docker" && props.exists && props.vaultId) {
        return (
            <WorkbenchToolbar
                vaultName={props.vaultName}
                status={props.status}
                exists={props.exists}
                vaultId={props.vaultId}
                onStart={props.onStart}
                onStop={props.onStop}
                stopping={props.stopping}
                starting={props.starting}
                view={props.view}
                onViewChange={props.onViewChange}
                chatLocked={props.chatLocked}
                onToggleHistory={props.onToggleHistory}
            />
        )
    }

    if (!props.isLoading && props.effectiveMode === "simple-chat" && props.vaultId) {
        return (
            <SimpleChatTopBar
                vaultName={props.vaultName}
                models={props.simpleChatTopBar.models}
                currentModel={props.simpleChatTopBar.currentModel}
                modelsLoading={props.simpleChatTopBar.modelsLoading}
                onChangeModel={props.simpleChatTopBar.onChangeModel}
                onNewChat={props.simpleChatTopBar.onNewChat}
                onToggleHistory={props.simpleChatTopBar.onToggleHistory}
            />
        )
    }

    return <Link className={cls.BackLink} to={Path.HomePage}>← Back to vaults</Link>
}
