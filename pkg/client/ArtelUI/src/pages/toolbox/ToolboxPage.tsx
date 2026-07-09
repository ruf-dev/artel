import {useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/toolbox/ToolboxPage.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import useUser from "@/hooks/user/User.ts"
import HeroSegment from "@/pages/toolbox/components/HeroSegment/HeroSegment.tsx"
import ContentSegment from "@/pages/toolbox/components/ContentSegment/ContentSegment.tsx"

export default function ToolboxPage() {
    const navigate = useNavigate()
    const {auth} = useUser()
    const {fetchMomCandidates} = useMcpKeys()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetchMomCandidates()
        }
    }, [auth, fetchMomCandidates])

    return (
        <div className={cls.Root}>
            <HeroSegment/>
            <ContentSegment/>
        </div>
    )
}
