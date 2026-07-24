import {Dropdown} from "@vervstack/chures"

import NotesSearchBar from "@/pages/notes/components/NotesSidebar/components/NotesSearchBar/NotesSearchBar.tsx"
import cls from "@/pages/notes/components/NotesSidebar/components/SidebarTopBar/SidebarTopBar.module.css"

interface VaultOption {
    id: string
    name: string
}

interface SidebarTopBarProps {
    vaults: VaultOption[]
    vaultId: string | null
    onVaultChange: (value: string[]) => void
    treeView: boolean
    onToggleTreeView: () => void
}

export default function SidebarTopBar(
    {vaults, vaultId, onVaultChange, treeView, onToggleTreeView}: SidebarTopBarProps
) {
    return (
        <div className={cls.SidebarTopBarContainer}>
            <div className={cls.SearchWrapper}>
                <NotesSearchBar treeView={treeView} onToggleTreeView={onToggleTreeView}/>
            </div>
            <div className={cls.VaultPickerWrapper}>
                <Dropdown
                    placeholder="Select vault…"
                    options={vaults}
                    value={vaultId ? [vaultId] : []}
                    onChange={onVaultChange}
                />
            </div>
        </div>
    )
}
