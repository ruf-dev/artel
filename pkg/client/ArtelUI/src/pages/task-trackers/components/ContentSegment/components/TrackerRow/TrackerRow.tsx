import {useState} from "react"
import {Button} from "@vervstack/chures"

import {TaskTrackerInfo} from "@/app/api/artel/task_trackers.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import SelectRoadmapRootDialog from "@/dialogs/SelectRoadmapRootDialog/SelectRoadmapRootDialog.tsx"
import cls from "@/pages/task-trackers/components/ContentSegment/components/TrackerRow/TrackerRow.module.css"

export default function TrackerRow({tracker, onRemove}: {
    tracker: TaskTrackerInfo
    onRemove: (id: string) => Promise<void>
}) {
    const [removing, setRemoving] = useState(false)
    const {OpenDialog} = useDialog()

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
            <div className={cls.RowActions}>
                {tracker.type === "trello" && (
                    <Button
                        variant="secondary"
                        onClick={() => tracker.id && OpenDialog(<SelectRoadmapRootDialog trackerId={tracker.id}/>)}
                    >
                        View roadmap
                    </Button>
                )}
                <Button variant="danger" onClick={handleRemove} disabled={removing}>
                    {removing ? "Removing…" : "Remove"}
                </Button>
            </div>
        </div>
    )
}
