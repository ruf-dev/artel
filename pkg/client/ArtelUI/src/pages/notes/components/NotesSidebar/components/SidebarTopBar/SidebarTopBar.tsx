import {useMemo} from "react"
import {Dropdown} from "@vervstack/chures"

import NotesSearchBar from "@/pages/notes/components/NotesSidebar/components/NotesSearchBar/NotesSearchBar.tsx"
// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import PublicBadge from "@/pages/notes/components/NotesSidebar/components/SidebarTopBar/components/PublicBadge/PublicBadge.tsx"
import cls from "@/pages/notes/components/NotesSidebar/components/SidebarTopBar/SidebarTopBar.module.css"

interface VaultOption {
    id: string
    name: string
    isPublic?: boolean
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
    // chures' Dropdown has no per-option renderer (only a plain text label per row), so a
    // public vault can't get a real inline badge in the option list — suffix the label
    // instead. The selected-vault badge next to the trigger (below) covers the closed state.
    const options = useMemo(
        () => vaults.map(v => ({id: v.id, name: v.isPublic ? `${v.name} · Public` : v.name})),
        [vaults],
    )
    const selectedVault = vaults.find(v => v.id === vaultId)

    return (
        <div className={cls.SidebarTopBarContainer}>
            <div className={cls.SearchWrapper}>
                <NotesSearchBar treeView={treeView} onToggleTreeView={onToggleTreeView}/>
            </div>
            <div className={cls.VaultPickerWrapper}>
                {selectedVault?.isPublic && <PublicBadge/>}
                <Dropdown
                    className={cls.VaultDropdown}
                    placeholder="Select vault…"
                    options={options}
                    value={vaultId ? [vaultId] : []}
                    onChange={onVaultChange}
                />
            </div>
        </div>
    )
}
