import {useState} from "react"

import FileTree from "@/components/FileTree/FileTree.tsx"
import {VaultPaneBundle} from "@/pages/workbench/processes/useWorkbenchSidebar.ts"
import VaultPaneHeader from
    "@/pages/workbench/components/WorkbenchSidebar/components/VaultPane/components/VaultPaneHeader.tsx"
import VaultPaneSearch from
    "@/pages/workbench/components/WorkbenchSidebar/components/VaultPane/components/VaultPaneSearch.tsx"
import VaultPaneFooter from
    "@/pages/workbench/components/WorkbenchSidebar/components/VaultPane/components/VaultPaneFooter.tsx"
import {filterVaultTree} from
    "@/pages/workbench/components/WorkbenchSidebar/components/VaultPane/vaultPaneFilter.ts"
import cls from "@/pages/workbench/components/WorkbenchSidebar/components/VaultPane/VaultPane.module.css"

// The Vault tab body: the current vault's folders/notes as a FileTree, filtered
// by a local search box. Clicking a note row toggles it in the chat's attachment
// set (attachedPaths / onToggleAttach come from useWorkbenchAttachments), and an
// attached row renders active. Structure mirrors HistoryPane.
export default function VaultPane(props: VaultPaneBundle) {
    const {vaultName, folders, notes, isLoading, attachedPaths, onToggleAttach} = props
    const [query, setQuery] = useState("")

    const tree = filterVaultTree(query, folders, notes)

    return (
        <div className={cls.VaultPaneContainer}>
            <VaultPaneHeader vaultName={vaultName} noteCount={notes.length}/>
            <VaultPaneSearch value={query} onChange={setQuery}/>
            {isLoading && <p className={cls.State}>Loading vault…</p>}
            {!isLoading && notes.length === 0 && (
                <p className={cls.State}>This vault has no notes.</p>
            )}
            {!isLoading && notes.length > 0 && (
                <div className={cls.ScrollArea}>
                    <FileTree
                        folders={tree.folders}
                        items={tree.notes}
                        isActive={path => attachedPaths.includes(path)}
                        onSelectItem={onToggleAttach}
                    />
                </div>
            )}
            <VaultPaneFooter attachedCount={attachedPaths.length}/>
        </div>
    )
}
