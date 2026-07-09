import {Button} from "@vervstack/chures"

import cls from "@/pages/admin/components/TabBar/TabBar.module.css"
import {Tab} from "@/pages/admin/adminTypes.ts"
import {cn} from "@/app/utils/cn.ts"

interface TabBarProps {
    tab: Tab
    onTabChange: (t: Tab) => void
}

export default function TabBar({tab, onTabChange}: TabBarProps) {
    return (
        <div className={cls.TabBarContainer}>
            <Button
                variant="ghost"
                className={cn(cls.Tab, tab === "instances" && cls.TabActive)}
                onClick={() => onTabChange("instances")}
            >
                Instances
            </Button>
            <Button
                variant="ghost"
                className={cn(cls.Tab, tab === "couch_users" && cls.TabActive)}
                onClick={() => onTabChange("couch_users")}
            >
                CouchUsers
            </Button>
            <Button
                variant="ghost"
                className={cn(cls.Tab, tab === "users" && cls.TabActive)}
                onClick={() => onTabChange("users")}
            >
                Users
            </Button>
            <Button
                variant="ghost"
                className={cn(cls.Tab, tab === "s3_instances" && cls.TabActive)}
                onClick={() => onTabChange("s3_instances")}
            >
                S3 storage
            </Button>
        </div>
    )
}
