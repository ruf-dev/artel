import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import {NODE_HEIGHT, NODE_WIDTH, RoadmapLayoutNode} from "@/pages/roadmap/processes/roadmapLayout.ts"
import cls from "@/pages/roadmap/components/RoadmapCanvasNode/RoadmapCanvasNode.module.css"

interface Props {
    layoutNode: RoadmapLayoutNode
    isRoot: boolean
    selected: boolean
    onClick: () => void
    onExpand: () => void
}

// No board/list *names* are available on a RoadmapNode — get_card only returns idBoard/idList,
// and resolving those to human names would need an extra list_lists call per board. Kept simple
// per the plan's note that a follow-up design pass will refine the visual treatment.
function boardListLabel(idBoard: string, idList: string): string {
    return `board ${idBoard.slice(-6)} · list ${idList.slice(-6)}`
}

export default function RoadmapCanvasNode({layoutNode, isRoot, selected, onClick, onExpand}: Props) {
    const {node} = layoutNode

    return (
        <div
            className={cn(cls.RoadmapCanvasNodeContainer, selected && cls.Selected, isRoot && cls.RootNode)}
            style={{left: layoutNode.x, top: layoutNode.y, width: NODE_WIDTH, height: NODE_HEIGHT}}
            data-roadmap-node
            onClick={e => {
                e.stopPropagation()
                onClick()
            }}
            role="button"
            tabIndex={0}
        >
            <span className={cn(cls.StatusDot, node.closed ? cls.StatusClosed : cls.StatusOpen)}/>

            <div className={cls.Head}>
                <div className={cls.Title}>{node.name}</div>
                <div className={cls.SubLabel}>{boardListLabel(node.idBoard, node.idList)}</div>
            </div>

            {node.links.length > 0 && (
                <span className={cls.LinkBadge}>{node.links.length} link{node.links.length === 1 ? "" : "s"}</span>
            )}

            {!node.expanded && (
                <Button
                    variant="ghost"
                    className={cls.ExpandBtn}
                    onClick={e => {
                        e.stopPropagation()
                        onExpand()
                    }}
                >
                    Expand
                </Button>
            )}
        </div>
    )
}
