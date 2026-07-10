import { Button } from "@vervstack/chures"

import Input from "@/components/atoms/Input/Input.tsx"
import SearchIcon from "@/pages/notes/components/icons/SearchIcon.tsx"
import CloseIcon from "@/pages/notes/components/icons/CloseIcon.tsx"
import { useNotesSearchQuery } from "@/pages/notes/components/NotesSidebar/processes/useNotesSearchQuery.ts"
import cls from "@/pages/notes/components/NotesSidebar/components/NotesSearchBar/NotesSearchBar.module.css"

export default function NotesSearchBar() {
    const [searchQuery, setSearchQuery] = useNotesSearchQuery()

    return (
        <div className={cls.NotesSearchBarContainer}>
            <SearchIcon className={cls.SearchIcon}/>
            <Input
                className={cls.SearchInput}
                placeholder="Search notes…"
                value={searchQuery}
                setValue={setSearchQuery}
            />
            {searchQuery && (
                <Button
                    variant="ghost"
                    className={cls.ClearSearchBtn}
                    onClick={() => setSearchQuery("")}
                    aria-label="Clear search"
                >
                    <CloseIcon/>
                </Button>
            )}
        </div>
    )
}
