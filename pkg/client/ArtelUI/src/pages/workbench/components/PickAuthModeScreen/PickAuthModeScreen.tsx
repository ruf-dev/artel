import {useEffect, useState} from "react"
import {useNavigate} from "react-router-dom"
import {Button} from "@vervstack/chures"

import cls from "@/pages/workbench/components/PickAuthModeScreen/PickAuthModeScreen.module.css"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useWorkbenchMutations} from "@/app/hooks/Workbench.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {Path} from "@/app/routing/Router.tsx"
import AuthModeOption from "@/pages/workbench/components/AuthModeOption/AuthModeOption.tsx"

type AuthMode = "api_key" | "subscription_login"

interface Props {
    vaultId: string
    onStarted: (authMode: string) => void
}

export default function PickAuthModeScreen({vaultId, onStarted}: Props) {
    const {connections, fetch: fetchConnections} = useExternalConnections()
    const {start} = useWorkbenchMutations(vaultId)
    const navigate = useNavigate()
    const bakeError = useBakeError()

    const [authMode, setAuthMode] = useState<AuthMode>("subscription_login")
    const [touched, setTouched] = useState(false)
    const [starting, setStarting] = useState(false)

    useEffect(() => {
        void fetchConnections()
    }, [fetchConnections])

    const hasAnthropicConnection = connections
        .some(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_ANTHROPIC)

    // Default to the API key flow once we've confirmed a connection exists, unless
    // the user has already made an explicit choice.
    useEffect(() => {
        if (!touched && hasAnthropicConnection) {
            setAuthMode("api_key")
        }
    }, [hasAnthropicConnection, touched])

    function handleSelect(mode: AuthMode) {
        setTouched(true)
        setAuthMode(mode)
    }

    function handleStart() {
        setStarting(true)
        start(authMode)
            .then(() => {
                if (authMode === "subscription_login") {
                    onStarted(authMode)
                }
            })
            .catch(e => bakeError("Failed to start workbench", e))
            .finally(() => setStarting(false))
    }

    return (
        <div className={cls.PickAuthModeScreenContainer}>
            <p className={cls.ModalSub}>Choose how Claude Code should authenticate.</p>
            <div className={cls.OptionList}>
                <AuthModeOption
                    selected={authMode === "api_key"}
                    label="Use my Anthropic API key"
                    desc="Runs Claude Code with a connected Anthropic API key."
                    disabled={!hasAnthropicConnection}
                    hint="Connect an Anthropic key in Connections first"
                    onSelect={() => handleSelect("api_key")}
                />
                <AuthModeOption
                    selected={authMode === "subscription_login"}
                    label="Log in with my Claude subscription"
                    desc="Sign in with your Claude.ai account inside the Workbench."
                    onSelect={() => handleSelect("subscription_login")}
                />
            </div>
            <div className={cls.ModalFooter}>
                <Button variant="ghost" onClick={() => navigate(Path.HomePage)}>Cancel</Button>
                <Button variant="primary" onClick={handleStart} disabled={starting}>
                    {starting ? "Starting…" : "Start"}
                </Button>
            </div>
        </div>
    )
}
