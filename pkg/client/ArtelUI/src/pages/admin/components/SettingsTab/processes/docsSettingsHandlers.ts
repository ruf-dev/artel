import {AdminSystemSettingsAPI} from "@/app/api/artel/admin_system_settings.pb.ts"
import {DocsSource} from "@/app/api/artel/public_docs.pb.ts"
import {AuthMiddleware} from "@/processes/AuthMiddleware.ts"

interface DocsSettingsHandlersParams {
    auth: AuthMiddleware
    defaultDocsVaultId: string
    docsSource: DocsSource
    setDefaultDocsVaultId: (vaultId: string) => void
    setDocsSource: (source: DocsSource) => void
    bakeError: (title: string, err: unknown) => void
}

/**
 * Orchestrates the two ways `/docs` default source gets changed from SettingsTab
 * (picking a vault, or flipping the Vault/GitHub toggle) — both end up calling the
 * same UpdateDefaultDocsVault RPC with an optimistic update + revert-on-error,
 * differing only in which fields change. Colocated here (rather than inline in
 * SettingsTab.tsx) to keep the component focused on rendering.
 */
export function createDocsSettingsHandlers({
    auth,
    defaultDocsVaultId,
    docsSource,
    setDefaultDocsVaultId,
    setDocsSource,
    bakeError,
}: DocsSettingsHandlersParams) {
    function persist(vaultId: string, source: DocsSource, oldVaultId: string, oldSource: DocsSource) {
        AdminSystemSettingsAPI.UpdateDefaultDocsVault({vaultId, source}, auth.getInitReq())
            .catch(err => {
                bakeError("Failed to update default docs settings", err)
                setDefaultDocsVaultId(oldVaultId)
                setDocsSource(oldSource)
            })
    }

    function handleDefaultDocsVaultChange(vaultId: string) {
        if (defaultDocsVaultId === vaultId && docsSource === DocsSource.VAULT) return

        const oldVaultId = defaultDocsVaultId
        const oldSource = docsSource
        setDefaultDocsVaultId(vaultId)
        setDocsSource(DocsSource.VAULT)
        persist(vaultId, DocsSource.VAULT, oldVaultId, oldSource)
    }

    function handleDocsSourceChange(source: DocsSource) {
        if (docsSource === source) return

        const oldSource = docsSource
        setDocsSource(source)
        persist(defaultDocsVaultId, source, defaultDocsVaultId, oldSource)
    }

    return {handleDefaultDocsVaultChange, handleDocsSourceChange}
}
