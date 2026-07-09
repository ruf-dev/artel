import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import ToolsDialog from "@/pages/toolbox/components/ToolsDialog/ToolsDialog.tsx"
import cls from "@/pages/toolbox/components/ContentSegment/components/MomCandidateRow/MomCandidateRow.module.css"

export default function MomCandidateRow({candidate}: { candidate: MomCandidate }) {
    const {OpenDialog, SetClosable} = useDialog()
    const connCount = candidate.connections?.length ?? 0
    const toolCount = candidate.tools?.length ?? 0

    function onToolboxClick() {
        OpenDialog(<ToolsDialog candidate={candidate}/>)
        SetClosable(true)
    }

    return (
        <div className={cls.MomCandidateRowContainer} onClick={onToolboxClick} role="button" tabIndex={0}>
            <div className={cls.RowHeader}>
                <span className={cls.RowName}>{candidate.name}</span>
                {toolCount > 0 && (
                    <span className={cls.Chip}>{toolCount} {toolCount === 1 ? "tool" : "tools"}</span>
                )}
                {connCount > 0 && (
                    <span className={cls.Chip}>{connCount} {connCount === 1 ? "connection" : "connections"}</span>
                )}
            </div>
            <div className={cls.RowAuthor}>{candidate.author}</div>
            <div className={cls.RowDesc}>{candidate.description}</div>
        </div>
    )
}
