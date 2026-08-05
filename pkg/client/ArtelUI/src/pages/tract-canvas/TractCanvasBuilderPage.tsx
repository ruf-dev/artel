import {useEffect} from "react"
import {useNavigate, useParams} from "react-router-dom"
import {Button} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/TractCanvasBuilderPage.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import useUser from "@/hooks/user/User.ts"
import TractCanvasBuilder from "@/pages/tract-canvas/components/TractCanvasBuilder/TractCanvasBuilder.tsx"

export default function TractCanvasBuilderPage() {
    const {tractUuid} = useParams<{ tractUuid: string }>()
    const navigate = useNavigate()
    const {auth} = useUser()
    const tractsStore = useTracts()
    const {fetch: fetchConnections} = useExternalConnections()
    const {momCandidates, fetchMomCandidates} = useMcpKeys()

    // Unauthenticated → /init is handled at the router level (HomeLayout.tsx), not per-page.
    useEffect(() => {
        if (!auth.isAuthenticated() || !tractUuid) return
        if (tractsStore.tracts.length === 0) void tractsStore.fetch()
        void tractsStore.fetchRuns(tractUuid)
        void tractsStore.fetchTools()
        void tractsStore.fetchTriggers()
        void fetchConnections()
        void fetchMomCandidates()
    }, [auth, tractUuid])

    const tract = tractsStore.tracts.find(t => t.uuid === tractUuid)
    const runs = tractUuid ? tractsStore.runsByTract[tractUuid] ?? [] : []

    if (!tract) {
        return (
            <div className={cls.Root}>
                <div className={cls.Bar}>
                    <Button variant="ghost" onClick={() => navigate(Path.TractsPage)}>← Tracts</Button>
                    <span className={cls.Divider}/>
                    <span className={cls.Name}>…</span>
                </div>
            </div>
        )
    }

    return (
        <TractCanvasBuilder
            tract={tract}
            tools={tractsStore.tools}
            triggers={tractsStore.triggers}
            runs={runs}
            momCandidates={momCandidates}
            onBack={() => navigate(Path.TractsPage)}
        />
    )
}
