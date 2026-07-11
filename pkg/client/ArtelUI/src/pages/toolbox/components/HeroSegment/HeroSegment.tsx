import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import cls from "@/pages/toolbox/components/HeroSegment/HeroSegment.module.css"

export default function HeroSegment() {
    const {momCandidates, momCandidatesLoading} = useMcpKeys()

    return (
        <div className={cls.HeroSegmentContainer}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>MCP</div>
                <h1 className={cls.HeroTitle}>Toolbox</h1>
                <p className={cls.HeroSub}>
                    <b>
                        {momCandidatesLoading
                            ? "…"
                            : `${momCandidates.length} ${momCandidates.length === 1 ? "tool" : "tools"}`}
                    </b>
                    {" · "}<span>available in this installation</span>
                </p>
            </div>
        </div>
    )
}
