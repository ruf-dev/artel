import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import cls from "@/components/ManageVaultS3Section/ManageVaultS3Section.module.css"
import LinkedView from "@/components/ManageVaultS3Section/components/LinkedView"
import LinkForm from "@/components/ManageVaultS3Section/components/LinkForm"

export interface S3LinkPatch {
    s3InstanceId?: string
    s3BucketName?: string
}

interface Props {
    vault: VaultItem
    onLinked: (patch: S3LinkPatch) => void
}

export default function ManageVaultS3Section({vault, onLinked}: Props) {
    const linked = !!vault.s3InstanceId

    return (
        <section className={cls.ManageVaultS3SectionContainer}>
            <div className={cls.SectionHead}>
                <span className={cls.SectionTitle}>S3 storage</span>
            </div>
            {linked
                ? <LinkedView vault={vault} onLinked={onLinked}/>
                : <LinkForm vault={vault} onLinked={onLinked}/>
            }
        </section>
    )
}
