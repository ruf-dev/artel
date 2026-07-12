import {useEffect, useState} from "react"

import cls from "@/dialogs/ManageVaultDialog/widgets/VaultSettingsSection/VaultSettingsSection.module.css"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {AuthAPI} from "@/app/api/artel"
import {useVaultMutations} from "@/app/hooks/Vaults.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import useUser from "@/hooks/user/User.ts"
import BinaryStorageToggle from
    "@/dialogs/ManageVaultDialog/widgets/VaultSettingsSection/components/BinaryStorageToggle/BinaryStorageToggle.tsx"

interface Props {
    vault: VaultItem
    onChanged: (patch: Partial<VaultItem>) => void
}

export default function VaultSettingsSection({vault, onChanged}: Props) {
    const {auth} = useUser()
    const {setBinaryStorage} = useVaultMutations()
    const bakeError = useBakeError()

    // Default to "unavailable" while loading so the toggle never briefly shows
    // enabled before we've confirmed S3 actually exists.
    const [isS3Available, setIsS3Available] = useState(false)
    const [useCouchDb, setUseCouchDb] = useState(vault.useCouchdbForBinaries ?? true)

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

    return (
        <section className={cls.VaultSettingsSectionContainer}>
            <div className={cls.SectionHead}>
                <span className={cls.SectionTitle}>Vault settings</span>
            </div>
            <BinaryStorageToggle
                checked={checked}
                disabled={!isS3Available}
                onChange={handleChange}
            />
        </section>
    )
}
