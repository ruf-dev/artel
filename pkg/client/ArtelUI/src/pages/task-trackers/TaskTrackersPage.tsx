import {useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/task-trackers/TaskTrackersPage.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useTaskTrackers} from "@/app/hooks/TaskTrackers.ts"
import useUser from "@/hooks/user/User.ts"
import HeroSegment from "@/pages/task-trackers/components/HeroSegment/HeroSegment.tsx"
import ContentSegment from "@/pages/task-trackers/components/ContentSegment/ContentSegment.tsx"
import AddTrackerDialog from "@/pages/task-trackers/components/AddTrackerDialog/AddTrackerDialog.tsx"

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
