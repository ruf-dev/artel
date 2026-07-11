import {Button, CheckmarkIcon, ModalClose} from "@vervstack/chures"

import {TrelloBoardInfo} from "@/app/api/artel/task_trackers.pb.ts"
import SuccessBannerText
    from "@/pages/task-trackers/components/AddTrackerDialog/screens/components/SuccessBannerText/SuccessBannerText.tsx"
import cls from "@/pages/task-trackers/components/AddTrackerDialog/screens/DoneStep.module.css"

export default function DoneStep({boards, onClose}: {
    boards: TrelloBoardInfo[]
    onClose: () => void
}) {
    return (
        <div className={cls.DoneStepContainer}>
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
                    <Button variant="primary" onClick={onClose}>Done</Button>
                </div>
            </div>
        </div>
    )
}
