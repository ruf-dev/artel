import cls from "@/pages/admin/components/AdminHero/AdminHero.module.css"
import {Tab} from "@/pages/admin/adminTypes.ts"

interface AdminHeroProps {
    tab: Tab
}

export default function AdminHero({tab}: AdminHeroProps) {
    const titleMap: Record<Tab, string> = {
        instances: "CouchDB instances",
        couch_users: "CouchDB users",
        users: "Artel users",
        s3_instances: "S3 storage",
        docker_api: "Docker API",
    }
    const title = titleMap[tab]
    return (
        <div className={cls.AdminHeroContainer}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>Admin</div>
                <h1 className={cls.HeroTitle}>{title}</h1>
            </div>
        </div>
    )
}
