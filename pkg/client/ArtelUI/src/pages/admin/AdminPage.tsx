import {useState, useEffect} from "react"
import {useNavigate} from "react-router-dom"

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

export default function AdminPage() {
    const navigate = useNavigate()
    const {auth, isAdmin} = useUser()
    const [tab, setTab] = useState<Tab>("instances")

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
            return
        }
        if (!isAdmin) {
            navigate(Path.HomePage)
        }
    }, [auth, isAdmin, navigate])

    return (
        <div className={cls.AdminPageContainer}>
            <AdminHero tab={tab} />
            <TabBar tab={tab} onTabChange={setTab} />
            {tab === "instances" && <InstancesTab />}
            {tab === "couch_users" && <UsersTab />}
            {tab === "users" && <ArtelUsersTab />}
            {tab === "s3_instances" && (
                <div className={cls.ContentContainer}>
                    <S3InstancesTab />
                </div>
            )}
        </div>
    )
}
