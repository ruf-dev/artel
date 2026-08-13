import {useEffect, useState} from "react"
import {Button, LoadingWrapper} from "@vervstack/chures"

import {DocsSource} from "@/app/api/artel/public_docs.pb.ts"
import {VaultItem, VaultsAPI} from "@/app/api/artel/vaults.pb.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {cn} from "@/app/utils/cn.ts"
import useUser from "@/hooks/user/User.ts"
import cls from "@/pages/admin/components/SettingsTab/components/DefaultDocsVaultSection/DefaultDocsVaultSection.module.css" // eslint-disable-line max-len
import PublishedVaultsGroup from "@/pages/admin/components/SettingsTab/components/DefaultDocsVaultSection/components/PublishedVaultsGroup/PublishedVaultsGroup.tsx" // eslint-disable-line max-len
import PrivateVaultsGroup from "@/pages/admin/components/SettingsTab/components/DefaultDocsVaultSection/components/PrivateVaultsGroup/PrivateVaultsGroup.tsx" // eslint-disable-line max-len

interface Props {
    defaultDocsVaultId: string
    onDefaultDocsVaultChange: (vaultId: string) => void
    docsSource: DocsSource
    onDocsSourceChange: (source: DocsSource) => void
}

export default function DefaultDocsVaultSection({
    defaultDocsVaultId,
    onDefaultDocsVaultChange,
    docsSource,
    onDocsSourceChange,
}: Props) {
    const {auth} = useUser()
    const bakeError = useBakeError()
    const [vaults, setVaults] = useState<VaultItem[]>([])
    const [loading, setLoading] = useState(true)

    // Fetched locally instead of threaded down from SettingsTab: the vault list is
    // this section's own concern (no sibling section needs it), so it follows the
    // same "owns its related data" reasoning widgets use — even though this file
    // lives in the page-local `components/` tier rather than `widgets/`.
    useEffect(() => {
        VaultsAPI.ListVaults({}, auth.getInitReq())
            .then(res => setVaults(res.vaults ?? []))
            .catch(err => bakeError("Failed to load vaults", err))
            .finally(() => setLoading(false))
    }, [auth, bakeError])

    function handleVaultPublished(vault: VaultItem) {
        setVaults(prev => prev.map(v => (v.id === vault.id ? vault : v)))
    }

    const publishedVaults = vaults.filter(v => v.isPublic)
    const privateVaults = vaults.filter(v => !v.isPublic)
    const isVaultSource = docsSource === DocsSource.VAULT
    const isGithubSource = docsSource === DocsSource.GITHUB

    return (
        <section className={cls.DefaultDocsVaultSectionContainer}>
            <h2 className={cls.SectionTitle}>Default /docs Vault</h2>
            <div className={cls.SourceButtonGroup}>
                <Button
                    variant={isVaultSource ? "primary" : "secondary"}
                    onClick={() => onDocsSourceChange(DocsSource.VAULT)}
                    className={cn(cls.SourceButton, isVaultSource && cls.SourceButtonActive)}
                >
                    Vault
                </Button>
                <Button
                    variant={isGithubSource ? "primary" : "secondary"}
                    onClick={() => onDocsSourceChange(DocsSource.GITHUB)}
                    className={cn(cls.SourceButton, isGithubSource && cls.SourceButtonActive)}
                >
                    GitHub
                </Button>
            </div>
            <p className={cls.SectionDescription}>
                Choose what unauthenticated visitors land on at the bare /docs URL — a published vault, or the
                built-in guide.
            </p>
            {isGithubSource ? (
                <p className={cls.GithubNote}>
                    Using the built-in Artel Quick Start guide, served from GitHub. Switch to Vault to pick a
                    published vault instead.
                </p>
            ) : (
                <LoadingWrapper isLoading={loading}>
                    {vaults.length === 0 ? (
                        <p className={cls.Empty}>
                            You don&apos;t have any vaults yet — create one and publish it to set it as the default.
                        </p>
                    ) : (
                        <div className={cls.Groups}>
                            <PublishedVaultsGroup
                                vaults={publishedVaults}
                                defaultDocsVaultId={defaultDocsVaultId}
                                onSelect={onDefaultDocsVaultChange}
                                onClear={() => onDefaultDocsVaultChange("")}
                            />
                            <div className={cls.Separator}/>
                            <PrivateVaultsGroup vaults={privateVaults} onPublished={handleVaultPublished}/>
                        </div>
                    )}
                </LoadingWrapper>
            )}
        </section>
    )
}
