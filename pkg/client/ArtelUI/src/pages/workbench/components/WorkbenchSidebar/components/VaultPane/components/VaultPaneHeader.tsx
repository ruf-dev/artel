import cls from
    "@/pages/workbench/components/WorkbenchSidebar/components/VaultPane/components/VaultPaneHeader.module.css"

interface Props {
    vaultName: string
    noteCount: number
}

// The Vault pane's title row: the current vault's name (ellipsised when it can't
// fit) and a plain count of how many notes it holds.
export default function VaultPaneHeader({vaultName, noteCount}: Props) {
    return (
        <div className={cls.VaultPaneHeaderContainer}>
            <span className={cls.Name} title={vaultName}>{vaultName}</span>
            <span className={cls.Count}>{noteCount}</span>
        </div>
    )
}
