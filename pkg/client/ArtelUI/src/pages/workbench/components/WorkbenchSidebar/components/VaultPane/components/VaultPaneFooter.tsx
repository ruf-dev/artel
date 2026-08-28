import cls from
    "@/pages/workbench/components/WorkbenchSidebar/components/VaultPane/components/VaultPaneFooter.module.css"

interface Props {
    attachedCount: number
}

// Pinned strip under the file tree that echoes how many vault files are currently
// attached to the chat as context chips.
export default function VaultPaneFooter({attachedCount}: Props) {
    const label = attachedCount === 0 ? "Nothing attached" : `${attachedCount} attached`

    return <div className={cls.VaultPaneFooterContainer}>{label}</div>
}
