import {useNavigate} from "react-router-dom"

import cls from "@/pages/tract-canvas/segments/ContentSegment/ContentSegment.module.css"

import {useDialog} from "@/app/hooks/Dialog"
import {useTracts} from "@/app/hooks/Tracts.ts"

import TractCanvasCard from "@/pages/tract-canvas/widgets/TractCanvasCard/TractCanvasCard.tsx"
import NewTractButton from "@/pages/tract-canvas/components/NewTractButton/NewTractButton.tsx"
import NewTractDialog from "@/pages/tract-canvas/dialogs/NewTractDialog/NewTractDialog.tsx"

export default function ContentSegment() {
    const {tracts, loading} = useTracts()
    const {OpenDialog} = useDialog()
    const navigate = useNavigate()

    if (loading) {
        return <div className={cls.ContentSegmentContainer}><p className={cls.Empty}>Loading…</p></div>
    }

    if (tracts.length === 0) {
        return (
            <div className={cls.ContentSegmentContainer}>
                <div className={cls.EmptyState}>
                    <h2 className={cls.EmptyStateTitle}>No tracts yet</h2>
                    <p className={cls.EmptyStateSub}>
                        Tracts connect a trigger — a webhook or a manual fire — to a tree of
                        actions and conditions. Create your first one to get started.
                    </p>
                    <NewTractButton onClick={() => OpenDialog(<NewTractDialog/>)}/>
                </div>
            </div>
        )
    }

    return (
        <div className={cls.ContentSegmentContainer}>
            <div className={cls.List}>
                {tracts.map(t => (
                    <TractCanvasCard
                        key={t.uuid}
                        tract={t}
                        onClick={() => navigate(`/tracts/${t.uuid}`)}/>
                ))}
            </div>
        </div>
    )
}
