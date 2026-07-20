import {useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/tract-templates/TractTemplatesListPage.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {useTractTemplates} from "@/app/hooks/TractTemplates.ts"
import useUser from "@/hooks/user/User.ts"
import HeroSegment from "@/components/HeroSegment/HeroSegment.tsx"
import ContentSegment from "@/pages/tract-templates/segments/ContentSegment/ContentSegment.tsx"

export default function TractTemplatesListPage() {
    const navigate = useNavigate()
    const {auth} = useUser()
    const {templates, loading, mineOnly, fetch} = useTractTemplates()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetch()
        }
    }, [auth, fetch])

    return (
        <div className={cls.Root}>
            <HeroSegment
                eyebrow="Workspace · Automations"
                title="Templates"
                subtitle={
                    <>
                        <b>
                            {loading
                                ? "…"
                                : `${templates.length} ${templates.length === 1 ? "template" : "templates"}`}
                        </b>
                        {" · "}
                        <span>{mineOnly ? "yours" : "public gallery"}</span>
                    </>
                }
                compact
            />
            <ContentSegment/>
        </div>
    )
}
