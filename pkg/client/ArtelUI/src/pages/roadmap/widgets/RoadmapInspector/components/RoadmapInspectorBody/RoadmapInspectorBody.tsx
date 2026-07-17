import {Button} from "@vervstack/chures"

import {RoadmapNode} from "@/pages/roadmap/processes/roadmapGraph.ts"
import {ParsedTaskLink} from "@/pages/roadmap/processes/taskLinkParser.ts"
import LinkRow from "@/pages/roadmap/widgets/RoadmapInspector/components/LinkRow/LinkRow.tsx"
import cls
    from "@/pages/roadmap/widgets/RoadmapInspector/components/RoadmapInspectorBody/RoadmapInspectorBody.module.css"

interface Props {
    node: RoadmapNode
    nodes: Record<string, RoadmapNode>
    onSelectNode: (id: string) => void
    onAddLink: () => void
    onClose: () => void
}

function resolveLinkTarget(link: ParsedTaskLink, nodes: Record<string, RoadmapNode>): RoadmapNode | undefined {
    if (!link.shortLink) return undefined
    return Object.values(nodes).find(n => n.shortLink === link.shortLink || n.id === link.shortLink)
}

export default function RoadmapInspectorBody({node, nodes, onSelectNode, onAddLink, onClose}: Props) {
    return (
        <div className={cls.RoadmapInspectorBodyContainer}>
            <div className={cls.Head}>
                <div className={cls.Titles}>
                    <div className={cls.Title}>{node.name}</div>
                    <div className={cls.Sub}>{node.closed ? "Closed" : "Open"}</div>
                </div>
                <Button variant="iconDanger" className={cls.CloseBtn} onClick={onClose} aria-label="Close inspector">
                    ✕
                </Button>
            </div>

            <div className={cls.Body}>
                <p className={cls.Desc}>{node.desc || "No description."}</p>

                <a className={cls.OpenLink} href={node.url} target="_blank" rel="noopener noreferrer">
                    Open in Trello ↗
                </a>

                <div className={cls.LinksHead}>
                    <span className={cls.LinksTitle}>Links</span>
                    <Button variant="ghost" onClick={onAddLink}>+ Add link</Button>
                </div>

                {node.links.length === 0 && <p className={cls.Empty}>No links yet.</p>}

                <div className={cls.LinksList}>
                    {node.links.map((link, i) => (
                        <LinkRow
                            key={`${link.relation}-${link.url}-${i}`}
                            link={link}
                            target={resolveLinkTarget(link, nodes)}
                            onSelectNode={onSelectNode}
                        />
                    ))}
                </div>
            </div>
        </div>
    )
}
