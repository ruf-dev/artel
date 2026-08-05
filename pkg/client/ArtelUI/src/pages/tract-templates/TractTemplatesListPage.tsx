import {useEffect} from "react"

import cls from "@/pages/tract-templates/TractTemplatesListPage.module.css"
import {useTractTemplates} from "@/app/hooks/TractTemplates.ts"
import useUser from "@/hooks/user/User.ts"
import HeroSegment from "@/components/HeroSegment/HeroSegment.tsx"
import ContentSegment from "@/pages/tract-templates/segments/ContentSegment/ContentSegment.tsx"

export default function TractTemplatesListPage() {
    const {auth} = useUser()
    const {templates, loading, mineOnly, fetch} = useTractTemplates()

    // Unauthenticated → /init is handled at the router level (HomeLayout.tsx), not per-page.
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
