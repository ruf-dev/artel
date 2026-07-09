import {ModalClose} from "@vervstack/chures"

import TypeGrid from "@/pages/task-trackers/components/AddTrackerDialog/screens/components/TypeGrid/TypeGrid.tsx"
import cls from "@/pages/task-trackers/components/AddTrackerDialog/screens/ChooseTypeStep.module.css"

export default function ChooseTypeStep({onChooseTrello, onClose}: {
    onChooseTrello: () => void
    onClose: () => void
}) {
    return (
        <div className={cls.ChooseTypeStepContainer}>
            <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
                 aria-labelledby="chooseTrackerTitle">
                <div className={cls.ModalHead}>
                    <h2 className={cls.ModalTitle} id="chooseTrackerTitle">Add connection</h2>
                    <ModalClose onClick={onClose}/>
                </div>
                <p className={cls.ModalSub}>Choose a task tracking service to connect.</p>

                <TypeGrid onChooseTrello={onChooseTrello}/>
            </div>
        </div>
    )
}
