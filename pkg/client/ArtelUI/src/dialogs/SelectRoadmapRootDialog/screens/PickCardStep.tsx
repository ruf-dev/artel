import {useEffect, useState} from "react"
import {Button, ModalClose} from "@vervstack/chures"

import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import SearchInput from "@/components/SearchInput/SearchInput.tsx"
import {
    fetchCards, TrelloBoardLite, TrelloCardLite, TrelloListLite,
} from "@/dialogs/SelectRoadmapRootDialog/processes/trelloBrowse.ts"
import cls from "@/dialogs/SelectRoadmapRootDialog/screens/PickCardStep.module.css"

interface Props {
    trackerId: string
    board: TrelloBoardLite
    list: TrelloListLite
    onConfirm: (card: TrelloCardLite) => void
    onBack: () => void
    onClose: () => void
}

export default function PickCardStep({trackerId, board, list, onConfirm, onBack, onClose}: Props) {
    const [cards, setCards] = useState<TrelloCardLite[]>([])
    const [loading, setLoading] = useState(true)
    const [query, setQuery] = useState("")
    const bakeError = useBakeError()

    useEffect(() => {
        setLoading(true)
        fetchCards(trackerId, board.id)
            .then(setCards)
            .catch(err => bakeError("Failed to load cards", err))
            .finally(() => setLoading(false))
    }, [trackerId, board.id])

    const listCards = cards.filter(c => c.idList === list.id)
    const query_ = query.trim().toLowerCase()
    const visibleCards = query_ ? listCards.filter(c => c.name.toLowerCase().includes(query_)) : listCards

    return (
        <div className={cls.PickCardStepContainer} role="dialog" aria-modal="true" aria-labelledby="pickCardTitle">
            <div className={cls.Head}>
                <h2 className={cls.Title} id="pickCardTitle">Choose a card on "{list.name}"</h2>
                <ModalClose onClick={onClose}/>
            </div>

            <SearchInput value={query} setValue={setQuery} placeholder="Search cards…" className={cls.Search}/>

            <div className={cls.List}>
                {visibleCards.map(c => (
                    <SelectOption key={c.id} label={c.name} selected={false} onSelect={() => onConfirm(c)}/>
                ))}
                {!loading && visibleCards.length === 0 && <p className={cls.Empty}>No cards found.</p>}
                {loading && <p className={cls.Empty}>Loading…</p>}
            </div>

            <div className={cls.Actions}>
                <Button variant="secondary" onClick={onBack}>Back</Button>
            </div>
        </div>
    )
}
