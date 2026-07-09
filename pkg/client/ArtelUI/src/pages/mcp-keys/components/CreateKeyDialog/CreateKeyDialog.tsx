import {useState} from "react"
import {Button, ModalClose} from "@vervstack/chures"

import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useVaults} from "@/app/hooks/Vaults.ts"
import FormField from "@/components/FormField/FormField.tsx"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import cls from "@/pages/mcp-keys/components/CreateKeyDialog/CreateKeyDialog.module.css"

export default function CreateKeyDialog() {
    const [creating, setCreating] = useState(false)
    const [name, setName] = useState("")
    const [selectedVaultId, setSelectedVaultId] = useState("")
    const [rawToken, setRawToken] = useState("")
    const [copied, setCopied] = useState(false)

    const {create} = useMcpKeys()
    const {CloseDialog} = useDialog()
    const {vaults} = useVaults()

    async function handleCreate() {
        if (!name || !selectedVaultId) return
        setCreating(true)
        try {
            const resp = await create(name, selectedVaultId)
            setRawToken(resp.rawToken ?? "")
        } finally {
            setCreating(false)
        }
    }

    async function handleCopy() {
        await navigator.clipboard.writeText(rawToken)
        setCopied(true)
        setTimeout(() => setCopied(false), 2000)
    }

    if (rawToken) {
        return (
            <div className={cls.Overlay}>
                <div
                    className={cls.Modal}
                    onClick={e =>
                        e.stopPropagation()}
                    role="dialog"
                    aria-modal="true"
                    aria-labelledby="createdKeyTitle">
                    <div className={cls.ModalHead}>
                        <h2 className={cls.ModalTitle} id="createdKeyTitle">Key created</h2>
                    </div>
                    <p className={cls.ModalSub}>
                        Copy this key now — it will not be shown again.
                    </p>
                    <div className={cls.TokenBox}>
                        <span className={cls.TokenText}>{rawToken}</span>
                        <Button variant="secondary" onClick={handleCopy}>
                            {copied ? "Copied!" : "Copy"}
                        </Button>
                    </div>
                    <div className={cls.ModalActions}>
                        <Button variant="primary" onClick={CloseDialog}>Done</Button>
                    </div>
                </div>
            </div>
        )
    }

    return (
        <div className={cls.Overlay}>
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                 aria-labelledby="createKeyTitle">
                <div className={cls.ModalHead}>
                    <h2 className={cls.ModalTitle} id="createKeyTitle">New API key</h2>
                    <ModalClose onClick={CloseDialog} disabled={creating} className={cls.ModalClose}/>
                </div>
                <p className={cls.ModalSub}>Name your key and pick the vault it will access.</p>

                <FormField
                    label="Key name"
                    placeholder="e.g. My Assistant"
                    onChange={setName}
                    disabled={creating}
                    fieldClassName={cls.Field}
                    labelClassName={cls.FieldLabel}
                    inputClassName={cls.Input}
                />

                <div className={cls.Field}>
                    <span className={cls.FieldLabel}>Vault</span>
                    <div className={cls.OptionList}>
                        {vaults.map((v: VaultItem) => (
                            <SelectOption
                                key={v.id}
                                label={v.name ?? ""}
                                selected={selectedVaultId === v.id}
                                onSelect={() => setSelectedVaultId(v.id ?? "")}
                            />
                        ))}
                        {vaults?.length === 0 && (
                            <p className={cls.Empty}>No vaults available. Create one first.</p>
                        )}
                    </div>
                </div>

                <div className={cls.ModalActions}>
                    <Button variant="primary" onClick={handleCreate} disabled={creating || !name || !selectedVaultId}>
                        {creating ? "Creating…" : "Create key"}
                    </Button>
                </div>
            </div>
        </div>
    )
}
