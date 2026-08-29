import TopbarLeft from "@/pages/workbench/components/WorkbenchTopbar/components/TopbarLeft/TopbarLeft.tsx"
// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import TopbarNavToggle from "@/pages/workbench/components/WorkbenchTopbar/components/TopbarNavToggle/TopbarNavToggle.tsx"
import TopbarRight from "@/pages/workbench/components/WorkbenchTopbar/components/TopbarRight/TopbarRight.tsx"
import type {WorkbenchMode} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"
import type {WorkbenchView} from "@/pages/workbench/processes/workbenchView.ts"
import cls from "@/pages/workbench/components/WorkbenchTopbar/WorkbenchTopbar.module.css"

interface ModelProps {
    models: string[]
    value: string
    isLoading?: boolean
    onChange: (model: string) => void
}

interface Props {
    effectiveMode: WorkbenchMode | "picking"
    exists: boolean
    status: string
    vaultId?: string
    view: WorkbenchView
    onViewChange: (view: WorkbenchView) => void
    chatLocked: boolean
    onStart: () => void
    onStop: () => void
    stopping: boolean
    starting: boolean
    model: ModelProps
    onToggleNav: () => void
    showNavToggle: boolean
    onNewChat: () => void
}

// The workbench top bar: status + model switcher on the left, the New Chat button
// + Chat/Terminal switch + start/stop + settings on the right. Renders for api mode
// and for a provisioned docker workbench — the loading/picking states are gated out
// by the caller, and the internal guard drops docker-without-a-workbench.
//
// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function WorkbenchTopbar(props: Props) {
    const visible = props.effectiveMode === "api"
        || (props.effectiveMode === "docker" && props.exists)
    if (!visible) return null

    return (
        <div className={cls.WorkbenchTopbarContainer}>
            {props.showNavToggle && <TopbarNavToggle onToggle={props.onToggleNav}/>}
            <TopbarLeft
                effectiveMode={props.effectiveMode}
                exists={props.exists}
                status={props.status}
                model={props.model}
            />
            <div className={cls.Spacer}/>
            <TopbarRight
                effectiveMode={props.effectiveMode}
                exists={props.exists}
                status={props.status}
                vaultId={props.vaultId}
                view={props.view}
                onViewChange={props.onViewChange}
                chatLocked={props.chatLocked}
                onStart={props.onStart}
                onStop={props.onStop}
                stopping={props.stopping}
                starting={props.starting}
                onNewChat={props.onNewChat}
            />
        </div>
    )
}
