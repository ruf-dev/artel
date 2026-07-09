import {Button} from "@vervstack/chures"

import {Vault} from "@/pages/mcp-auth/processes/vault.ts"
import cls from "@/pages/mcp-auth/components/VaultList/VaultList.module.css"

interface VaultListProps {
    vaults: Vault[]
    loading: boolean
    onSelect: (id: string) => void
}

export default function VaultList({vaults, loading, onSelect}: VaultListProps) {
    if (vaults.length === 0) {
        return (
            <p className={cls.Empty}>
                No vaults found. Create one in Artel first.
            </p>
        )
    }
    return (
        <div className={cls.VaultListContainer}>
            {vaults.map(function (v) {
                return (
                    <Button
                        key={v.id}
                        variant="ghost"
                        className={cls.SubmitBtn}
                        disabled={loading}
                        onClick={function () { onSelect(v.id) }}
                    >
                        {v.name}
                    </Button>
                )
            })}
        </div>
    )
}
