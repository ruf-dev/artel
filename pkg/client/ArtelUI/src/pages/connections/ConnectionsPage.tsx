import {useEffect} from "react"
import {useNavigate, useSearchParams} from "react-router-dom"

import cls from "@/pages/connections/ConnectionsPage.module.css"

import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import useUser from "@/hooks/user/User.ts"

import ConnectionDetailDialog from "@/dialogs/ConnectionDetailDialog/ConnectionDetailDialog.tsx"
import ManageEmailDialog from "@/dialogs/ManageEmailDialog/ManageEmailDialog.tsx"
import ManageGitlabDialog from "@/dialogs/ManageGitlabDialog/ManageGitlabDialog.tsx"
import ProviderCard from "@/widgets/ProviderCard/ProviderCard.tsx"

const PROVIDERS: {provider: ExternalProvider; name: string}[] = [
    {provider: ExternalProvider.EXTERNAL_PROVIDER_EMAIL, name: "Email"},
    {provider: ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS, name: "Google Sheets"},
    {provider: ExternalProvider.EXTERNAL_PROVIDER_MIRO, name: "Miro"},
    {provider: ExternalProvider.EXTERNAL_PROVIDER_TRELLO, name: "Trello"},
    {provider: ExternalProvider.EXTERNAL_PROVIDER_GITLAB, name: "GitLab"},
]

export default function ConnectionsPage() {
    const navigate = useNavigate()
    const {auth} = useUser()
    const {fetch: fetchConnections} = useExternalConnections()
    const bakeError = useBakeError()
    const [searchParams, setSearchParams] = useSearchParams()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    useEffect(() => {
        if (!auth.isAuthenticated()) return

        const status = searchParams.get("status")
        if (status === "error") {
            bakeError("Connection failed", new Error("Google OAuth authorization was denied or failed."))
        }
        if (status) {
            setSearchParams({}, {replace: true})
        }

        void fetchConnections()
    }, [auth])

    return (
        <div className={cls.Root}>
            <HeroSegment/>
            <ContentSegment/>
        </div>
    )
}

function HeroSegment() {
    const {connections, loading} = useExternalConnections()
    const connectedCount = connections.length

    return (
        <div className={cls.Hero}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>Workspace</div>
                <h1 className={cls.HeroTitle}>Connected services</h1>
                <p className={cls.HeroSub}>
                    <b>{loading ? "…" : `${connectedCount} ${connectedCount === 1 ? "service" : "services"}`}</b>
                    {" · "}<span>linked external integrations</span>
                </p>
            </div>
        </div>
    )
}

function ContentSegment() {
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
        <div className={cls.Content}>
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
