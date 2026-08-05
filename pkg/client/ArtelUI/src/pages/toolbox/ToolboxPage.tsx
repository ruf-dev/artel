import {useEffect} from "react"

import cls from "@/pages/toolbox/ToolboxPage.module.css"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import useUser from "@/hooks/user/User.ts"
import HeroSegment from "@/components/HeroSegment/HeroSegment.tsx"
import ContentSegment from "@/pages/toolbox/components/ContentSegment/ContentSegment.tsx"

export default function ToolboxPage() {
    const {auth} = useUser()
    const {momCandidates, momCandidatesLoading, fetchMomCandidates} = useMcpKeys()

    // Unauthenticated → /init is handled at the router level (HomeLayout.tsx), not per-page.
    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetchMomCandidates()
        }
    }, [auth, fetchMomCandidates])

    return (
        <div className={cls.Root}>
            <HeroSegment
                eyebrow="MCP"
                title="Toolbox"
                subtitle={
                    <>
                        <b>
                            {momCandidatesLoading
                                ? "…"
                                : `${momCandidates.length} ${momCandidates.length === 1 ? "tool" : "tools"}`}
                        </b>
                        {" · "}<span>available in this installation</span>
                    </>
                }
            />
            <ContentSegment/>
        </div>
    )
}
