import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import ConnectionDetailDialog from "@/dialogs/ConnectionDetailDialog/ConnectionDetailDialog.tsx"
import ManageEmailDialog from "@/dialogs/ManageEmailDialog/ManageEmailDialog.tsx"
import ManageGitlabDialog from "@/dialogs/ManageGitlabDialog/ManageGitlabDialog.tsx"
import ProviderCard from "@/widgets/ProviderCard/ProviderCard.tsx"
import cls from "@/pages/connections/components/ContentSegment/ContentSegment.module.css"

const PROVIDERS: {provider: ExternalProvider; name: string}[] = [
    {provider: ExternalProvider.EXTERNAL_PROVIDER_EMAIL, name: "Email"},
    {provider: ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS, name: "Google Sheets"},
    {provider: ExternalProvider.EXTERNAL_PROVIDER_MIRO, name: "Miro"},
    {provider: ExternalProvider.EXTERNAL_PROVIDER_TRELLO, name: "Trello"},
    {provider: ExternalProvider.EXTERNAL_PROVIDER_GITLAB, name: "GitLab"},
]

export default function ContentSegment() {
    const {connections, loading} = useExternalConnections()
    const {OpenDialog} = useDialog()

    function getConnections(p: ExternalProvider) {
        return connections.filter(c => c.provider === p)
    }

    function getDialog(provider: ExternalProvider) {
        if (provider === ExternalProvider.EXTERNAL_PROVIDER_EMAIL) {
            return <ManageEmailDialog/>
        }
        if (provider === ExternalProvider.EXTERNAL_PROVIDER_GITLAB) {
            return <ManageGitlabDialog/>
        }
        return <ConnectionDetailDialog provider={provider}/>
    }

    const sorted = [...PROVIDERS].sort((a, b) => {
        const aConnected = getConnections(a.provider).length > 0
        const bConnected = getConnections(b.provider).length > 0
        if (aConnected !== bConnected) return aConnected ? -1 : 1
        return a.name.localeCompare(b.name)
    })

    return (
        <div className={cls.ContentSegmentContainer}>
            <div className={cls.Grid}>
                {sorted.map(({provider, name}) => (
                    <ProviderCard
                        key={provider}
                        provider={provider}
                        name={name}
                        connections={getConnections(provider)}
                        loading={loading}
                        onClick={() => OpenDialog(getDialog(provider))}
                    />
                ))}
            </div>
        </div>
    )
}
