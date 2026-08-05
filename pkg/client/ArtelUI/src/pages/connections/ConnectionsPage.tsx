import {useEffect} from "react"
import {useSearchParams} from "react-router-dom"

import cls from "@/pages/connections/ConnectionsPage.module.css"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import useUser from "@/hooks/user/User.ts"
import HeroSegment from "@/components/HeroSegment/HeroSegment.tsx"
import Tabs from "@/components/atoms/Tabs/Tabs.tsx"
import ContentSegment from "@/pages/connections/components/ContentSegment/ContentSegment.tsx"
import BYOKSection from "@/pages/connections/components/BYOKSection/BYOKSection.tsx"
import CommunitySection from "@/pages/connections/components/CommunitySection/CommunitySection.tsx"
import ExternalConnectionsTabIcon
    from "@/pages/connections/components/ExternalConnectionsTabIcon/ExternalConnectionsTabIcon.tsx"
import ByokTabIcon from "@/pages/connections/components/ByokTabIcon/ByokTabIcon.tsx"
import CommunityTabIcon from "@/pages/connections/components/CommunityTabIcon/CommunityTabIcon.tsx"

type ConnectionsTab = "external" | "byok" | "community"

function resolveTab(value: string | null): ConnectionsTab {
    if (value === "byok" || value === "community") return value
    return "external"
}

export default function ConnectionsPage() {
    const {auth} = useUser()
    const {connections, loading, fetch: fetchConnections} = useExternalConnections()
    const bakeError = useBakeError()
    const [searchParams, setSearchParams] = useSearchParams()

    const activeTab = resolveTab(searchParams.get("tab"))

    // Unauthenticated → /init is handled at the router level (HomeLayout.tsx), not per-page.
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
        if (tab === "byok" || tab === "community") {
            setSearchParams({tab})
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
                    {key: "community", label: "Community", icon: <CommunityTabIcon/>},
                ]}
                active={activeTab}
                onChange={handleTabChange}
            />
            {activeTab === "external" && <ContentSegment/>}
            {activeTab === "byok" && <BYOKSection/>}
            {activeTab === "community" && <CommunitySection/>}
        </div>
    )
}
