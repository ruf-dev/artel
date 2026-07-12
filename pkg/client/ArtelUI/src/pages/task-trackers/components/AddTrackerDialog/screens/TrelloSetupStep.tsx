import {useState} from "react"
import {Button, ModalClose} from "@vervstack/chures"

import {TrelloBoardInfo} from "@/app/api/artel/task_trackers.pb.ts"
import {useTaskTrackers} from "@/app/hooks/TaskTrackers.ts"
import Input from "@/components/atoms/Input/Input.tsx"
import cls from "@/pages/task-trackers/components/AddTrackerDialog/screens/TrelloSetupStep.module.css"

export default function TrelloSetupStep({onSuccess, onBack, onClose}: {
    onSuccess: (boards: TrelloBoardInfo[]) => void
    onBack: () => void
    onClose: () => void
}) {
    const [apiKey, setApiKey] = useState("")
    const [apiToken, setApiToken] = useState("")
    const [connecting, setConnecting] = useState(false)
    const [error, setError] = useState("")

    const {add} = useTaskTrackers()

    async function handleConnect() {
        setConnecting(true)
        setError("")
        try {
            const resp = await add({type: "trello", apiKey, apiToken})
            onSuccess(resp.boards ?? [])
        } catch (e: unknown) {
            const msg = e instanceof Error ? e.message : "Failed to connect. Check your credentials."
            setError(msg)
        } finally {
            setConnecting(false)
        }
    }

    return (
        <div className={cls.TrelloSetupStepContainer}>
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                 aria-labelledby="trelloSetupTitle">
                <div className={cls.ModalHead}>
                    <h2 className={cls.ModalTitle} id="trelloSetupTitle">Connect Trello</h2>
                    <ModalClose onClick={onClose} disabled={connecting}/>
                </div>
                <p className={cls.ModalSub}>
                    Get your credentials from the{" "}
                    <a href="https://trello.com/app-key" target="_blank" rel="noopener noreferrer">Trello API key page</a>
                    {". "}
                    See the{" "}
                    <a href="https://developer.atlassian.com/cloud/trello/guides/rest-api/api-introduction/" target="_blank" rel="noopener noreferrer">
                        Trello REST API docs
                    </a>
                    {" "}for details.
                </p>

                <label className={cls.Field}>
                    <span className={cls.FieldLabel}>API Key</span>
                    <Input
                        placeholder="Your Trello API key"
                        value={apiKey}
                        setValue={setApiKey}
                        disabled={connecting}
                        autoComplete="off"
                    />
                </label>

                <label className={cls.Field}>
                    <span className={cls.FieldLabel}>API Token</span>
                    <Input
                        placeholder="Your Trello API token"
                        value={apiToken}
                        setValue={setApiToken}
                        disabled={connecting}
                        autoComplete="off"
                    />
                </label>

                {error && (
                    <p className={cls.ErrorText}>{error}</p>
                )}

                <div className={cls.ModalActions}>
                    <Button variant="secondary" onClick={onBack} disabled={connecting}>
                        Back
                    </Button>
                    <Button
                        variant="primary"
                        onClick={handleConnect}
                        disabled={connecting || !apiKey || !apiToken}
                    >
                        {connecting ? "Connecting…" : "Connect"}
                    </Button>
                </div>
            </div>
        </div>
    )
}
