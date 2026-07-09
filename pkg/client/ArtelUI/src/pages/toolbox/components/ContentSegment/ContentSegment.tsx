import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import MomCandidateRow from "@/pages/toolbox/components/ContentSegment/components/MomCandidateRow/MomCandidateRow.tsx"
import cls from "@/pages/toolbox/components/ContentSegment/ContentSegment.module.css"

export default function ContentSegment() {
    const {momCandidates, momCandidatesLoading} = useMcpKeys()

    if (momCandidatesLoading) {
        return (
            <div className={cls.ContentSegmentContainer}>
                <p className={cls.Empty}>Loading…</p>
            </div>
        )
    }

    return (
        <div className={cls.ContentSegmentContainer}>
            <div className={cls.List}>
                {momCandidates.map(c => (
                    <MomCandidateRow key={c.name} candidate={c}/>
                ))}
                {momCandidates.length === 0 && (
                    <p className={cls.Empty}>No tools available.</p>
                )}
            </div>
        </div>
    )
}
