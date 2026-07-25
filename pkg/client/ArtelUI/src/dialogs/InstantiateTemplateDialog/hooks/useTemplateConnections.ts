import {useEffect, useMemo, useState, SetStateAction, Dispatch} from "react"

import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {ExternalConnectionInfo, ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {
    ConnectionRequirement,
    requiredConnections,
} from "@/dialogs/InstantiateTemplateDialog/processes/requiredConnections.ts"
import type {TractStep} from "@/processes/Tracts.ts"

interface UseTemplateConnectionsResult {
    requirements: ConnectionRequirement[]
    candidatesByKey: Record<string, ExternalConnectionInfo[]>
    connectionsByKey: Record<string, string>
    setConnectionsByKey: Dispatch<SetStateAction<Record<string, string>>>
}

export function useTemplateConnections(steps: TractStep[]): UseTemplateConnectionsResult {
    const {momCandidates, fetchMomCandidates} = useMcpKeys()
    const {connections: externalConnections, fetch: fetchExternalConnections} = useExternalConnections()
    const [connectionsByKey, setConnectionsByKey] = useState<Record<string, string>>({})

    useEffect(() => {
        void fetchMomCandidates()
        void fetchExternalConnections()
    }, [fetchMomCandidates, fetchExternalConnections])

    useEffect(() => {
        function handleVisibilityChange() {
            if (document.visibilityState === "visible") {
                void fetchMomCandidates()
                void fetchExternalConnections()
            }
        }
        document.addEventListener("visibilitychange", handleVisibilityChange)
        return () => document.removeEventListener("visibilitychange", handleVisibilityChange)
    }, [fetchMomCandidates, fetchExternalConnections])

    const requirements = useMemo(
        () => requiredConnections(steps),
        [steps],
    )

    const anthropicConnections = useMemo(
        () => externalConnections.filter(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_ANTHROPIC),
        [externalConnections],
    )

    const candidatesByKey = useMemo(() => {
        const map: Record<string, ExternalConnectionInfo[]> = {}
        for (const req of requirements) {
            map[req.key] = req.kind === "llm"
                ? anthropicConnections
                : (momCandidates.find(c => c.name === req.key)?.connections ?? [])
        }
        return map
    }, [requirements, momCandidates, anthropicConnections])

    useEffect(() => {
        setConnectionsByKey(prev => {
            let changed = false
            const next = {...prev}
            for (const req of requirements) {
                const id = next[req.key]
                if (id && !candidatesByKey[req.key]?.some(c => c.id === id)) {
                    delete next[req.key]
                    changed = true
                }
            }
            return changed ? next : prev
        })
    }, [candidatesByKey, requirements])

    return {requirements, candidatesByKey, connectionsByKey, setConnectionsByKey}
}
