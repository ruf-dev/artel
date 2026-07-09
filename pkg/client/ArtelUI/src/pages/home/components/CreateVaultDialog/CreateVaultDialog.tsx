import {useState} from "react"
import {Button, ModalClose} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import {useVaultMutations} from "@/app/hooks/Vaults.ts"
import {useBakeError} from "@/app/hooks/useErrorToast"
import FormField from "@/components/FormField/FormField.tsx"
import cls from "@/pages/home/components/CreateVaultDialog/CreateVaultDialog.module.css"

export default function CreateVaultDialog() {
    const [isCreating, setIsCreating] = useState(false)
    const [vaultName, setVaultName] = useState("")

    const {create} = useVaultMutations()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

    async function handleCreate() {
        setIsCreating(true)
        const name = vaultName.trim()
        create(name)
            .then(CloseDialog)
            .catch(e => bakeError('Error creating vault', e))
            .finally(() => {
                setIsCreating(false)
            })
    }

    return (
        <div className={cls.CreateVaultDialogContainer}>
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                 aria-labelledby="createModalTitle">
                <div className={cls.ModalHead}>
                    <h2 className={cls.ModalTitle} id="createModalTitle">New vault</h2>
                    <ModalClose onClick={CloseDialog} disabled={isCreating} className={cls.ModalClose}/>
                </div>
                <p className={cls.ModalSub}>Give it a name — that's all we need to spin one up.</p>

                <FormField
                    label="Vault name"
                    placeholder="e.g. Marketplace inventory"
                    onChange={setVaultName}
                    disabled={isCreating}
                    maxLength={48}
                    fieldClassName={cls.Field}
                    labelClassName={cls.FieldLabel}
                    inputClassName={cls.Input}
                />

                <div className={cls.ModalActions}>
                    <Button variant="primary" onClick={handleCreate} disabled={isCreating}>
                        {isCreating ? "Creating…" : "Create vault"}
                    </Button>
                </div>
            </div>
        </div>
    )
}
