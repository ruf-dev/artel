import {useEffect} from "react"
import {useNavigate, useSearchParams} from "react-router-dom"
import {externalConnectionsService} from "@/processes/ExternalConnections.ts"

export default function GoogleOAuthCallbackPage() {
    const [searchParams] = useSearchParams()
    const navigate = useNavigate()

    useEffect(() => {
        const code = searchParams.get("code") ?? ""
        const state = searchParams.get("state") ?? ""

        externalConnectionsService.exchangeGoogleOAuth(code, state)
            .then(() => navigate("/connections?status=connected", {replace: true}))
            .catch(() => navigate("/connections?status=error", {replace: true}))
    }, [])

    return null
}
