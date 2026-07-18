import {DragEvent, useState} from "react"

import {CanvasNode} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import {DragOverSide} from "@/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx"
import {DropSide} from "@/processes/tractStepsMove.ts"

const STEP_DRAG_MIME = "application/x-artel-tract-step-id"

// Owns pointer/hover state for dragging a step node onto another node's left half ("before") or
// right half ("after") to reorder it — see TractCanvasArea, which renders the indicator this
// exposes per-node and calls onMoveStep on drop.
export function useTractCanvasDrag(onMoveStep: (sourceId: string, targetId: string, side: DropSide) => void) {
    const [draggingId, setDraggingId] = useState<string | null>(null)
    const [dragOverId, setDragOverId] = useState<string | null>(null)
    const [dragOverSide, setDragOverSide] = useState<DragOverSide | null>(null)

    function clearDragState() {
        setDraggingId(null)
        setDragOverId(null)
        setDragOverSide(null)
    }

    function handleNodeDragStart(e: DragEvent<HTMLDivElement>, node: CanvasNode) {
        e.dataTransfer.effectAllowed = "move"
        e.dataTransfer.setData(STEP_DRAG_MIME, node.id)
        setDraggingId(node.id)
    }

    function handleNodeDragOver(e: DragEvent<HTMLDivElement>, node: CanvasNode) {
        if (!draggingId || node.id === draggingId) return
        e.preventDefault()
        e.dataTransfer.dropEffect = "move"
        // The trigger has nothing before it — hovering anywhere over it always means "after".
        const rect = e.currentTarget.getBoundingClientRect()
        const side: DragOverSide = node.kind === "trigger"
            ? "after"
            : (e.clientX < rect.left + rect.width / 2 ? "before" : "after")
        setDragOverId(node.id)
        setDragOverSide(side)
    }

    function handleNodeDragLeave(e: DragEvent<HTMLDivElement>, node: CanvasNode) {
        if (e.currentTarget.contains(e.relatedTarget as Node | null)) return
        if (dragOverId === node.id) {
            setDragOverId(null)
            setDragOverSide(null)
        }
    }

    function handleNodeDrop(e: DragEvent<HTMLDivElement>, node: CanvasNode) {
        e.preventDefault()
        e.stopPropagation()
        const sourceId = e.dataTransfer.getData(STEP_DRAG_MIME)
        const side = dragOverSide
        clearDragState()
        if (!sourceId || !side || sourceId === node.id) return
        onMoveStep(sourceId, node.id, side)
    }

    return {
        draggingId, dragOverId, dragOverSide,
        handleNodeDragStart, handleNodeDragOver, handleNodeDragLeave, handleNodeDrop, clearDragState,
    }
}
