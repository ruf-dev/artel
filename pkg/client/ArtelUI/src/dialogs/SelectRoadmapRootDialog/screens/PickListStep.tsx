import {useEffect, useState} from "react"
import {Button, ModalClose} from "@vervstack/chures"

import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import {fetchLists, TrelloBoardLite, TrelloListLite} from "@/dialogs/SelectRoadmapRootDialog/processes/trelloBrowse.ts"
import cls from "@/dialogs/SelectRoadmapRootDialog/screens/PickListStep.module.css"

interface Props {
    trackerId: string
    board: TrelloBoardLite
    onSelect: (list: TrelloListLite) => void
    onBack: () => void
    onClose: () => void
}

export default function PickListStep({trackerId, board, onSelect, onBack, onClose}: Props) {
    const [lists, setLists] = useState<TrelloListLite[]>([])
    const [loading, setLoading] = useState(true)
    const bakeError = useBakeError()

    useEffect(() => {
        setLoading(true)
        fetchLists(trackerId, board.id)
            .then(setLists)
            .catch(err => bakeError("Failed to load lists", err))
            .finally(() => setLoading(false))
    }, [trackerId, board.id])

    return (
        <div className={cls.PickListStepContainer} role="dialog" aria-modal="true" aria-labelledby="pickListTitle">
            <div className={cls.Head}>
                <h2 className={cls.Title} id="pickListTitle">Choose a list on "{board.name}"</h2>
                <ModalClose onClick={onClose}/>
            </div>

            <div className={cls.List}>
                {lists.map(l => (
                    <SelectOption key={l.id} label={l.name} selected={false} onSelect={() => onSelect(l)}/>
                ))}
                {!loading && lists.length === 0 && <p className={cls.Empty}>No lists found.</p>}
                {loading && <p className={cls.Empty}>Loading…</p>}
            </div>

            <div className={cls.Actions}>
                <Button variant="secondary" onClick={onBack}>Back</Button>
            </div>
        </div>
    )
}
