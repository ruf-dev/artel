import {cn} from "@/app/utils/cn.ts"
import {NODE_HEIGHT, NODE_WIDTH} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import cls from "@/pages/tract-canvas/components/TractCanvasArea/components/ConnectorPath/ConnectorPath.module.css"

interface ConnectorPathProps {
    from: {x: number; y: number}
    to: {x: number; y: number}
    running: boolean
}

export default function ConnectorPath({from, to, running}: ConnectorPathProps) {
    const sx = from.x + NODE_WIDTH
    const sy = from.y + NODE_HEIGHT / 2
    const tx = to.x
    const ty = to.y + NODE_HEIGHT / 2
    const midX = (sx + tx) / 2
    const d = `M ${sx} ${sy} C ${midX} ${sy}, ${midX} ${ty}, ${tx} ${ty}`
    return <path className={cn(cls.Conn, running && cls.ConnRunning)} d={d}/>
}
