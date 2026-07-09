import {CanvasNode} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import {BranchIcon} from "@/pages/tract-canvas/icons/BranchIcon/BranchIcon.tsx"
import {ChatIcon} from "@/pages/tract-canvas/icons/ChatIcon/ChatIcon.tsx"
import {ForkIcon} from "@/pages/tract-canvas/icons/ForkIcon/ForkIcon.tsx"
import {LayersIcon} from "@/pages/tract-canvas/icons/LayersIcon/LayersIcon.tsx"
import {ManualTriggerIcon} from "@/pages/tract-canvas/icons/ManualTriggerIcon/ManualTriggerIcon.tsx"
import {WebhookIcon} from "@/pages/tract-canvas/icons/WebhookIcon/WebhookIcon.tsx"
import {iconForTool} from "@/pages/tract-canvas/icons/tractIconHelpers.ts"

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
