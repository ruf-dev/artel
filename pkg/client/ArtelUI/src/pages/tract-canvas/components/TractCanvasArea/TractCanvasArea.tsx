import React, {useRef, useState} from "react"

import cls from "@/pages/tract-canvas/components/TractCanvasArea/TractCanvasArea.module.css"
import {cn} from "@/app/utils/cn.ts"
import {CanvasLayout} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import {Location} from "@/processes/tractSteps.ts"
import {DropSide} from "@/processes/tractStepsMove.ts"
import {TractTool} from "@/processes/Tracts.ts"
import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import TractCanvasNode, {NodeStatus} from "@/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx"
import ConnectorPath from "@/pages/tract-canvas/components/TractCanvasArea/components/ConnectorPath/ConnectorPath.tsx"
import ParallelBoxes from "@/pages/tract-canvas/components/TractCanvasArea/components/ParallelBoxes/ParallelBoxes.tsx"
import {useTractCanvasDrag} from "@/pages/tract-canvas/components/TractCanvasArea/useTractCanvasDrag.ts"

const DRAG_THRESHOLD_PX = 4

interface Props {
    layout: CanvasLayout
    tools: TractTool[]
    triggerInfo?: {name: string; kind: string; source: string}
    momCandidates: MomCandidate[]
    selectedNodeId: string | null
    onSelectNode: (id: string) => void
    onBackgroundClick: () => void
    nodeStatus: (id: string) => NodeStatus
    runningEdgeIds: Set<string>
    onAddBlock: (location: Location, index: number) => void
    onMoveStep: (sourceId: string, targetId: string, side: DropSide) => void
}

export default function TractCanvasArea(props: Props) {
    const {layout, tools, triggerInfo, momCandidates, selectedNodeId} = props
    const {onSelectNode, onBackgroundClick, nodeStatus, runningEdgeIds, onAddBlock, onMoveStep} = props
    const wrapRef = useRef<HTMLDivElement>(null)
    const dragRef = useRef<
        {startX: number; startY: number; scrollLeft: number; scrollTop: number; moved: boolean} | null
    >(null)
    const suppressClickRef = useRef(false)
    const [panning, setPanning] = useState(false)
    const boxedNodeIds = new Set(layout.parallelBoxes.map(b => b.id))
    const drag = useTractCanvasDrag(onMoveStep)

    function handlePointerDown(e: React.PointerEvent<HTMLDivElement>) {
        if (e.button !== 0) return
        if ((e.target as HTMLElement).closest("[data-tract-node]")) return
        const wrap = wrapRef.current
        if (!wrap) return
        dragRef.current = {
            startX: e.clientX, startY: e.clientY, scrollLeft: wrap.scrollLeft, scrollTop: wrap.scrollTop, moved: false,
        }
        wrap.setPointerCapture(e.pointerId)
    }

    function handlePointerMove(e: React.PointerEvent<HTMLDivElement>) {
        const drag = dragRef.current
        const wrap = wrapRef.current
        if (!drag || !wrap) return
        const dx = e.clientX - drag.startX
        const dy = e.clientY - drag.startY
        if (!drag.moved && Math.hypot(dx, dy) > DRAG_THRESHOLD_PX) {
            drag.moved = true
            setPanning(true)
        }
        if (drag.moved) {
            wrap.scrollLeft = drag.scrollLeft - dx
            wrap.scrollTop = drag.scrollTop - dy
        }
    }

    function handlePointerUp(e: React.PointerEvent<HTMLDivElement>) {
        const wrap = wrapRef.current
        if (wrap) wrap.releasePointerCapture(e.pointerId)
        if (dragRef.current?.moved) suppressClickRef.current = true
        dragRef.current = null
        setPanning(false)
    }

    function handleClickCapture(e: React.MouseEvent<HTMLDivElement>) {
        if (suppressClickRef.current) {
            suppressClickRef.current = false
            e.stopPropagation()
        }
    }

    return (
        <div
            ref={wrapRef}
            className={cn(cls.Wrap, panning && cls.Panning)}
            onClick={onBackgroundClick}
            onClickCapture={handleClickCapture}
            onPointerDown={handlePointerDown}
            onPointerMove={handlePointerMove}
            onPointerUp={handlePointerUp}
        >
            <div className={cls.Canvas} style={{width: layout.width, height: layout.height}}>
                <svg className={cls.Svg} width={layout.width} height={layout.height}>
                    {layout.edges.map(edge => {
                        const from = layout.nodes.find(n => n.id === edge.fromId)
                        const to = layout.nodes.find(n => n.id === edge.toId)
                        if (!from || !to) return null
                        return <ConnectorPath key={edge.id} from={from} to={to} running={runningEdgeIds.has(edge.id)}/>
                    })}
                </svg>
                <ParallelBoxes
                    boxes={layout.parallelBoxes} selectedNodeId={selectedNodeId} onSelectNode={onSelectNode}/>
                {layout.nodes.filter(node => !boxedNodeIds.has(node.id)).map(node => (
                    <TractCanvasNode
                        key={node.id}
                        node={node}
                        tools={tools}
                        triggerInfo={triggerInfo}
                        momCandidates={momCandidates}
                        status={nodeStatus(node.id)}
                        selected={node.id === selectedNodeId}
                        dragging={node.id === drag.draggingId}
                        dragOverSide={node.id === drag.dragOverId ? drag.dragOverSide : null}
                        onClick={() => onSelectNode(node.id)}
                        onAddBlock={() => onAddBlock(node.nextLocation, node.nextIndex)}
                        onDragStart={e => drag.handleNodeDragStart(e, node)}
                        onDragOver={e => drag.handleNodeDragOver(e, node)}
                        onDragLeave={e => drag.handleNodeDragLeave(e, node)}
                        onDrop={e => drag.handleNodeDrop(e, node)}
                        onDragEnd={drag.clearDragState}
                    />
                ))}
            </div>
        </div>
    )
}
