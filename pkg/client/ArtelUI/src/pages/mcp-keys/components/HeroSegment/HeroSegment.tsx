import {Button} from "@vervstack/chures"

import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import cls from "@/pages/mcp-keys/components/HeroSegment/HeroSegment.module.css"

interface HeroSegmentProps {
    onCreateClick: () => void
}

export default function HeroSegment({onCreateClick}: HeroSegmentProps) {
    const {keys, loading} = useMcpKeys()

    return (
        <div className={cls.HeroSegmentContainer}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>MCP</div>
                <h1 className={cls.HeroTitle}>API Keys</h1>
                <p className={cls.HeroSub}>
                    <b>{loading ? "…" : `${keys.length} ${keys.length === 1 ? "key" : "keys"}`}</b>
                    {" · "}<span>bridge your MCP agents to Artel</span>
                </p>
            </div>
            <Button variant="primary" onClick={onCreateClick}>
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2"
                     strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                New key
            </Button>
        </div>
    )
}
