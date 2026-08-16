import {useEffect} from "react"
import {useSearchParams} from "react-router-dom"

import cls from "@/pages/admin/AdminPage.module.css"
import S3InstancesTab from "@/widgets/S3InstancesTab/S3InstancesTab.tsx"
import {Tab, CouchSubTab} from "@/pages/admin/adminTypes.ts"
import AdminHero from "@/pages/admin/components/AdminHero/AdminHero.tsx"
import TabBar from "@/pages/admin/components/TabBar/TabBar.tsx"
import CouchDbTab from "@/pages/admin/components/CouchDbTab/CouchDbTab.tsx"
import ArtelUsersTab from "@/pages/admin/components/ArtelUsersTab/ArtelUsersTab.tsx"
import DockerApiTab from "@/pages/admin/components/DockerApiTab/DockerApiTab.tsx"
import SettingsTab from "@/pages/admin/components/SettingsTab/SettingsTab.tsx"

function resolveTab(value: string | null): Tab {
    if (
        value === "users" ||
        value === "s3_instances" ||
        value === "docker_api" ||
        value === "settings"
    ) {
        return value
    }
    return "couchdb"
}

function resolveSubTab(value: string | null): CouchSubTab {
    if (value === "users") {
        return "users"
    }
    return "instances"
}

export default function AdminPage() {
    const [searchParams, setSearchParams] = useSearchParams()
    const tab = resolveTab(searchParams.get("tab"))
    const subTab = resolveSubTab(searchParams.get("subtab"))

    useEffect(() => {
        if (searchParams.get("dialog") === "true") {
            setSearchParams(prev => {
                const next = new URLSearchParams(prev)
                next.delete("dialog")
                return next
            }, {replace: true})
        }
    }, [searchParams, setSearchParams])

    function handleTabChange(t: Tab) {
        if (t === "couchdb") {
            setSearchParams({})
        } else {
            setSearchParams({tab: t})
        }
    }

    function handleSubTabChange(st: CouchSubTab) {
        setSearchParams(prev => {
            const next = new URLSearchParams(prev)
            if (st === "instances") {
                next.delete("subtab")
            } else {
                next.set("subtab", st)
            }
            return next
        })
    }

    const autoOpenDockerDialog = tab === "docker_api" && searchParams.get("dialog") === "true"

    return (
        <div className={cls.AdminPageContainer}>
            <AdminHero tab={tab} subTab={subTab}/>
            <div className={cls.StickyTabsWrapper}>
                {tab === "couchdb" && <CouchDbTab subTab={subTab} onSubTabChange={handleSubTabChange}/>}
                {tab === "users" && <ArtelUsersTab/>}
                {tab === "s3_instances" && (
                    <div className={cls.ContentContainer}>
                        <S3InstancesTab/>
                    </div>
                )}
                {tab === "docker_api" && <DockerApiTab autoOpenAddDialog={autoOpenDockerDialog}/>}
                {tab === "settings" && <SettingsTab/>}
                <TabBar tab={tab} onTabChange={handleTabChange}/>
            </div>
        </div>
    )
}
