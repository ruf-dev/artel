import {useEffect, useState} from "react"
import {ModalClose} from "@vervstack/chures"

import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import {fetchBoards, TrelloBoardLite} from "@/dialogs/SelectRoadmapRootDialog/processes/trelloBrowse.ts"
import cls from "@/dialogs/SelectRoadmapRootDialog/screens/PickBoardStep.module.css"

interface Props {
    trackerId: string
    onSelect: (board: TrelloBoardLite) => void
    onClose: () => void
}

export default function PickBoardStep({trackerId, onSelect, onClose}: Props) {
    const [boards, setBoards] = useState<TrelloBoardLite[]>([])
    const [loading, setLoading] = useState(true)
    const bakeError = useBakeError()

    useEffect(() => {
        setLoading(true)
        fetchBoards(trackerId)
            .then(setBoards)
            .catch(err => bakeError("Failed to load boards", err))
            .finally(() => setLoading(false))
    }, [trackerId])

    return (
        <div className={cls.PickBoardStepContainer} role="dialog" aria-modal="true" aria-labelledby="pickBoardTitle">
            <div className={cls.Head}>
                <h2 className={cls.Title} id="pickBoardTitle">View roadmap — choose a board</h2>
                <ModalClose onClick={onClose}/>
            </div>
            <p className={cls.Sub}>Pick the Trello board that holds the card you want to start the roadmap from.</p>

            <div className={cls.List}>
                {boards.map(b => (
                    <SelectOption key={b.id} label={b.name} selected={false} onSelect={() => onSelect(b)}/>
                ))}
                {!loading && boards.length === 0 && <p className={cls.Empty}>No boards found.</p>}
                {loading && <p className={cls.Empty}>Loading…</p>}
            </div>
        </div>
    )
}
