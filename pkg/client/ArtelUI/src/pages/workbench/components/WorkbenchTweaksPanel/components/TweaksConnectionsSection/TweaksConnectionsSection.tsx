import {useEffect} from "react"

import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import ManageOpenRouterDialog from "@/dialogs/ManageOpenRouterDialog/ManageOpenRouterDialog.tsx"
import type {WorkbenchMode} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"
import TweaksSection from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksSection/TweaksSection.tsx"
import ConnectionRow
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksConnectionsSection/components/ConnectionRow/ConnectionRow.tsx"
import cls
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksConnectionsSection/TweaksConnectionsSection.module.css"

interface Props {
    effectiveMode: WorkbenchMode | "picking"
    status: string
}

// Connections row of the Tweaks panel: the Docker container's state (docker mode
// only) plus the OpenRouter BYOK key, with "Manage"/"Connect" opening the shared
// ManageOpenRouterDialog.
export default function TweaksConnectionsSection({effectiveMode, status}: Props) {
    const {connections, fetch} = useExternalConnections()
    const {OpenDialog} = useDialog()

    useEffect(() => {
        if (connections.length === 0) void fetch()
    }, [])

    const openrouter = connections.find(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_OPENROUTER)
    const keyPreview = openrouter?.generic?.fields?.key_preview

    return (
        <TweaksSection label="Connections">
            <div className={cls.TweaksConnectionsSectionContainer}>
                {effectiveMode === "docker" && (
                    <ConnectionRow
                        name="Claude Code container"
                        detail={status}
                        online={status === "running"}
                    />
                )}
                {openrouter
                    ? (
                        <ConnectionRow
                            name="OpenRouter"
                            detail={keyPreview}
                            online
                            actionLabel="Manage"
                            onAction={() => OpenDialog(<ManageOpenRouterDialog/>)}
                        />
                    )
                    : (
                        <ConnectionRow
                            name="OpenRouter"
                            detail="Not connected"
                            actionLabel="Connect"
                            onAction={() => OpenDialog(<ManageOpenRouterDialog/>)}
                        />
                    )}
            </div>
        </TweaksSection>
    )
}
