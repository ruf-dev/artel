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
import ManageEmailDialog from "@/components/ManageEmailDialog/ManageEmailDialog.tsx"
import ProviderCard from "@/widgets/ProviderCard/ProviderCard.tsx"
import EmailCard from "@/widgets/EmailCard/EmailCard.tsx"

const PROVIDERS: {provider: ExternalProvider; name: string}[] = [
    {provider: ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS, name: "Google Sheets"},
    {provider: ExternalProvider.EXTERNAL_PROVIDER_TRELLO, name: "Trello"},
    {provider: ExternalProvider.EXTERNAL_PROVIDER_MIRO, name: "Miro"},
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

    const emailConnections = connections.filter(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_EMAIL)

    function findConnection(p: ExternalProvider) {
        return connections.find(c => c.provider === p)
    }

    return (
        <div className={cls.Content}>
            <div className={cls.Grid}>
                {PROVIDERS.map(({provider, name}) => (
                    <ProviderCard
                        key={provider}
                        provider={provider}
                        name={name}
                        connection={findConnection(provider)}
                        loading={loading}
                        onClick={() => OpenDialog(<ConnectionDetailDialog provider={provider}/>)}
                    />
                ))}
                <EmailCard
                    connections={emailConnections}
                    loading={loading}
                    onClick={() => OpenDialog(<ManageEmailDialog/>)}
                />
            </div>
        </div>
    )
}
