import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import {ParsedTaskLink, RoadmapRelation} from "@/pages/roadmap/processes/taskLinkParser.ts"
import {RoadmapNode} from "@/pages/roadmap/processes/roadmapGraph.ts"
import cls from "@/pages/roadmap/widgets/RoadmapInspector/components/LinkRow/LinkRow.module.css"

const RELATION_LABEL: Record<RoadmapRelation, string> = {
    depends_on: "Depends on",
    blocks: "Blocks",
    blocked_by: "Blocked by",
    related_to: "Related to",
}

interface Props {
    link: ParsedTaskLink
    /** The already-loaded node this link resolves to, if any — undefined means "not in the graph
     * yet" and renders as an external-link chip instead of a jump-to-node row. */
    target?: RoadmapNode
    onSelectNode: (id: string) => void
}

export default function LinkRow({link, target, onSelectNode}: Props) {
    if (target) {
        return (
            <Button variant="ghost" className={cls.LinkRowContainer} onClick={() => onSelectNode(target.id)}>
                <span className={cls.RelationTag}>{RELATION_LABEL[link.relation]}</span>
                <span className={cls.TargetName}>{target.name}</span>
            </Button>
        )
    }

    return (
        <a
            className={cn(cls.LinkRowContainer, cls.ExternalChip)}
            href={link.url}
            target="_blank"
            rel="noopener noreferrer"
        >
            <span className={cls.RelationTag}>{RELATION_LABEL[link.relation]}</span>
            <span className={cls.TargetName}>{link.url}</span>
        </a>
    )
}
