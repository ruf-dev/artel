import {ReactNode} from "react"

import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import ProviderCard from "@/widgets/ProviderCard/ProviderCard.tsx"
import ComingSoonCard from "@/components/ComingSoonCard/ComingSoonCard.tsx"
import ManageAnthropicDialog from "@/dialogs/ManageAnthropicDialog/ManageAnthropicDialog.tsx"
import claudeIcon from "@/icons/llm/claude-color.svg"
import cls from "@/pages/connections/components/BYOKSection/BYOKSection.module.css"

const LLM_BYOK_PROVIDERS: {provider: ExternalProvider; name: string; icon: ReactNode}[] = [
    {
        provider: ExternalProvider.EXTERNAL_PROVIDER_ANTHROPIC,
        name: "Claude (Anthropic)",
        icon: <img src={claudeIcon} alt="Claude" className={cls.ProviderIconImage}/>,
    },
]

const COMING_SOON_CARDS: {icon: string; name: string; hint: string}[] = [
    {icon: "🗄️", name: "CouchDB", hint: "Bring your own CouchDB instance"},
    {icon: "🪣", name: "S3 / MinIO", hint: "Bring your own S3-compatible bucket"},
    {icon: "📁", name: "WebDAV", hint: "Bring your own WebDAV server"},
]

export default function BYOKSection() {
    const {connections, loading} = useExternalConnections()
    const {OpenDialog} = useDialog()

    function getConnections(p: ExternalProvider) {
        return connections.filter(c => c.provider === p)
    }

    return (
        <div className={cls.BYOKSectionContainer}>
            <section className={cls.Section}>
                <h2 className={cls.SectionTitle}>LLM Keys</h2>
                <div className={cls.Grid}>
                    {LLM_BYOK_PROVIDERS.map(({provider, name, icon}) => (
                        <ProviderCard
                            key={provider}
                            provider={provider}
                            name={name}
                            icon={icon}
                            connections={getConnections(provider)}
                            loading={loading}
                            onClick={() => OpenDialog(<ManageAnthropicDialog/>)}
                        />
                    ))}
                </div>
            </section>

            <section className={cls.Section}>
                <h2 className={cls.SectionTitle}>Infrastructure</h2>
                <div className={cls.Grid}>
                    {COMING_SOON_CARDS.map(({icon, name, hint}) => (
                        <ComingSoonCard
                            key={name}
                            icon={icon}
                            name={name}
                            hint={hint}
                        />
                    ))}
                </div>
            </section>
        </div>
    )
}
