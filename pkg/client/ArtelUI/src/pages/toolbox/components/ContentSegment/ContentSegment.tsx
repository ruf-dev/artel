import {useState} from "react"

import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import SearchInput from "@/components/SearchInput/SearchInput.tsx"
import MomCandidateRow from "@/pages/toolbox/components/ContentSegment/components/MomCandidateRow/MomCandidateRow.tsx"
import cls from "@/pages/toolbox/components/ContentSegment/ContentSegment.module.css"

export default function ContentSegment() {
    const {momCandidates, momCandidatesLoading} = useMcpKeys()
    const [query, setQuery] = useState("")

    if (momCandidatesLoading) {
        return (
            <div className={cls.ContentSegmentContainer}>
                <p className={cls.Empty}>Loading…</p>
            </div>
        )
    }

    const q = query.trim().toLowerCase()
    const filtered = q
        ? momCandidates.filter(c =>
            c.name?.toLowerCase().includes(q)
            || c.author?.toLowerCase().includes(q)
            || c.description?.toLowerCase().includes(q))
        : momCandidates

    return (
        <div className={cls.ContentSegmentContainer}>
            <SearchInput
                className={cls.SearchWrapper}
                placeholder="Search MCPs…"
                value={query}
                setValue={setQuery}
            />
            <div className={cls.List}>
                {filtered.map(c => (
                    <MomCandidateRow key={c.name} candidate={c}/>
                ))}
                {filtered.length === 0 && (
                    <p className={cls.Empty}>
                        {momCandidates.length === 0 ? "No tools available." : "No MCPs match your search."}
                    </p>
                )}
            </div>
        </div>
    )
}
