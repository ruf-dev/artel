import {Button} from "@vervstack/chures"

import {useTaskTrackers} from "@/app/hooks/TaskTrackers.ts"
import cls from "@/pages/task-trackers/components/HeroSegment/HeroSegment.module.css"

export default function HeroSegment({onAddClick}: { onAddClick: () => void }) {
    const {trackers, loading} = useTaskTrackers()

    return (
        <div className={cls.HeroSegmentContainer}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>Workspace</div>
                <h1 className={cls.HeroTitle}>Task trackers</h1>
                <p className={cls.HeroSub}>
                    <b>{loading ? "…" : `${trackers.length} ${trackers.length === 1 ? "connection" : "connections"}`}</b>
                    {" · "}<span>connected task tracking services</span>
                </p>
            </div>
            <Button variant="primary" onClick={onAddClick}>
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2"
                     strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Add connection
            </Button>
        </div>
    )
}
