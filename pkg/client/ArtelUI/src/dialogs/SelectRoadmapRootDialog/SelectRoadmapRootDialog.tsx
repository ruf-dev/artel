import {useState} from "react"
import {useNavigate} from "react-router-dom"

import {useDialog} from "@/app/hooks/Dialog"
import {
    TrelloBoardLite, TrelloCardLite, TrelloListLite,
} from "@/dialogs/SelectRoadmapRootDialog/processes/trelloBrowse.ts"
import PickBoardStep from "@/dialogs/SelectRoadmapRootDialog/screens/PickBoardStep.tsx"
import PickListStep from "@/dialogs/SelectRoadmapRootDialog/screens/PickListStep.tsx"
import PickCardStep from "@/dialogs/SelectRoadmapRootDialog/screens/PickCardStep.tsx"

type Step = "board" | "list" | "card"

// Multi-screen picker: board -> list -> card. Confirming a card navigates straight to the
// roadmap route for it and closes the dialog — there's no separate "confirm" step because
// picking a card row IS the confirmation (see PickCardStep's onSelect wired to onConfirm).
export default function SelectRoadmapRootDialog({trackerId}: {trackerId: string}) {
    const [step, setStep] = useState<Step>("board")
    const [board, setBoard] = useState<TrelloBoardLite | null>(null)
    const [list, setList] = useState<TrelloListLite | null>(null)

    const {CloseDialog} = useDialog()
    const navigate = useNavigate()

    function handleConfirmCard(card: TrelloCardLite) {
        CloseDialog()
        navigate(`/task-trackers/${trackerId}/roadmap/${card.id}`)
    }

    if (step === "board") {
        return (
            <PickBoardStep
                trackerId={trackerId}
                onSelect={b => {
                    setBoard(b)
                    setStep("list")
                }}
                onClose={CloseDialog}
            />
        )
    }

    if (step === "list" && board) {
        return (
            <PickListStep
                trackerId={trackerId}
                board={board}
                onSelect={l => {
                    setList(l)
                    setStep("card")
                }}
                onBack={() => setStep("board")}
                onClose={CloseDialog}
            />
        )
    }

    if (step === "card" && board && list) {
        return (
            <PickCardStep
                trackerId={trackerId}
                board={board}
                list={list}
                onConfirm={handleConfirmCard}
                onBack={() => setStep("list")}
                onClose={CloseDialog}
            />
        )
    }

    return null
}
