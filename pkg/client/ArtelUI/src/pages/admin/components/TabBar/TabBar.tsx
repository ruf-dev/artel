import Tabs from "@/components/atoms/Tabs/Tabs.tsx"
import cls from "@/pages/admin/components/TabBar/TabBar.module.css"
import {Tab} from "@/pages/admin/adminTypes.ts"
import CouchDbTabIcon from "@/pages/admin/components/icons/CouchDbTabIcon.tsx"
import ArtelUsersTabIcon from "@/pages/admin/components/icons/ArtelUsersTabIcon.tsx"
import S3TabIcon from "@/pages/admin/components/icons/S3TabIcon.tsx"
import DockerApiTabIcon from "@/pages/admin/components/icons/DockerApiTabIcon.tsx"
import SettingsTabIcon from "@/pages/admin/components/icons/SettingsTabIcon.tsx"

interface TabBarProps {
    tab: Tab
    onTabChange: (t: Tab) => void
}

export default function TabBar({tab, onTabChange}: TabBarProps) {
    return (
        <div className={cls.TabBarContainer}>
            <Tabs
                tabs={[
                    {key: "couchdb", label: "CouchDB", icon: <CouchDbTabIcon/>},
                    {key: "users", label: "Users", icon: <ArtelUsersTabIcon/>},
                    {key: "s3_instances", label: "S3 storage", icon: <S3TabIcon/>},
                    {key: "docker_api", label: "Docker API", icon: <DockerApiTabIcon/>},
                    {key: "settings", label: "Settings", icon: <SettingsTabIcon/>},
                ]}
                active={tab}
                onChange={(key) => onTabChange(key as Tab)}
            />
        </div>
    )
}
