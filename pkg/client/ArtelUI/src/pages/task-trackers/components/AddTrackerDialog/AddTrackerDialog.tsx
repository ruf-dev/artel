import {useState} from "react"

import {TrelloBoardInfo} from "@/app/api/artel/task_trackers.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import ChooseTypeStep from "@/pages/task-trackers/components/AddTrackerDialog/screens/ChooseTypeStep.tsx"
import TrelloSetupStep from "@/pages/task-trackers/components/AddTrackerDialog/screens/TrelloSetupStep.tsx"
import DoneStep from "@/pages/task-trackers/components/AddTrackerDialog/screens/DoneStep.tsx"

type Step = "choose" | "trello" | "done"

export default function AddTrackerDialog() {
    const [step, setStep] = useState<Step>("choose")
    const [boards, setBoards] = useState<TrelloBoardInfo[]>([])

    const {CloseDialog} = useDialog()

    function handleSuccess(connectedBoards: TrelloBoardInfo[]) {
        setBoards(connectedBoards)
        setStep("done")
    }

    if (step === "choose") {
        return <ChooseTypeStep onChooseTrello={() => setStep("trello")} onClose={CloseDialog}/>
    }

    if (step === "trello") {
        return <TrelloSetupStep onSuccess={handleSuccess} onBack={() => setStep("choose")} onClose={CloseDialog}/>
    }

    return <DoneStep boards={boards} onClose={CloseDialog}/>
}
