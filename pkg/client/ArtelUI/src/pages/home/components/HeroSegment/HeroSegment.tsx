import {Button} from "@vervstack/chures"

import {useVaults} from "@/app/hooks/Vaults.ts"
import cls from "@/pages/home/components/HeroSegment/HeroSegment.module.css"

interface HeroSegmentProps {
    onCreateClick: () => void
}

export default function HeroSegment({onCreateClick}: HeroSegmentProps) {
    const {isLoading, vaults} = useVaults()

    return (
        <div className={cls.HeroSegmentContainer}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>Workspace</div>
                <h1 className={cls.HeroTitle}>Your vaults</h1>
                <p className={cls.HeroSub}>
                    <b>{isLoading ? "…" : `${vaults.length} ${vaults.length === 1 ? "vault" : "vaults"}`}</b>
                    {" · "}<span>all systems operational</span>
                </p>
            </div>
            <Button variant="primary" onClick={onCreateClick}>
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2"
                     strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                New vault
            </Button>
        </div>
    )
}
