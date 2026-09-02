import Tabs from "@/components/atoms/Tabs/Tabs.tsx"
import cls from "@/pages/admin/components/CouchDbTab/CouchDbTab.module.css"
import {CouchSubTab} from "@/pages/admin/adminTypes.ts"
import InstancesTab from "@/pages/admin/components/InstancesTab/InstancesTab.tsx"
import UsersTab from "@/pages/admin/components/UsersTab/UsersTab.tsx"
import InstancesSubTabIcon from "@/pages/admin/components/icons/InstancesSubTabIcon.tsx"
import CouchUsersSubTabIcon from "@/pages/admin/components/icons/CouchUsersSubTabIcon.tsx"

interface CouchDbTabProps {
    subTab: CouchSubTab
    onSubTabChange: (t: CouchSubTab) => void
    autoOpenInstanceDialog?: boolean
}

export default function CouchDbTab({subTab, onSubTabChange, autoOpenInstanceDialog}: CouchDbTabProps) {
    return (
        <div className={cls.CouchDbTabContainer}>
            <Tabs
                tabs={[
                    {key: "instances", label: "Instances", icon: <InstancesSubTabIcon/>},
                    {key: "users", label: "Users", icon: <CouchUsersSubTabIcon/>},
                ]}
                active={subTab}
                onChange={(key) => onSubTabChange(key as CouchSubTab)}
            />
            {subTab === "instances" && <InstancesTab autoOpenAddDialog={autoOpenInstanceDialog}/>}
            {subTab === "users" && <UsersTab/>}
        </div>
    )
}
