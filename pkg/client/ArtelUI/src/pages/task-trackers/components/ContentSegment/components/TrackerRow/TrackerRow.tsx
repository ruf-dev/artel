import {useState} from "react"
import {Button} from "@vervstack/chures"

import {TaskTrackerInfo} from "@/app/api/artel/task_trackers.pb.ts"
import cls from "@/pages/task-trackers/components/ContentSegment/components/TrackerRow/TrackerRow.module.css"

export default function TrackerRow({tracker, onRemove}: {
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
        <div className={cls.TrackerRowContainer}>
            <div className={cls.RowInfo}>
                <span className={cls.RowName}>{tracker.name}</span>
                <span className={cls.RowMeta}>{tracker.type} · connected {tracker.createdAt}</span>
            </div>
            <Button variant="danger" onClick={handleRemove} disabled={removing}>
                {removing ? "Removing…" : "Remove"}
            </Button>
        </div>
    )
}
