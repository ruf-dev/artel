import {useEffect} from "react"
import {useNavigate, useSearchParams} from "react-router-dom"

import cls from "@/pages/connections/ConnectionsPage.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import useUser from "@/hooks/user/User.ts"
import HeroSegment from "@/components/HeroSegment/HeroSegment.tsx"
import Tabs from "@/components/atoms/Tabs/Tabs.tsx"
import ContentSegment from "@/pages/connections/components/ContentSegment/ContentSegment.tsx"
import BYOKSection from "@/pages/connections/components/BYOKSection/BYOKSection.tsx"
import ExternalConnectionsTabIcon
    from "@/pages/connections/components/ExternalConnectionsTabIcon/ExternalConnectionsTabIcon.tsx"
import ByokTabIcon from "@/pages/connections/components/ByokTabIcon/ByokTabIcon.tsx"

export default function ConnectionsPage() {
    const navigate = useNavigate()
    const {auth} = useUser()
    const {connections, loading, fetch: fetchConnections} = useExternalConnections()
    const bakeError = useBakeError()
    const [searchParams, setSearchParams] = useSearchParams()

    const activeTab = (searchParams.get("tab") === "byok" ? "byok" : "external")

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
            setSearchParams(prev => {
                const next = new URLSearchParams(prev)
                next.delete("status")
                return next
            }, {replace: true})
        }

        void fetchConnections()
    }, [auth])

    function handleTabChange(tab: string) {
        if (tab === "byok") {
            setSearchParams({tab: "byok"})
        } else {
            setSearchParams({})
        }
    }

    return (
        <div className={cls.Root}>
            <HeroSegment
                eyebrow="Workspace"
                title="Connected services"
                subtitle={
                    <>
                        <b>
                            {loading
                                ? "…"
                                : `${connections.length} ${connections.length === 1 ? "service" : "services"}`}
                        </b>
                        {" · "}<span>linked external integrations</span>
                    </>
                }
            />
            <Tabs
                tabs={[
                    {key: "external", label: "External Connections", icon: <ExternalConnectionsTabIcon/>},
                    {key: "byok", label: "BYOK", icon: <ByokTabIcon/>},
                ]}
                active={activeTab}
                onChange={handleTabChange}
            />
            {activeTab === "external" ? <ContentSegment/> : <BYOKSection/>}
        </div>
    )
}
