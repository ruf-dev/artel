import {useEffect} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/pages/mcp-keys/McpKeysPage.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import useUser from "@/hooks/user/User.ts"
import HeroSegment from "@/components/HeroSegment/HeroSegment.tsx"
import ContentSegment from "@/pages/mcp-keys/components/ContentSegment/ContentSegment.tsx"
import CreateKeyDialog from "@/pages/mcp-keys/components/CreateKeyDialog/CreateKeyDialog.tsx"

export default function McpKeysPage() {
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const {keys, loading, fetch: fetchKeys} = useMcpKeys()
    const {fetch: fetchExternalConnections} = useExternalConnections()

    // Unauthenticated → /init is handled at the router level (HomeLayout.tsx), not per-page.
    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetchKeys()
            void fetchExternalConnections()
        }
    }, [auth, fetchKeys, fetchExternalConnections])

    return (
        <div className={cls.Root}>
            <HeroSegment
                eyebrow="MCP"
                title="API Keys"
                subtitle={
                    <>
                        <b>{loading ? "…" : `${keys.length} ${keys.length === 1 ? "key" : "keys"}`}</b>
                        {" · "}<span>bridge your MCP agents to Artel</span>
                    </>
                }
                action={
                    <Button variant="primary" onClick={() => OpenDialog(<CreateKeyDialog/>)}>
                        <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor"
                             strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <line x1="12" y1="5" x2="12" y2="19"/>
                            <line x1="5" y1="12" x2="19" y2="12"/>
                        </svg>
                        New key
                    </Button>
                }
            />
            <ContentSegment/>
        </div>
    )
}
