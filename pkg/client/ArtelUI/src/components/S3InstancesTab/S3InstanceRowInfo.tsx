import cls from "@/components/S3InstancesTab/S3InstanceRowInfo.module.css"
import {GetS3InstanceResponse} from "@/app/api/artel/s3_instances.pb.ts"
import type {TestStatus} from "@/components/S3InstancesTab/S3InstanceRow.tsx"

export default function S3InstanceRowInfo(
    {instance, testStatus}: {instance: GetS3InstanceResponse; testStatus: TestStatus},
) {
    return (
        <div className={cls.S3InstanceRowInfoContainer}>
            <span className={cls.RowEndpoint}>{instance.endpoint}</span>
            <span className={cls.RowMeta}>
                {instance.region || "no region"} · added {instance.createdAt?.slice(0, 10)}
            </span>
            <div className={cls.BadgeRow}>
                {instance.useSsl && <span className={cls.Badge}>SSL</span>}
                {instance.pathStyle && <span className={cls.Badge}>Path-style</span>}
                {testStatus === "ok" && <span className={cls.BadgeOk}>Connected</span>}
                {testStatus === "fail" && <span className={cls.BadgeFail}>Failed</span>}
            </div>
        </div>
    )
}
