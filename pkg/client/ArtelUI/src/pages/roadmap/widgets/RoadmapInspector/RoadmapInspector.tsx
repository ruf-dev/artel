import {cn} from "@/app/utils/cn.ts"
import {RoadmapNode} from "@/pages/roadmap/processes/roadmapGraph.ts"
import RoadmapInspectorBody
    from "@/pages/roadmap/widgets/RoadmapInspector/components/RoadmapInspectorBody/RoadmapInspectorBody.tsx"
import cls from "@/pages/roadmap/widgets/RoadmapInspector/RoadmapInspector.module.css"

interface Props {
    node: RoadmapNode | null
    nodes: Record<string, RoadmapNode>
    onSelectNode: (id: string) => void
    onAddLink: () => void
    onClose: () => void
}

// Colocated under pages/roadmap/widgets rather than the global src/widgets — it owns no data
// fetching of its own (the graph is fetched/held by RoadmapPage and passed down) and has exactly
// one consumer, so promoting it to the global widgets/ tier would add indirection with no reuse
// payoff. See pkg/client/ArtelUI/CLAUDE.md's widget-vs-colocation guidance.
export default function RoadmapInspector({node, nodes, onSelectNode, onAddLink, onClose}: Props) {
    return (
        <div className={cls.RoadmapInspectorContainer}>
            <div className={cn(cls.Panel, node && cls.Open)}>
                {node && (
                    <RoadmapInspectorBody
                        key={node.id}
                        node={node}
                        nodes={nodes}
                        onSelectNode={onSelectNode}
                        onAddLink={onAddLink}
                        onClose={onClose}
                    />
                )}
            </div>
        </div>
    )
}
