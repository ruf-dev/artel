import {ReactNode} from "react"

import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import ProviderCard from "@/widgets/ProviderCard/ProviderCard.tsx"
import ComingSoonCard from "@/components/ComingSoonCard/ComingSoonCard.tsx"
import ManageAnthropicDialog from "@/dialogs/ManageAnthropicDialog/ManageAnthropicDialog.tsx"
import ManageOpenAIDialog from "@/dialogs/ManageOpenAIDialog/ManageOpenAIDialog.tsx"
import ManageOpenRouterDialog from "@/dialogs/ManageOpenRouterDialog/ManageOpenRouterDialog.tsx"
import ManageCouchDBDialog from "@/dialogs/ManageCouchDBDialog/ManageCouchDBDialog.tsx"
import ManageS3Dialog from "@/dialogs/ManageS3Dialog/ManageS3Dialog.tsx"
import ManagePostgresDialog from "@/dialogs/ManagePostgresDialog/ManagePostgresDialog.tsx"
import claudeIcon from "@/icons/llm/claude-color.svg"
import gptIcon from "@/icons/llm/gpt-color.svg"
import openrouterIcon from "@/icons/llm/openrouter-color.svg"
import cls from "@/pages/connections/components/BYOKSection/BYOKSection.module.css"

type BYOKProvider = {
    provider: ExternalProvider
    name: string
    icon?: ReactNode
    dialogFactory: () => React.ReactElement
}

const LLM_BYOK_PROVIDERS: BYOKProvider[] = [
    {
        provider: ExternalProvider.EXTERNAL_PROVIDER_OPENROUTER,
        name: "OpenRouter",
        icon: <img src={openrouterIcon} alt="OpenRouter" className={cls.ProviderIconImage}/>,
        dialogFactory: () => <ManageOpenRouterDialog/>,
    },
    {
        provider: ExternalProvider.EXTERNAL_PROVIDER_ANTHROPIC,
        name: "Claude (Anthropic)",
        icon: <img src={claudeIcon} alt="Claude" className={cls.ProviderIconImage}/>,
        dialogFactory: () => <ManageAnthropicDialog/>,
    },
    {
        provider: ExternalProvider.EXTERNAL_PROVIDER_OPENAI,
        name: "GPT (OpenAI)",
        icon: <img src={gptIcon} alt="GPT" className={cls.ProviderIconImage}/>,
        dialogFactory: () => <ManageOpenAIDialog/>,
    },
]

const INFRA_BYOK_PROVIDERS: BYOKProvider[] = [
    {
        provider: ExternalProvider.EXTERNAL_PROVIDER_COUCHDB,
        name: "CouchDB",
        dialogFactory: () => <ManageCouchDBDialog/>,
    },
    {
        provider: ExternalProvider.EXTERNAL_PROVIDER_S3,
        name: "S3 / MinIO",
        dialogFactory: () => <ManageS3Dialog/>,
    },
    {
        provider: ExternalProvider.EXTERNAL_PROVIDER_POSTGRES,
        name: "PostgreSQL",
        dialogFactory: () => <ManagePostgresDialog/>,
    },
]

const COMING_SOON_CARDS: {icon: string; name: string; hint: string}[] = [
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
                    {LLM_BYOK_PROVIDERS.map(({provider, name, icon, dialogFactory}) => (
                        <ProviderCard
                            key={provider}
                            provider={provider}
                            name={name}
                            icon={icon}
                            connections={getConnections(provider)}
                            loading={loading}
                            onClick={() => OpenDialog(dialogFactory())}
                        />
                    ))}
                </div>
            </section>

            <section className={cls.Section}>
                <h2 className={cls.SectionTitle}>Infrastructure</h2>
                <div className={cls.Grid}>
                    {INFRA_BYOK_PROVIDERS.map(({provider, name, icon, dialogFactory}) => (
                        <ProviderCard
                            key={provider}
                            provider={provider}
                            name={name}
                            icon={icon}
                            connections={getConnections(provider)}
                            loading={loading}
                            onClick={() => OpenDialog(dialogFactory())}
                        />
                    ))}
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
