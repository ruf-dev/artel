import {cn} from "@/app/utils/cn.ts"
import {NODE_HEIGHT, NODE_WIDTH} from "@/pages/roadmap/processes/roadmapLayout.ts"
import {RoadmapRelation} from "@/pages/roadmap/processes/taskLinkParser.ts"
import cls
    from "@/pages/roadmap/components/RoadmapCanvasArea/components/RoadmapConnectorPath/RoadmapConnectorPath.module.css"

const RELATION_CLASS: Record<RoadmapRelation, string> = {
    depends_on: cls.DependsOn,
    blocks: cls.Blocks,
    blocked_by: cls.BlockedBy,
    related_to: cls.RelatedTo,
}

interface Props {
    from: {x: number; y: number}
    to: {x: number; y: number}
    relation: RoadmapRelation
}

export default function RoadmapConnectorPath({from, to, relation}: Props) {
    const sx = from.x + NODE_WIDTH
    const sy = from.y + NODE_HEIGHT / 2
    const tx = to.x
    const ty = to.y + NODE_HEIGHT / 2
    const midX = (sx + tx) / 2
    const d = `M ${sx} ${sy} C ${midX} ${sy}, ${midX} ${ty}, ${tx} ${ty}`
    return <path className={cn(cls.Conn, RELATION_CLASS[relation])} d={d}/>
}
