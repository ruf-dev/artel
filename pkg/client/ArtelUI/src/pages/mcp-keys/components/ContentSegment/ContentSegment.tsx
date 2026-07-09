import {useDialog} from "@/app/hooks/Dialog"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import ManageKeyDialog from "@/dialogs/ManageKeyDialog/ManageKeyDialog.tsx"
import McpKeyCard from "@/widgets/McpKeyCard/McpKeyCard.tsx"
import cls from "@/pages/mcp-keys/components/ContentSegment/ContentSegment.module.css"

export default function ContentSegment() {
    const {keys, loading} = useMcpKeys()
    const {OpenDialog} = useDialog()

    // TODO add loader from chures
    if (loading) {
        return (
            <div className={cls.ContentSegmentContainer}>
                <p className={cls.Empty}>Loading…</p>
            </div>
        )
    }

    return (
        <div className={cls.ContentSegmentContainer}>
            <div className={cls.List}>
                {keys.map(key => (
                    <McpKeyCard
                        key={key.id}
                        mcpKey={key}
                        onManage={() => OpenDialog(<ManageKeyDialog mcpKey={key}/>)}
                    />
                ))}
                {keys.length === 0 && (
                    <p className={cls.Empty}>No API keys yet. Create one to get started.</p>
                )}
            </div>
        </div>
    )
}
