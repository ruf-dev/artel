import {Button} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.module.css"
import {cn} from "@/app/utils/cn.ts"
import {CanvasNode, NODE_HEIGHT, NODE_WIDTH} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import {TractTool} from "@/processes/Tracts.ts"
import {PlusIcon} from "@/pages/tract-canvas/components/TractIcons/TractIcons.tsx"
import {colorForKind} from "@/pages/tract-canvas/components/TractIcons/tractIconHelpers.ts"
import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import NodeIcon from "@/pages/tract-canvas/icons/NodeIcon/NodeIcon.tsx"
import NodeChips from "@/pages/tract-canvas/components/NodeChips/NodeChips.tsx"
import {cap, title, typeLabel} from "@/pages/tract-canvas/components/TractCanvasNode/tractCanvasNode.ts"

export type NodeStatus = "ok" | "err" | "running" | "idle"

interface Props {
    node: CanvasNode
    tools: TractTool[]
    triggerInfo?: { name: string; kind: string; source: string }
    momCandidates: MomCandidate[]
    status: NodeStatus
    selected: boolean
    onClick: () => void
    onAddBlock: () => void
}

export default function TractCanvasNode(
    {
        node, tools, triggerInfo, momCandidates,
        status, selected, onClick, onAddBlock
    }: Props) {
    const step = node.step
    const color = colorForKind(node.kind, step?.mcp)

    return (
        <div
            className={cn(cls.Node, selected && cls.Selected, node.kind === "trigger" && cls.TriggerNode)}
            style={{left: node.x, top: node.y, width: NODE_WIDTH, height: NODE_HEIGHT}}
            data-tract-node
            onClick={e => {
                e.stopPropagation()
                onClick()
            }}
            role="button"
            tabIndex={0}
        >
            <span className={cn(cls.StatusDot, cls[`Status${cap(status)}`])}/>

            <div className={cls.Head}>
                <div className={cls.IconBox}
                     style={{
                         background: color.bg,
                         borderColor: color.border,
                         color: color.fg,
                     }}>
                    <NodeIcon node={node} triggerInfo={triggerInfo}/>
                </div>
                <div className={cls.Titles}>
                    <div className={cls.TypeLabel}>{typeLabel(node, triggerInfo)}</div>
                    <div className={cls.Title}>{title(node, triggerInfo)}</div>
                </div>
            </div>

            <div className={cls.Chips}>
                <NodeChips
                    node={node}
                    tools={tools}
                    triggerInfo={triggerInfo}
                    momCandidates={momCandidates}/>
            </div>
            <Button
                className={cls.AddButton}
                aria-label="Add block after"
                title="Add block after"
                onClick={e => {
                    e.stopPropagation()
                    onAddBlock()
                }}
            >
                <PlusIcon/>
            </Button>
        </div>
    )
}
