import {useEffect, useMemo, useState} from "react"
import {useNavigate, useParams} from "react-router-dom"
import {Button} from "@vervstack/chures"

import cls from "@/pages/roadmap/RoadmapPage.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import useUser from "@/hooks/user/User.ts"
import {createRoadmapGraph, expandNode, RoadmapGraph} from "@/pages/roadmap/processes/roadmapGraph.ts"
import {layoutRoadmap} from "@/pages/roadmap/processes/roadmapLayout.ts"
import RoadmapCanvasArea from "@/pages/roadmap/components/RoadmapCanvasArea/RoadmapCanvasArea.tsx"
import RoadmapInspector from "@/pages/roadmap/widgets/RoadmapInspector/RoadmapInspector.tsx"
import AddTaskLinkDialog, {RoadmapLinkTarget} from "@/dialogs/AddTaskLinkDialog/AddTaskLinkDialog.tsx"

// Graph state lives here as useState rather than a global app/hooks store — see the comment atop
// pages/roadmap/processes/roadmapGraph.ts for why (colocation rule: a global hook reaching into a
// page-local processes/ folder from outside pages/roadmap would violate CLAUDE.md's colocation
// scope). This page owns the graph and everything below it (canvas, inspector, link dialog)
// receives it — or a narrowed slice of it — as props.
export default function RoadmapPage() {
    const {trackerId, cardId} = useParams<{trackerId: string; cardId: string}>()
    const navigate = useNavigate()
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const bakeError = useBakeError()
    const {executeMomTool} = useMcpKeys()

    const [graph, setGraph] = useState<RoadmapGraph | null>(null)
    const [loading, setLoading] = useState(true)
    const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null)

    // Unauthenticated → /init is handled at the router level (HomeLayout.tsx), not per-page.

    function callTrello(toolName: string, params: Record<string, unknown>): Promise<string> {
        if (!trackerId) return Promise.reject(new Error("no tracker selected"))
        return executeMomTool("trello", toolName, trackerId, params)
    }

    useEffect(() => {
        if (!auth.isAuthenticated() || !trackerId || !cardId) return
        const freshGraph = createRoadmapGraph(cardId)

        setGraph(freshGraph)
        setSelectedNodeId(cardId)
        setLoading(true)
        expandNode(freshGraph, cardId, callTrello)
            .then(setGraph)
            .catch(err => bakeError("Failed to load roadmap", err))
            .finally(() => setLoading(false))
    }, [auth, trackerId, cardId])

    function handleExpandNode(id: string) {
        if (!trackerId || !graph) return
        setLoading(true)
        expandNode(graph, id, callTrello)
            .then(setGraph)
            .catch(err => bakeError("Failed to load roadmap card", err))
            .finally(() => setLoading(false))
    }

    function handleOpenAddLink() {
        if (!trackerId || !graph || !selectedNodeId) return
        const targets: RoadmapLinkTarget[] = Object.values(graph.nodes).map(n => ({id: n.id, name: n.name, url: n.url}))
        OpenDialog(
            <AddTaskLinkDialog
                trackerId={trackerId}
                sourceCardId={selectedNodeId}
                targets={targets}
                onLinked={() => handleExpandNode(selectedNodeId)}
            />,
        )
    }

    const layout = useMemo(() => graph ? layoutRoadmap(graph) : null, [graph])
    const selectedNode = graph && selectedNodeId ? graph.nodes[selectedNodeId] ?? null : null

    if (!trackerId || !cardId) return null

    return (
        <div className={cls.RoadmapPageContainer}>
            <div className={cls.Bar}>
                <Button variant="ghost" onClick={() => navigate(Path.TaskTrackersPage)}>← Task trackers</Button>
                <span className={cls.BarDivider}/>
                <span className={cls.Name}>Roadmap</span>
            </div>
            <div className={cls.Main}>
                {graph && layout ? (
                    <RoadmapCanvasArea
                        layout={layout}
                        rootId={graph.rootId}
                        selectedNodeId={selectedNodeId}
                        onSelectNode={setSelectedNodeId}
                        onBackgroundClick={() => setSelectedNodeId(null)}
                        onExpandNode={handleExpandNode}
                    />
                ) : (
                    <p className={cls.LoadingState}>Loading…</p>
                )}
                <RoadmapInspector
                    node={selectedNode}
                    nodes={graph?.nodes ?? {}}
                    onSelectNode={setSelectedNodeId}
                    onAddLink={handleOpenAddLink}
                    onClose={() => setSelectedNodeId(null)}
                />
            </div>
            {loading && <p className={cls.LoadingBadge}>Loading…</p>}
        </div>
    )
}
