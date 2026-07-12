import { Button } from "@vervstack/chures"

import { cn } from "@/app/utils/cn.ts"
import Input from "@/components/atoms/Input/Input.tsx"
import SearchIcon from "@/pages/notes/components/icons/SearchIcon.tsx"
import CloseIcon from "@/pages/notes/components/icons/CloseIcon.tsx"
import ListIcon from "@/pages/notes/components/NotesSidebar/components/icons/ListIcon.tsx"
import TreeIcon from "@/pages/notes/components/NotesSidebar/components/icons/TreeIcon.tsx"
import { useNotesSearchQuery } from "@/pages/notes/components/NotesSidebar/processes/useNotesSearchQuery.ts"
import cls from "@/pages/notes/components/NotesSidebar/components/NotesSearchBar/NotesSearchBar.module.css"

interface NotesSearchBarProps {
    treeView: boolean
    onToggleTreeView: () => void
}

export default function NotesSearchBar({ treeView, onToggleTreeView }: NotesSearchBarProps) {
    const [searchQuery, setSearchQuery] = useNotesSearchQuery()

    return (
        <div className={cls.NotesSearchBarContainer}>
            <SearchIcon className={cls.SearchIcon}/>
            <Input
                className={cls.SearchInputWrapper}
                inputClassName={cls.SearchInput}
                placeholder="Search notes…"
                value={searchQuery}
                setValue={setSearchQuery}
            />
            {searchQuery && (
                <Button
                    variant="ghost"
                    className={cn(cls.TreeToggleBtn, treeView && cls.TreeToggleBtnActive)}
                    onClick={onToggleTreeView}
                    aria-label={treeView ? "Switch to list view" : "Switch to tree view"}
                    title={treeView ? "Switch to list view" : "Switch to tree view"}
                >
                    {treeView ? <TreeIcon/> : <ListIcon/>}
                </Button>
            )}
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
