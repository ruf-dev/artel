import {Button} from "@vervstack/chures"

import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import cls from "@/dialogs/ManageKeyDialog/components/CandidateOptionList/components/MomCandidateCard/MomCandidateCard.module.css"

interface MomCandidateCardProps {
    candidate: MomCandidate
    connected: boolean
    onSelect: () => void
}

export default function MomCandidateCard({candidate, connected, onSelect}: MomCandidateCardProps) {
    const count = candidate.connections?.length ?? 0
    return (
        <Button
            variant="ghost"
            className={cls.EntityCard}
            onClick={onSelect}
            disabled={connected}
            title={connected ? "Already connected" : undefined}
        >
            <div className={cls.MomCardHeader}>
                <span className={cls.EntityCardTitle}>{candidate.name}</span>
                <span className={count > 0 ? cls.Chip : `${cls.Chip} ${cls.ChipMuted}`}>
                    {connected
                        ? "Already connected"
                        : count > 0 ? `${count} connection${count > 1 ? "s" : ""}` : "No connection"}
                </span>
            </div>
            {candidate.description && <span className={cls.EntityCardDesc}>{candidate.description}</span>}
        </Button>
    )
}
