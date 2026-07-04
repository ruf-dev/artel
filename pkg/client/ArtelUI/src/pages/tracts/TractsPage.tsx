import {useEffect, useState} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/tracts/TractsPage.module.css"

import {Path} from "@/app/routing/Router.tsx"
import {useDialog, useDialogKeyboard} from "@/app/hooks/Dialog"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import useUser from "@/hooks/user/User.ts"

import Button from "@/components/shared/Button/Button.tsx"
import TractCard from "@/widgets/TractCard/TractCard.tsx"

export default function TractsPage() {
    const navigate = useNavigate()
    const {auth} = useUser()
    const {fetch} = useTracts()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetch()
        }
    }, [auth, fetch])

    return (
        <div className={cls.Root}>
            <HeroSegment/>
            <ContentSegment/>
        </div>
    )
}

function HeroSegment() {
    const {tracts, loading} = useTracts()
    const {OpenDialog} = useDialog()
    const active = tracts.filter(t => t.enabled).length
    const paused = tracts.length - active

    return (
        <div className={cls.Hero}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>Workspace</div>
                <h1 className={cls.HeroTitle}>Tracts</h1>
                <p className={cls.HeroSub}>
                    <b>{loading ? "…" : `${tracts.length} ${tracts.length === 1 ? "tract" : "tracts"}`}</b>
                    {tracts.length > 0 && <>{" · "}<span>{active} active, {paused} paused</span></>}
                </p>
            </div>
            <Button variant="primary" onClick={() => OpenDialog(<NewTractDialog/>)}>
                New tract
            </Button>
        </div>
    )
}

function ContentSegment() {
    const {tracts, loading} = useTracts()
    const {OpenDialog} = useDialog()
    const navigate = useNavigate()

    if (loading) {
        return <div className={cls.Content}><p className={cls.Empty}>Loading…</p></div>
    }

    if (tracts.length === 0) {
        return (
            <div className={cls.Content}>
                <div className={cls.EmptyState}>
                    <h2 className={cls.EmptyStateTitle}>No tracts yet</h2>
                    <p className={cls.EmptyStateSub}>
                        Tracts connect a trigger — a webhook or a manual fire — to a tree of
                        actions and conditions. Create your first one to get started.
                    </p>
                    <Button variant="primary" onClick={() => OpenDialog(<NewTractDialog/>)}>
                        New tract
                    </Button>
                </div>
            </div>
        )
    }

    return (
        <div className={cls.Content}>
            <div className={cls.List}>
                {tracts.map(t => (
                    <TractCard key={t.uuid} tract={t} onClick={() => navigate(`/tracts/${t.uuid}`)}/>
                ))}
            </div>
        </div>
    )
}

function NewTractDialog() {
    const {CloseDialog} = useDialog()
    const {createTract} = useTracts()
    const bakeError = useBakeError()
    const navigate = useNavigate()
    const [name, setName] = useState("")
    const [creating, setCreating] = useState(false)

    function handleCreate() {
        const trimmed = name.trim()
        if (!trimmed) return
        setCreating(true)
        createTract(trimmed, "", {steps: []})
            .then(({tract}) => {
                CloseDialog()
                navigate(`/tracts/${tract.uuid}`)
            })
            .catch(err => bakeError("Failed to create tract", err))
            .finally(() => setCreating(false))
    }

    const onKeyDown = useDialogKeyboard(handleCreate)

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>New tract</h2>
            <input
                className={cls.DialogInput}
                value={name}
                onChange={e => setName(e.target.value)}
                onKeyDown={onKeyDown}
                placeholder="Tract name"
                autoFocus
                disabled={creating}
            />
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={CloseDialog} disabled={creating}>Cancel</Button>
                <Button variant="primary" onClick={handleCreate} disabled={creating || !name.trim()}>
                    {creating ? "…" : "Create"}
                </Button>
            </div>
        </div>
    )
}
