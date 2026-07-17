import {useRef} from "react"

import cls from "@/pages/tract-canvas/components/TractCanvasLogPanel/components/ResizeHandle/ResizeHandle.module.css"

interface ResizeHandleProps {
    onResize: (newHeight: number) => void
    onDragStart: () => void
    onDragEnd: () => void
}

export default function ResizeHandle({onResize, onDragStart, onDragEnd}: ResizeHandleProps) {
    const draggingRef = useRef(false)

    function handleMouseMove(e: MouseEvent) {
        if (!draggingRef.current) return
        onResize(window.innerHeight - e.clientY)
    }

    function handleMouseUp() {
        draggingRef.current = false
        onDragEnd()
        window.removeEventListener("mousemove", handleMouseMove)
        window.removeEventListener("mouseup", handleMouseUp)
    }

    function handleMouseDown() {
        draggingRef.current = true
        onDragStart()
        window.addEventListener("mousemove", handleMouseMove)
        window.addEventListener("mouseup", handleMouseUp)
    }

    return (
        <div className={cls.ResizeHandleContainer} onMouseDown={handleMouseDown}/>
    )
}
