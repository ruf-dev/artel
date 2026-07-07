import {CanvasNode} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import {
    BranchIcon,
    ChatIcon,
    ForkIcon,
    LayersIcon,
    ManualTriggerIcon,
    WebhookIcon,
} from "@/pages/tract-canvas/components/TractIcons/TractIcons.tsx"
import {iconForTool} from "@/pages/tract-canvas/components/TractIcons/tractIconHelpers.ts"

export default function NodeIcon({node, triggerInfo}: { node: CanvasNode; triggerInfo?: { kind: string } }) {
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
