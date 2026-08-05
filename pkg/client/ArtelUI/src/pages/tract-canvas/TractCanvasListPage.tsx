import {useEffect} from "react"

import cls from "@/pages/tract-canvas/TractCanvasListPage.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useTracts} from "@/app/hooks/Tracts.ts"
import useUser from "@/hooks/user/User.ts"
import HeroSegment from "@/components/HeroSegment/HeroSegment.tsx"
import ContentSegment from "@/pages/tract-canvas/segments/ContentSegment/ContentSegment.tsx"
import NewTractButton from "@/pages/tract-canvas/components/NewTractButton/NewTractButton.tsx"
import NewTractDialog from "@/pages/tract-canvas/dialogs/NewTractDialog/NewTractDialog.tsx"

export default function TractCanvasListPage() {
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const {tracts, loading, fetch} = useTracts()
    const active = tracts.filter(t => t.enabled).length
    const paused = tracts.length - active

    // Unauthenticated → /init is handled at the router level (HomeLayout.tsx), not per-page.
    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetch()
        }
    }, [auth, fetch])

    return (
        <div className={cls.Root}>
            <HeroSegment
                eyebrow="Workspace · Automations"
                title="Tracts"
                subtitle={
                    <>
                        <b>{loading ? "…" : `${tracts.length} ${tracts.length === 1 ? "tract" : "tracts"}`}</b>
                        {tracts.length > 0 && <>{" · "}<span>{active} active, {paused} paused</span></>}
                    </>
                }
                action={<NewTractButton onClick={() => OpenDialog(<NewTractDialog/>)}/>}
                compact
            />
            <ContentSegment/>
        </div>
    )
}
