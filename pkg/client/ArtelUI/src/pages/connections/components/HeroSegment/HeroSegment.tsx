import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import cls from "@/pages/connections/components/HeroSegment/HeroSegment.module.css"

export default function HeroSegment() {
    const {connections, loading} = useExternalConnections()
    const connectedCount = connections.length

    return (
        <div className={cls.HeroSegmentContainer}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>Workspace</div>
                <h1 className={cls.HeroTitle}>Connected services</h1>
                <p className={cls.HeroSub}>
                    <b>{loading ? "…" : `${connectedCount} ${connectedCount === 1 ? "service" : "services"}`}</b>
                    {" · "}<span>linked external integrations</span>
                </p>
            </div>
        </div>
    )
}
