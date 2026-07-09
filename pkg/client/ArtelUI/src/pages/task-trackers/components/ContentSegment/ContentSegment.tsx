import {useTaskTrackers} from "@/app/hooks/TaskTrackers.ts"
import TrackerRow from "@/pages/task-trackers/components/ContentSegment/components/TrackerRow/TrackerRow.tsx"
import cls from "@/pages/task-trackers/components/ContentSegment/ContentSegment.module.css"

export default function ContentSegment() {
    const {trackers, loading, remove} = useTaskTrackers()

    if (loading) {
        return (
            <div className={cls.ContentSegmentContainer}>
                <p className={cls.Empty}>Loading…</p>
            </div>
        )
    }

    return (
        <div className={cls.ContentSegmentContainer}>
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
