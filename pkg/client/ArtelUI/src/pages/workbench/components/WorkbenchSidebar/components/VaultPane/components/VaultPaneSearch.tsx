import Input from "@/components/atoms/Input/Input.tsx"
import cls from
    "@/pages/workbench/components/WorkbenchSidebar/components/VaultPane/components/VaultPaneSearch.module.css"

interface Props {
    value: string
    onChange: (v: string) => void
}

// The Vault pane's filter box. Controlled by VaultPane, which owns the query
// string and feeds it to vaultPaneFilter. Mirrors NotesSearchBar's styling but
// without the list/tree toggle.
export default function VaultPaneSearch({value, onChange}: Props) {
    return (
        <div className={cls.VaultPaneSearchContainer}>
            <svg
                className={cls.SearchIcon} viewBox="0 0 16 16" width={16} height={16} fill="none"
                stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" aria-hidden="true"
            >
                <circle cx="6.5" cy="6.5" r="4.5"/>
                <path d="M10.5 10.5l3 3"/>
            </svg>
            <Input
                className={cls.SearchInputWrapper}
                inputClassName={cls.SearchInput}
                placeholder="Search files…"
                value={value}
                setValue={onChange}
            />
        </div>
    )
}
