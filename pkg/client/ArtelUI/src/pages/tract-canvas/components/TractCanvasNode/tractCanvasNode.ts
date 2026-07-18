import {CanvasNode} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"

export function cap(s: string): string {
    return s.charAt(0).toUpperCase() + s.slice(1)
}

export function typeLabel(node: CanvasNode, triggerInfo?: { kind: string; source: string }): string {
    if (node.kind === "trigger") return triggerInfo ? triggerInfo.kind : "trigger"
    if (node.kind === "condition") return "condition"
    if (node.kind === "parallel") return "parallel"
    if (node.kind === "group") return "group"
    if (node.kind === "script") return "script"
    const step = node.step
    if (step?.mcp && step.tool) return `${step.mcp}.${step.tool}`
    return "action"
}

export function title(node: CanvasNode, triggerInfo?: { name: string }): string {
    if (node.kind === "trigger") return triggerInfo?.name ?? "No trigger linked"
    const step = node.step
    if (!step) return node.kind
    return step.name || step.id
}
