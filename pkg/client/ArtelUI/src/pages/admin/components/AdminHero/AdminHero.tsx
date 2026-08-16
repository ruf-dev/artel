import cls from "@/pages/admin/components/AdminHero/AdminHero.module.css"
import {Tab, CouchSubTab} from "@/pages/admin/adminTypes.ts"

interface AdminHeroProps {
    tab: Tab
    subTab: CouchSubTab
}

export default function AdminHero({tab, subTab}: AdminHeroProps) {
    const titleMap: Record<Exclude<Tab, "couchdb">, string> = {
        users: "Artel users",
        s3_instances: "S3 storage",
        docker_api: "Docker API",
        settings: "Settings",
    }
    const title = tab === "couchdb"
        ? (subTab === "users" ? "CouchDB users" : "CouchDB instances")
        : titleMap[tab]
    return (
        <div className={cls.AdminHeroContainer}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>Admin</div>
                <h1 className={cls.HeroTitle}>{title}</h1>
            </div>
        </div>
    )
}
