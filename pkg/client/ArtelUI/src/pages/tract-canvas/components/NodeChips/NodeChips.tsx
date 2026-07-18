import cls from "@/pages/tract-canvas/components/NodeChips/NodeChips.module.css"
import {CanvasNode} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import {TractTool} from "@/processes/Tracts.ts"
import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"
import {cn} from "@/app/utils/cn.ts"

const BUILTIN_MCP = "artel"

export default function NodeChips({node, tools, triggerInfo, momCandidates}: {
    node: CanvasNode
    tools: TractTool[]
    triggerInfo?: { kind: string; source: string }
    momCandidates: MomCandidate[]
}) {
    if (node.kind === "trigger") {
        if (!triggerInfo) return <span className={cls.Chip}>unlinked</span>

        return <>
            <span className={cls.Chip}>{triggerInfo.source}</span>
        </>
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
    if (node.kind === "script") {
        const inN = step.inputParams?.length ?? 0
        const outN = step.outputParams?.length ?? 0
        return <span className={cls.Chip}>{inN} in · {outN} out</span>
    }

    const tool = tools.find(t => t.mcp === step.mcp && t.tool === step.tool)
    if (!tool) return <span className={cls.Chip}>unknown tool</span>

    if (step.mcp === BUILTIN_MCP) return <span className={cls.Chip}>built-in</span>

    const candidate = momCandidates.find(c => c.name === step.mcp)
    const connection = candidate?.connections?.find(c => c.id === step.connection_uuid)
    if (!connection) return <span className={cn(cls.Chip, cls.ChipWarn)}>no connection</span>

    return (
        <span className={cn(cls.Chip, cls.ChipConn)}>
            <span className={cls.ChipIcon}>
                <ProviderIcon provider={connection.provider}/>
            </span>

            {connectionLabel(connection)}
        </span>
    )
}
