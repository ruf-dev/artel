import cls from "@/widgets/TractCanvas/TractCanvasNode.module.css"

import {CanvasNode, NODE_WIDTH} from "@/processes/tractCanvasLayout.ts"
import {TractTool} from "@/processes/Tracts.ts"
import {
    BranchIcon,
    ChatIcon,
    ForkIcon,
    LayersIcon,
    ManualTriggerIcon,
    WebhookIcon,
} from "@/components/TractIcons/TractIcons.tsx"
import {colorForKind, iconForTool} from "@/components/TractIcons/tractIconHelpers.ts"

export type NodeStatus = "ok" | "err" | "running" | "idle"

interface Props {
    node: CanvasNode
    tools: TractTool[]
    triggerInfo?: {name: string; kind: string; source: string}
    status: NodeStatus
    selected: boolean
    onClick: () => void
}

export default function TractCanvasNode({node, tools, triggerInfo, status, selected, onClick}: Props) {
    const step = node.step
    const color = colorForKind(node.kind, step?.mcp)

    return (
        <div
            className={`${cls.Node} ${selected ? cls.Selected : ""} ${node.kind === "trigger" ? cls.TriggerNode : ""}`}
            style={{left: node.x, top: node.y, width: NODE_WIDTH}}
            onClick={e => {
                e.stopPropagation()
                onClick()
            }}
            role="button"
            tabIndex={0}
        >
            <span className={`${cls.StatusDot} ${cls[`Status${cap(status)}`]}`}/>
            <div className={cls.Head}>
                <div className={cls.IconBox} style={{background: color.bg, borderColor: color.border, color: color.fg}}>
                    <NodeIcon node={node} triggerInfo={triggerInfo}/>
                </div>
                <div className={cls.Titles}>
                    <div className={cls.TypeLabel}>{typeLabel(node, triggerInfo)}</div>
                    <div className={cls.Title}>{title(node, triggerInfo)}</div>
                </div>
            </div>
            <div className={cls.Chips}>
                <NodeChips node={node} tools={tools} triggerInfo={triggerInfo}/>
            </div>
        </div>
    )
}

function cap(s: string): string {
    return s.charAt(0).toUpperCase() + s.slice(1)
}

function NodeIcon({node, triggerInfo}: { node: CanvasNode; triggerInfo?: {kind: string} }) {
    if (node.kind === "trigger") {
        return triggerInfo?.kind === "webhook" ? <WebhookIcon/> : <ManualTriggerIcon/>
    }
    if (node.kind === "condition") return <BranchIcon/>
    if (node.kind === "parallel") return <ForkIcon/>
    if (node.kind === "group") return <LayersIcon/>
    const step = node.step
    if (!step) return <ChatIcon/>
    const Icon = iconForTool(step.mcp ?? "", step.tool ?? "")
    return <Icon/>
}

function typeLabel(node: CanvasNode, triggerInfo?: {kind: string; source: string}): string {
    if (node.kind === "trigger") return triggerInfo ? triggerInfo.kind : "trigger"
    if (node.kind === "condition") return "condition"
    if (node.kind === "parallel") return "parallel"
    if (node.kind === "group") return "group"
    const step = node.step
    if (step?.mcp && step.tool) return `${step.mcp}.${step.tool}`
    return "action"
}

function title(node: CanvasNode, triggerInfo?: {name: string}): string {
    if (node.kind === "trigger") return triggerInfo?.name ?? "No trigger linked"
    const step = node.step
    if (!step) return node.kind
    return step.name || step.id
}

function NodeChips({node, tools, triggerInfo}: {
    node: CanvasNode
    tools: TractTool[]
    triggerInfo?: {kind: string; source: string}
}) {
    if (node.kind === "trigger") {
        if (!triggerInfo) return <span className={cls.Chip}>unlinked</span>
        return <><span className={cls.Chip}>{triggerInfo.kind}</span><span className={cls.Chip}>{triggerInfo.source}</span></>
    }

    const step = node.step
    if (!step) return null

    if (node.kind === "condition") {
        const n = step.conditions?.length ?? 0
        return <span className={cls.Chip}>{n} {n === 1 ? "condition" : "conditions"}</span>
    }
    if (node.kind === "parallel") {
        const n = step.steps?.length ?? 0
        return <span className={cls.Chip}>{n === 0 ? "no lanes yet" : `${n} lanes`}</span>
    }
    if (node.kind === "group") {
        const n = step.steps?.length ?? 0
        return <span className={cls.Chip}>{n} {n === 1 ? "step" : "steps"}</span>
    }

    const tool = tools.find(t => t.mcp === step.mcp && t.tool === step.tool)
    if (!tool) return <span className={cls.Chip}>unknown tool</span>

    const inputs = Object.keys(tool.inputSchema.properties)
    const outputs = Object.keys(tool.outputSchema.properties)
    return (
        <>
            {inputs.map(n => <span key={`in-${n}`} className={`${cls.Chip} ${cls.ChipIn}`}>{n}</span>)}
            {outputs.map(n => <span key={`out-${n}`} className={`${cls.Chip} ${cls.ChipOut}`}>{n}</span>)}
        </>
    )
}
