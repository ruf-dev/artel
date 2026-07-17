import React, {useRef, useState} from "react"

import {cn} from "@/app/utils/cn.ts"
import {RoadmapLayout} from "@/pages/roadmap/processes/roadmapLayout.ts"
import RoadmapCanvasNode from "@/pages/roadmap/components/RoadmapCanvasNode/RoadmapCanvasNode.tsx"
import RoadmapConnectorPath
    from "@/pages/roadmap/components/RoadmapCanvasArea/components/RoadmapConnectorPath/RoadmapConnectorPath.tsx"
import cls from "@/pages/roadmap/components/RoadmapCanvasArea/RoadmapCanvasArea.module.css"

const DRAG_THRESHOLD_PX = 4

interface Props {
    layout: RoadmapLayout
    rootId: string
    selectedNodeId: string | null
    onSelectNode: (id: string) => void
    onBackgroundClick: () => void
    onExpandNode: (id: string) => void
}

export default function RoadmapCanvasArea(props: Props) {
    const {layout, rootId, selectedNodeId, onSelectNode, onBackgroundClick, onExpandNode} = props
    const wrapRef = useRef<HTMLDivElement>(null)
    const dragRef = useRef<
        {startX: number; startY: number; scrollLeft: number; scrollTop: number; moved: boolean} | null
    >(null)
    const suppressClickRef = useRef(false)
    const [panning, setPanning] = useState(false)

    function handlePointerDown(e: React.PointerEvent<HTMLDivElement>) {
        if (e.button !== 0) return
        if ((e.target as HTMLElement).closest("[data-roadmap-node]")) return
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
            className={cn(cls.RoadmapCanvasAreaContainer, panning && cls.Panning)}
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
                        return <RoadmapConnectorPath key={edge.id} from={from} to={to} relation={edge.relation}/>
                    })}
                </svg>
                {layout.nodes.map(layoutNode => (
                    <RoadmapCanvasNode
                        key={layoutNode.id}
                        layoutNode={layoutNode}
                        isRoot={layoutNode.id === rootId}
                        selected={layoutNode.id === selectedNodeId}
                        onClick={() => onSelectNode(layoutNode.id)}
                        onExpand={() => onExpandNode(layoutNode.id)}
                    />
                ))}
            </div>
        </div>
    )
}
