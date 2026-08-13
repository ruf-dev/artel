import {Navigate} from "react-router-dom"
import {useQuery} from "@tanstack/react-query"
import {LoadingWrapper} from "@vervstack/chures"

import {PublicDocsAPI} from "@/app/api/artel"
import DocsEmptyState from "@/pages/docs/components/DocsEmptyState/DocsEmptyState.tsx"
import cls from "@/pages/docs/DocsDefaultRedirect.module.css"

// Element rendered at the unauthenticated `/docs` route (Path.DocsPageDefault). Resolves the
// admin-configured default public vault and forwards to its real `/docs/:slug` URL so the
// shareable-URL behavior is unchanged — `/docs` still lands on a concrete slug. Any rejection
// (no default configured, or the configured vault is no longer public) falls back to the empty
// state instead of erroring out — same "anonymous, no auth" contract as DocsPage/usePublicDocs.
export default function DocsDefaultRedirect() {
    const defaultVaultQuery = useQuery({
        queryKey: ['public-docs', 'default-vault'],
        queryFn: () => PublicDocsAPI.GetDefaultVault({}),
        retry: false,
    })

    return (
        <div className={cls.DocsDefaultRedirectContainer}>
            <LoadingWrapper isLoading={defaultVaultQuery.isLoading}>
                {defaultVaultQuery.error || !defaultVaultQuery.data?.slug
                    ? <DocsEmptyState/>
                    : <Navigate to={`/docs/${defaultVaultQuery.data.slug}`} replace/>}
            </LoadingWrapper>
        </div>
    )
}
