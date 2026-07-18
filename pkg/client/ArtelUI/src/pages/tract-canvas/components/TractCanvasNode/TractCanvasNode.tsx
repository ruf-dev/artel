import type {DragEvent} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.module.css"
import {cn} from "@/app/utils/cn.ts"
import {CanvasNode, NODE_HEIGHT, NODE_WIDTH} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import {TractTool} from "@/processes/Tracts.ts"
import {PlusIcon} from "@/pages/tract-canvas/icons/PlusIcon/PlusIcon.tsx"
import {colorForKind} from "@/pages/tract-canvas/icons/tractIconHelpers.ts"
import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import NodeIcon from "@/pages/tract-canvas/icons/NodeIcon/NodeIcon.tsx"
import NodeChips from "@/pages/tract-canvas/components/NodeChips/NodeChips.tsx"
import {cap, title, typeLabel} from "@/pages/tract-canvas/components/TractCanvasNode/tractCanvasNode.ts"

export type NodeStatus = "ok" | "err" | "running" | "idle"
export type DragOverSide = "before" | "after"

interface Props {
    node: CanvasNode
    tools: TractTool[]
    triggerInfo?: { name: string; kind: string; source: string }
    momCandidates: MomCandidate[]
    status: NodeStatus
    selected: boolean
    dragging: boolean
    dragOverSide: DragOverSide | null
    onClick: () => void
    onAddBlock: () => void
    onDragStart: (e: DragEvent<HTMLDivElement>) => void
    onDragOver: (e: DragEvent<HTMLDivElement>) => void
    onDragLeave: (e: DragEvent<HTMLDivElement>) => void
    onDrop: (e: DragEvent<HTMLDivElement>) => void
    onDragEnd: (e: DragEvent<HTMLDivElement>) => void
}

export default function TractCanvasNode(props: Props) {
    const step = props.node.step
    const color = colorForKind(props.node.kind, step?.mcp)

    return (
        <div
            className={cn(
                cls.Node,
                props.selected && cls.Selected,
                props.node.kind === "trigger" && cls.TriggerNode,
                props.status === "running" && cls.Running,
                props.dragging && cls.Dragging,
                props.dragOverSide === "before" && cls.DropBefore,
                props.dragOverSide === "after" && cls.DropAfter,
            )}
            style={{left: props.node.x, top: props.node.y, width: NODE_WIDTH, height: NODE_HEIGHT}}
            data-tract-node
            draggable={props.node.kind !== "trigger"}
            onClick={e => {
                e.stopPropagation()
                props.onClick()
            }}
            onDragStart={props.onDragStart}
            onDragOver={props.onDragOver}
            onDragLeave={props.onDragLeave}
            onDrop={props.onDrop}
            onDragEnd={props.onDragEnd}
            role="button"
            tabIndex={0}
        >
            <span className={cn(cls.StatusDot, cls[`Status${cap(props.status)}`])}/>

            <div className={cls.Head}>
                <div className={cls.IconBox}
                     style={{
                         background: color.bg,
                         borderColor: color.border,
                         color: color.fg,
                     }}>
                    <NodeIcon node={props.node} triggerInfo={props.triggerInfo}/>
                </div>
                <div className={cls.Titles}>
                    <div className={cls.TypeLabel}>{typeLabel(props.node, props.triggerInfo)}</div>
                    <div className={cls.Title}>{title(props.node, props.triggerInfo)}</div>
                </div>
            </div>

            <div className={cls.Chips}>
                <NodeChips
                    node={props.node}
                    tools={props.tools}
                    triggerInfo={props.triggerInfo}
                    momCandidates={props.momCandidates}/>
            </div>
            <Button
                className={cls.AddButton}
                aria-label="Add block after"
                title="Add block after"
                draggable={false}
                onClick={e => {
                    e.stopPropagation()
                    props.onAddBlock()
                }}
            >
                <PlusIcon/>
            </Button>
        </div>
    )
}
