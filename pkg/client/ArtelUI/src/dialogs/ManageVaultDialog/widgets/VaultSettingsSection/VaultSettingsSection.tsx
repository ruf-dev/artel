import {useEffect, useState} from "react"
import {ConfirmDialog} from "@vervstack/chures"

import cls from "@/dialogs/ManageVaultDialog/widgets/VaultSettingsSection/VaultSettingsSection.module.css"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {AuthAPI} from "@/app/api/artel"
import {useVaultMutations} from "@/app/hooks/Vaults.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {useDialog} from "@/app/hooks/Dialog.ts"
import useUser from "@/hooks/user/User.ts"
import BinaryStorageToggle from
    "@/dialogs/ManageVaultDialog/widgets/VaultSettingsSection/components/BinaryStorageToggle/BinaryStorageToggle.tsx"
import PublishToggle from
    "@/dialogs/ManageVaultDialog/widgets/VaultSettingsSection/components/PublishToggle/PublishToggle.tsx"
import PublishSlugForm from
    "@/dialogs/ManageVaultDialog/widgets/VaultSettingsSection/components/PublishSlugForm/PublishSlugForm.tsx"
import PostgresSection from
    "@/dialogs/ManageVaultDialog/widgets/VaultSettingsSection/components/PostgresSection/PostgresSection.tsx"

interface Props {
    vault: VaultItem
    onChanged: (patch: Partial<VaultItem>) => void
}

export default function VaultSettingsSection({vault, onChanged}: Props) {
    const {auth} = useUser()
    const {setBinaryStorage, publish, unpublish} = useVaultMutations()
    const bakeError = useBakeError()
    const {OpenDialog, CloseDialog} = useDialog()

    // Default to "unavailable" while loading so the toggle never briefly shows
    // enabled before we've confirmed S3 actually exists.
    const [isS3Available, setIsS3Available] = useState(false)
    const [useCouchDb, setUseCouchDb] = useState(vault.useCouchdbForBinaries ?? true)

    const [slugFormOpen, setSlugFormOpen] = useState(false)
    const [publishing, setPublishing] = useState(false)
    const [postgresStatus, setPostgresStatus] = useState(vault.postgresStatus)

    useEffect(() => {
        AuthAPI.GetConfig({}, auth.getInitReq())
            .then(res => setIsS3Available(!!res.isS3Available))
            .catch(e => bakeError("Error loading storage config", e))
    }, [])

    function handleChange(v: boolean) {
        setUseCouchDb(v)
        setBinaryStorage(vault.id ?? "", v)
            .then(() => onChanged({useCouchdbForBinaries: v}))
            .catch(e => {
                setUseCouchDb(!v)
                bakeError("Error updating binary storage setting", e)
            })
    }

    // While S3 isn't confirmed available, force the toggle to appear on
    // (CouchDB) and disabled, without discarding the real underlying value.
    const checked = isS3Available ? useCouchDb : true

    function handlePublish(slug: string) {
        setPublishing(true)
        publish(vault.id ?? "", slug)
            .then(() => {
                onChanged({isPublic: true, slug})
                setSlugFormOpen(false)
            })
            .catch(e => bakeError("Error publishing vault", e))
            .finally(() => setPublishing(false))
    }

    function handleUnpublish() {
        return unpublish(vault.id ?? "")
            .then(() => onChanged({isPublic: false, slug: undefined}))
            .catch(e => bakeError("Error unpublishing vault", e))
    }

    function handlePublicToggle(next: boolean) {
        if (next) {
            setSlugFormOpen(true)
            return
        }
        OpenDialog(
            <ConfirmDialog
                title="Make vault private"
                message={"This removes public access to this vault. Anyone using the current "
                    + "public link will no longer be able to view it."}
                confirmLabel="Make private"
                danger
                onClose={CloseDialog}
                onConfirm={handleUnpublish}
            />
        )
    }

    return (
        <div className={cls.VaultSettingsSectionContainer}>
            <BinaryStorageToggle
                checked={checked}
                disabled={!isS3Available}
                onChange={handleChange}
            />
            <PublishToggle
                checked={vault.isPublic ?? false}
                disabled={slugFormOpen || publishing}
                onChange={handlePublicToggle}
            />
            {slugFormOpen && (
                <PublishSlugForm
                    vaultName={vault.name ?? ""}
                    submitting={publishing}
                    onPublish={handlePublish}
                    onCancel={() => setSlugFormOpen(false)}
                />
            )}
            <PostgresSection
                vaultId={vault.id ?? ""}
                postgresStatus={postgresStatus}
                onStatusChange={status => {
                    setPostgresStatus(status)
                    onChanged({postgresStatus: status})
                }}
            />
        </div>
    )
}
