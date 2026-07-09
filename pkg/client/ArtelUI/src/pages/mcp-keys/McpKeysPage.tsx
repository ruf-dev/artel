import {useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/mcp-keys/McpKeysPage.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import useUser from "@/hooks/user/User.ts"
import HeroSegment from "@/pages/mcp-keys/components/HeroSegment/HeroSegment.tsx"
import ContentSegment from "@/pages/mcp-keys/components/ContentSegment/ContentSegment.tsx"
import CreateKeyDialog from "@/pages/mcp-keys/components/CreateKeyDialog/CreateKeyDialog.tsx"

export default function McpKeysPage() {
    const navigate = useNavigate()
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const {fetch: fetchKeys} = useMcpKeys()
    const {fetch: fetchExternalConnections} = useExternalConnections()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetchKeys()
            void fetchExternalConnections()
        }
    }, [auth, fetchKeys, fetchExternalConnections])

    return (
        <div className={cls.Root}>
            <HeroSegment onCreateClick={() => OpenDialog(<CreateKeyDialog/>)}/>
            <ContentSegment/>
        </div>
    )
}
