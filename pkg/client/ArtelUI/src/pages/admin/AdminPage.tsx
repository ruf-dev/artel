import {useEffect} from "react"
import {useNavigate, useSearchParams} from "react-router-dom"

import cls from "@/pages/admin/AdminPage.module.css"
import {Path} from "@/app/routing/Router.tsx"
import useUser from "@/hooks/user/User.ts"
import S3InstancesTab from "@/widgets/S3InstancesTab/S3InstancesTab.tsx"
import {Tab} from "@/pages/admin/adminTypes.ts"
import AdminHero from "@/pages/admin/components/AdminHero/AdminHero.tsx"
import TabBar from "@/pages/admin/components/TabBar/TabBar.tsx"
import InstancesTab from "@/pages/admin/components/InstancesTab/InstancesTab.tsx"
import UsersTab from "@/pages/admin/components/UsersTab/UsersTab.tsx"
import ArtelUsersTab from "@/pages/admin/components/ArtelUsersTab/ArtelUsersTab.tsx"
import DockerApiTab from "@/pages/admin/components/DockerApiTab/DockerApiTab.tsx"

function resolveTab(value: string | null): Tab {
    if (value === "couch_users" || value === "users" || value === "s3_instances" || value === "docker_api") return value
    return "instances"
}

export default function AdminPage() {
    const navigate = useNavigate()
    const {auth, isAdmin} = useUser()
    const [searchParams, setSearchParams] = useSearchParams()
    const tab = resolveTab(searchParams.get("tab"))

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
            return
        }
        if (!isAdmin) {
            navigate(Path.HomePage)
        }
    }, [auth, isAdmin, navigate])

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
        if (t === "instances") {
            setSearchParams({})
        } else {
            setSearchParams({tab: t})
        }
    }

    const autoOpenDockerDialog = tab === "docker_api" && searchParams.get("dialog") === "true"

    return (
        <div className={cls.AdminPageContainer}>
            <AdminHero tab={tab} />
            <TabBar tab={tab} onTabChange={handleTabChange} />
            {tab === "instances" && <InstancesTab />}
            {tab === "couch_users" && <UsersTab />}
            {tab === "users" && <ArtelUsersTab />}
            {tab === "s3_instances" && (
                <div className={cls.ContentContainer}>
                    <S3InstancesTab />
                </div>
            )}
            {tab === "docker_api" && <DockerApiTab autoOpenAddDialog={autoOpenDockerDialog} />}
        </div>
    )
}
