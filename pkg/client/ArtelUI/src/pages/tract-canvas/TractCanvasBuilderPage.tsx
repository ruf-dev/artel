import {useEffect} from "react"
import {useNavigate, useParams} from "react-router-dom"

import cls from "@/pages/tract-canvas/TractCanvasBuilderPage.module.css"

import {Path} from "@/app/routing/Router.tsx"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import useUser from "@/hooks/user/User.ts"

import {Button} from "@vervstack/chures"
import TractCanvasBuilder from "@/pages/tract-canvas/components/TractCanvasBuilder/TractCanvasBuilder.tsx"

export default function TractCanvasBuilderPage() {
    const {tractUuid} = useParams<{ tractUuid: string }>()
    const navigate = useNavigate()
    const {auth} = useUser()
    const {tracts, fetch, fetchRuns, runsByTract, fetchTools, tools, triggers, fetchTriggers} = useTracts()
    const {fetch: fetchConnections} = useExternalConnections()

    useEffect(() => {
        if (!auth.isAuthenticated()) navigate(Path.InitPage)
    }, [auth, navigate])

    useEffect(() => {
        if (!auth.isAuthenticated() || !tractUuid) return
        if (tracts.length === 0) void fetch()
        void fetchRuns(tractUuid)
        void fetchTools()
        void fetchTriggers()
        void fetchConnections()
    }, [auth, tractUuid])

    const tract = tracts.find(t => t.uuid === tractUuid)
    const runs = tractUuid ? runsByTract[tractUuid] ?? [] : []

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
            tools={tools}
            triggers={triggers}
            runs={runs}
            onBack={() => navigate(Path.TractsPage)}
        />
    )
}
