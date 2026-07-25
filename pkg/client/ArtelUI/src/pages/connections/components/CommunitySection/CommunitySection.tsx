import {useEffect} from "react"
import {LoadingWrapper} from "@vervstack/chures"

import cls from "@/pages/connections/components/CommunitySection/CommunitySection.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import ConnectGenericDialog from "@/dialogs/ConnectGenericDialog/ConnectGenericDialog.tsx"
import CommunityConnectorCard from "@/widgets/CommunityConnectorCard/CommunityConnectorCard.tsx"

export default function CommunitySection() {
    const {OpenDialog} = useDialog()
    const bakeError = useBakeError()
    const {
        communityConnectors, communityConnectorsLoading, fetchCommunityConnectors, deleteCommunityConnector,
    } = useMcpKeys()

    useEffect(() => {
        fetchCommunityConnectors().catch(e => bakeError("Failed to load community connectors", e))
    }, [])

    function handleConnect(name: string) {
        OpenDialog(<ConnectGenericDialog name={name} onSuccess={fetchCommunityConnectors}/>)
    }

    return (
        <div className={cls.CommunitySectionContainer}>
            <LoadingWrapper isLoading={communityConnectorsLoading}>
                {communityConnectors.length === 0
                    ? (
                        <p className={cls.Empty}>
                            No community connectors yet. Ask an admin to add one via the MCP tool.
                        </p>
                    )
                    : (
                        <div className={cls.Grid}>
                            {communityConnectors.map(connector => (
                                <CommunityConnectorCard
                                    key={connector.name}
                                    name={connector.name ?? ""}
                                    author={connector.author ?? ""}
                                    description={connector.description ?? ""}
                                    toolCount={connector.tools?.length ?? 0}
                                    isConnected={connector.viewerIsConnected ?? false}
                                    viewerIsOwner={connector.viewerIsOwner ?? false}
                                    onConnect={() => handleConnect(connector.name ?? "")}
                                    onDelete={() => deleteCommunityConnector(connector.name ?? "")}
                                />
                            ))}
                        </div>
                    )}
            </LoadingWrapper>
        </div>
    )
}
