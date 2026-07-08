import {useEffect} from "react"

import cls from "@/widgets/McpKeyCard/McpKeyCard.module.css"

import {McpKeyInfo} from "@/app/api/artel/mcp_keys.pb.ts"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useVaults} from "@/app/hooks/Vaults.ts"

import CardHeader from "@/components/McpKeyCard/CardHeader.tsx"
import CardChips from "@/components/McpKeyCard/CardChips.tsx"
import CardMeta from "@/components/McpKeyCard/CardMeta.tsx"

interface Props {
    mcpKey: McpKeyInfo
    onManage: () => void
}

export default function McpKeyCard({mcpKey, onManage}: Props) {
    const {vaults} = useVaults()
    const {connectorsByKey, fetchConnectors} = useMcpKeys()
    const {connections} = useExternalConnections()

    const vault = vaults.find(v => v.id === mcpKey.vaultId)
    const connectors = mcpKey.id ? connectorsByKey[mcpKey.id] ?? [] : []

    useEffect(() => {
        if (mcpKey.id) {
            void fetchConnectors(mcpKey.id)
        }
    }, [mcpKey.id, fetchConnectors])

    return (
        <div className={cls.Card} onClick={onManage}>
            <div className={cls.CardMain}>
                <CardHeader name={mcpKey.name} keyPreview={mcpKey.keyPreview}/>
                <CardChips vault={vault} connectors={connectors} connections={connections}/>
                <CardMeta lastAccessedAt={mcpKey.lastAccessedAt}/>
            </div>
        </div>
    )
}
