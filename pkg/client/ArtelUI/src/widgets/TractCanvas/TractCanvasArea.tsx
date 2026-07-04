import cls from "@/widgets/TractCanvas/TractCanvasArea.module.css"

import {CanvasLayout, NODE_HEIGHT, NODE_WIDTH} from "@/processes/tractCanvasLayout.ts"
import {TractTool} from "@/processes/Tracts.ts"

import TractCanvasNode, {NodeStatus} from "@/widgets/TractCanvas/TractCanvasNode.tsx"

interface Props {
    layout: CanvasLayout
    tools: TractTool[]
    triggerInfo?: {name: string; kind: string; source: string}
    selectedNodeId: string | null
    onSelectNode: (id: string) => void
    onBackgroundClick: () => void
    nodeStatus: (id: string) => NodeStatus
    runningEdgeIds: Set<string>
}

export default function TractCanvasArea({layout, tools, triggerInfo, selectedNodeId, onSelectNode, onBackgroundClick, nodeStatus, runningEdgeIds}: Props) {
    return (
        <div className={cls.Wrap} onClick={onBackgroundClick}>
            <div className={cls.Canvas} style={{width: layout.width, height: layout.height}}>
                <svg className={cls.Svg} width={layout.width} height={layout.height}>
                    {layout.edges.map(edge => {
                        const from = layout.nodes.find(n => n.id === edge.fromId)
                        const to = layout.nodes.find(n => n.id === edge.toId)
                        if (!from || !to) return null
                        return <ConnectorPath key={edge.id} from={from} to={to} running={runningEdgeIds.has(edge.id)}/>
                    })}
                </svg>
                {layout.nodes.map(node => (
                    <TractCanvasNode
                        key={node.id}
                        node={node}
                        tools={tools}
                        triggerInfo={triggerInfo}
                        status={nodeStatus(node.id)}
                        selected={node.id === selectedNodeId}
                        onClick={() => onSelectNode(node.id)}
                    />
                ))}
            </div>
        </div>
    )
}

function ConnectorPath({from, to, running}: {
    from: {x: number; y: number}
    to: {x: number; y: number}
    running: boolean
}) {
    const sx = from.x + NODE_WIDTH
    const sy = from.y + NODE_HEIGHT / 2
    const tx = to.x
    const ty = to.y + NODE_HEIGHT / 2
    const midX = (sx + tx) / 2
    const d = `M ${sx} ${sy} C ${midX} ${sy}, ${midX} ${ty}, ${tx} ${ty}`
    return <path className={`${cls.Conn} ${running ? cls.ConnRunning : ""}`} d={d}/>
}
