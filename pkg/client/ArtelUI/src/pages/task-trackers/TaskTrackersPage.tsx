import {useEffect} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/pages/task-trackers/TaskTrackersPage.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useTaskTrackers} from "@/app/hooks/TaskTrackers.ts"
import useUser from "@/hooks/user/User.ts"
import HeroSegment from "@/components/HeroSegment/HeroSegment.tsx"
import ContentSegment from "@/pages/task-trackers/components/ContentSegment/ContentSegment.tsx"
import AddTrackerDialog from "@/pages/task-trackers/components/AddTrackerDialog/AddTrackerDialog.tsx"

export default function TaskTrackersPage() {
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const {trackers, loading, fetch: fetchTrackers} = useTaskTrackers()

    // Unauthenticated → /init is handled at the router level (HomeLayout.tsx), not per-page.
    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetchTrackers()
        }
    }, [auth, fetchTrackers])

    return (
        <div className={cls.Root}>
            <HeroSegment
                eyebrow="Workspace"
                title="Task trackers"
                subtitle={
                    <>
                        <b>
                            {loading
                                ? "…"
                                : `${trackers.length} ${trackers.length === 1 ? "connection" : "connections"}`}
                        </b>
                        {" · "}<span>connected task tracking services</span>
                    </>
                }
                action={
                    <Button variant="primary" onClick={() => OpenDialog(<AddTrackerDialog/>)}>
                        <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor"
                             strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <line x1="12" y1="5" x2="12" y2="19"/>
                            <line x1="5" y1="12" x2="19" y2="12"/>
                        </svg>
                        Add connection
                    </Button>
                }
            />
            <ContentSegment/>
        </div>
    )
}
