import {useEffect} from "react"
import {useNavigate, useSearchParams} from "react-router-dom"

import cls from "@/pages/connections/ConnectionsPage.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import useUser from "@/hooks/user/User.ts"
import HeroSegment from "@/pages/connections/components/HeroSegment/HeroSegment.tsx"
import ContentSegment from "@/pages/connections/components/ContentSegment/ContentSegment.tsx"

export default function ConnectionsPage() {
    const navigate = useNavigate()
    const {auth} = useUser()
    const {fetch: fetchConnections} = useExternalConnections()
    const bakeError = useBakeError()
    const [searchParams, setSearchParams] = useSearchParams()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    useEffect(() => {
        if (!auth.isAuthenticated()) return

        const status = searchParams.get("status")
        if (status === "error") {
            bakeError("Connection failed", new Error("Google OAuth authorization was denied or failed."))
        }
        if (status) {
            setSearchParams({}, {replace: true})
        }

        void fetchConnections()
    }, [auth])

    return (
        <div className={cls.Root}>
            <HeroSegment/>
            <ContentSegment/>
        </div>
    )
}
