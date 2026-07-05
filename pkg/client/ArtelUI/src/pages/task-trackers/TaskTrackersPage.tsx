import {useState, useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/task-trackers/TaskTrackersPage.module.css"

import {TaskTrackerInfo, TrelloBoardInfo} from "@/app/api/artel/task_trackers.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useTaskTrackers} from "@/app/hooks/TaskTrackers.ts"
import useUser from "@/hooks/user/User.ts"

import {ModalClose, CheckmarkIcon} from "@vervstack/chures"

export default function TaskTrackersPage() {
    const navigate = useNavigate()
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const {fetch: fetchTrackers} = useTaskTrackers()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetchTrackers()
        }
    }, [auth, fetchTrackers])

    return (
        <div className={cls.Root}>
            <HeroSegment onAddClick={() => OpenDialog(<AddTrackerDialog/>)}/>
            <ContentSegment/>
        </div>
    )
}

function HeroSegment({onAddClick}: { onAddClick: () => void }) {
    const {trackers, loading} = useTaskTrackers()

    return (
        <div className={cls.Hero}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>Workspace</div>
                <h1 className={cls.HeroTitle}>Task trackers</h1>
                <p className={cls.HeroSub}>
                    <b>{loading ? "…" : `${trackers.length} ${trackers.length === 1 ? "connection" : "connections"}`}</b>
                    {" · "}<span>connected task tracking services</span>
                </p>
            </div>
            <button className={cls.AddBtn} onClick={onAddClick}>
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2"
                     strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Add connection
            </button>
        </div>
    )
}

function ContentSegment() {
    const {trackers, loading, remove} = useTaskTrackers()

    if (loading) {
        return (
            <div className={cls.Content}>
                <p className={cls.Empty}>Loading…</p>
            </div>
        )
    }

    return (
        <div className={cls.Content}>
            <div className={cls.List}>
                {trackers.map(t => (
                    <TrackerRow key={t.id} tracker={t} onRemove={remove}/>
                ))}
                {trackers.length === 0 && (
                    <p className={cls.Empty}>No task trackers connected yet. Add one to get started.</p>
                )}
            </div>
        </div>
    )
}

function TrackerRow({tracker, onRemove}: {
    tracker: TaskTrackerInfo
    onRemove: (id: string) => Promise<void>
}) {
    const [removing, setRemoving] = useState(false)

    async function handleRemove() {
        if (!tracker.id) return
        setRemoving(true)
        try {
            await onRemove(tracker.id)
        } finally {
            setRemoving(false)
        }
    }

    return (
        <div className={cls.Row}>
            <div className={cls.RowInfo}>
                <span className={cls.RowName}>{tracker.name}</span>
                <span className={cls.RowMeta}>{tracker.type} · connected {tracker.createdAt}</span>
            </div>
            <button className={cls.BtnDanger} onClick={handleRemove} disabled={removing} type="button">
                {removing ? "Removing…" : "Remove"}
            </button>
        </div>
    )
}

type Step = "choose" | "trello" | "done"

function AddTrackerDialog() {
    const [step, setStep] = useState<Step>("choose")
    const [boards, setBoards] = useState<TrelloBoardInfo[]>([])

    const {CloseDialog} = useDialog()

    function handleSuccess(connectedBoards: TrelloBoardInfo[]) {
        setBoards(connectedBoards)
        setStep("done")
    }

    if (step === "choose") {
        return <ChooseTypeStep onChooseTrello={() => setStep("trello")} onClose={CloseDialog}/>
    }

    if (step === "trello") {
        return <TrelloSetupStep onSuccess={handleSuccess} onBack={() => setStep("choose")} onClose={CloseDialog}/>
    }

    return <DoneStep boards={boards} onClose={CloseDialog}/>
}

function ChooseTypeStep({onChooseTrello, onClose}: {
    onChooseTrello: () => void
    onClose: () => void
}) {
    return (
        <div className={cls.Overlay}>
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                 aria-labelledby="chooseTrackerTitle">
                <div className={cls.ModalHead}>
                    <h2 className={cls.ModalTitle} id="chooseTrackerTitle">Add connection</h2>
                    <ModalClose onClick={onClose}/>
                </div>
                <p className={cls.ModalSub}>Choose a task tracking service to connect.</p>

                <TypeGrid onChooseTrello={onChooseTrello}/>
            </div>
        </div>
    )
}

function TypeGrid({onChooseTrello}: {onChooseTrello: () => void}) {
    return (
        <div className={cls.TypeGrid}>
            <button className={cls.TypeCard} onClick={onChooseTrello} type="button">
                <div className={cls.TypeCardIcon}>T</div>
                <TypeCardText label="Trello" desc="Connect via API key and token"/>
            </button>
        </div>
    )
}

function TypeCardText({label, desc}: {label: string; desc: string}) {
    return (
        <div className={cls.TypeCardText}>
            <div className={cls.TypeCardLabel}>{label}</div>
            <div className={cls.TypeCardDesc}>{desc}</div>
        </div>
    )
}

function TrelloSetupStep({onSuccess, onBack, onClose}: {
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
        <div className={cls.Overlay}>
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
                    <input
                        className={cls.Input}
                        placeholder="Your Trello API key"
                        value={apiKey}
                        onChange={e => setApiKey(e.target.value)}
                        disabled={connecting}
                        autoComplete="off"
                    />
                </label>

                <label className={cls.Field}>
                    <span className={cls.FieldLabel}>API Token</span>
                    <input
                        className={cls.Input}
                        placeholder="Your Trello API token"
                        value={apiToken}
                        onChange={e => setApiToken(e.target.value)}
                        disabled={connecting}
                        autoComplete="off"
                    />
                </label>

                {error && (
                    <p style={{color: "var(--color-error)", fontSize: "0.8125rem", marginBottom: "1rem"}}>{error}</p>
                )}

                <div className={cls.ModalActions}>
                    <button className={cls.BtnSecondary} onClick={onBack} disabled={connecting} type="button">
                        Back
                    </button>
                    <button
                        className={cls.BtnPrimary}
                        onClick={handleConnect}
                        disabled={connecting || !apiKey || !apiToken}
                        type="button"
                    >
                        {connecting ? "Connecting…" : "Connect"}
                    </button>
                </div>
            </div>
        </div>
    )
}

function SuccessBannerText({boardCount}: {boardCount: number}) {
    return (
        <span>
            Successfully connected.{" "}
            <b>{boardCount} {boardCount === 1 ? "board" : "boards"}</b> available.
        </span>
    )
}

function DoneStep({boards, onClose}: {
    boards: TrelloBoardInfo[]
    onClose: () => void
}) {
    return (
        <div className={cls.Overlay}>
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                 aria-labelledby="doneTitleTT">
                <div className={cls.ModalHead}>
                    <h2 className={cls.ModalTitle} id="doneTitleTT">Trello connected</h2>
                    <ModalClose onClick={onClose}/>
                </div>

                <div className={cls.SuccessBanner}>
                    <CheckmarkIcon/>
                    <SuccessBannerText boardCount={boards.length}/>
                </div>

                <div className={cls.ModalActions}>
                    <button className={cls.BtnPrimary} onClick={onClose} type="button">Done</button>
                </div>
            </div>
        </div>
    )
}
