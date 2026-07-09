import {useState, useEffect} from "react"

import cls from "@/pages/admin/components/SetupStatusSection/SetupStatusSection.module.css"
import {CouchInstancesAPI, GetCouchInstanceStatusResponse} from "@/app/api/artel/couch_instances.pb.ts"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"

interface SetupStatusSectionProps {
    instanceId: string
    refreshKey: number
}

export default function SetupStatusSection({instanceId, refreshKey}: SetupStatusSectionProps) {
    const {auth} = useUser()
    const bakeError = useBakeError()
    const [status, setStatus] = useState<GetCouchInstanceStatusResponse | null>(null)

    useEffect(() => {
        async function loadStatus() {
            try {
                const res = await CouchInstancesAPI.GetCouchInstanceStatus({id: instanceId}, auth.getInitReq())
                setStatus(res)
            } catch (err) {
                bakeError("Failed to load instance status", err)
            }
        }
        void loadStatus()
    }, [instanceId, auth, refreshKey, bakeError])

    const checks = [
        {label: "Single-node cluster", ok: status?.clusterModeEnabled},
        {label: "_users database", ok: status?.usersDbInitialized},
        {label: "_replicator database", ok: status?.replicatorDbInitialized},
    ]

    return (
        <div className={cls.SetupStatusSectionContainer}>
            <div className={cls.StatusSectionTitle}>Setup status</div>
            {checks.map(check => (
                <div key={check.label} className={cls.StatusRow}>
                    <span className={cls.StatusLabel}>{check.label}</span>
                    <span className={status === null ? cls.StatusLoading : check.ok ? cls.StatusOk : cls.StatusFail}>
                        {status === null ? "…" : check.ok ? "✓" : "✗"}
                    </span>
                </div>
            ))}
        </div>
    )
}
