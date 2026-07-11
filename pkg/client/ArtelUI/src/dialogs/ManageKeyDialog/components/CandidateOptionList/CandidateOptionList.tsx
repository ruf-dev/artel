import {McpConnectorInfo, MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import MomCandidateCard
    from "@/dialogs/ManageKeyDialog/components/CandidateOptionList/components/MomCandidateCard/MomCandidateCard.tsx"
import cls from "@/dialogs/ManageKeyDialog/components/CandidateOptionList/CandidateOptionList.module.css"

interface CandidateOptionListProps {
    connectors: McpConnectorInfo[]
    onSelect: (c: MomCandidate) => void
}

export default function CandidateOptionList({connectors, onSelect}: CandidateOptionListProps) {
    const {momCandidates: candidates} = useMcpKeys()
    return (
        <div className={cls.OptionList}>
            {candidates.map(c => (
                <MomCandidateCard
                    key={c.name}
                    candidate={c}
                    connected={connectors.some(con => con.mcpName === c.name)}
                    onSelect={() => onSelect(c)}
                />
            ))}
            {candidates.length === 0 && <p className={cls.Empty}>No services available yet.</p>}
        </div>
    )
}
