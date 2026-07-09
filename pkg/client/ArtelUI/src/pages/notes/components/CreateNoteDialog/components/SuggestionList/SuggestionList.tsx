import { Button } from "@vervstack/chures"

import { cn } from "@/app/utils/cn"
import cls from "@/pages/notes/components/CreateNoteDialog/components/SuggestionList/SuggestionList.module.css"

interface SuggestionListProps {
    suggestions: string[]
    highlightedIdx: number
    onSelect: (s: string) => void
}

export default function SuggestionList({ suggestions, highlightedIdx, onSelect }: SuggestionListProps) {
    return (
        <div className={cls.SuggestionListContainer}>
            {suggestions.map((s, i) => (
                <Button
                    key={s}
                    variant="ghost"
                    className={cn(cls.SuggestionItem, i === highlightedIdx && cls.SuggestionItemHighlighted)}
                    onMouseDown={e => { e.preventDefault(); onSelect(s) }}
                >
                    {s}
                </Button>
            ))}
        </div>
    )
}
