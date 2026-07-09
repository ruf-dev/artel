import {useState} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/widgets/VaultDangerZone/VaultDangerZone.module.css"
import {useVaultMutations} from "@/app/hooks/Vaults.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

interface Props {
    vaultId: string
    onDeleted: () => void
}

export default function VaultDangerZone({vaultId, onDeleted}: Props) {
    const {remove} = useVaultMutations()
    const bakeError = useBakeError()
    const [deleting, setDeleting] = useState(false)

    function handleDelete() {
        setDeleting(true)
        remove(vaultId)
            .then(onDeleted)
            .catch(e => bakeError('Error deleting vault', e))
            .finally(() => setDeleting(false))
    }

    return (
        <section className={cls.VaultDangerZoneContainer}>
            <div className={cls.VaultDangerZoneRow}>
                <DangerZoneText/>
                <Button variant="danger" onClick={handleDelete} disabled={deleting}>
                    {deleting ? "Deleting…" : "Delete"}
                </Button>
            </div>
        </section>
    )
}

function DangerZoneText() {
    return (
        <div className={cls.DangerZoneTextContainer}>
            <div className={cls.DangerZoneTitle}>Delete this vault</div>
            <div className={cls.DangerZoneSub}>Permanent. Connection string stops working immediately.</div>
        </div>
    )
}
