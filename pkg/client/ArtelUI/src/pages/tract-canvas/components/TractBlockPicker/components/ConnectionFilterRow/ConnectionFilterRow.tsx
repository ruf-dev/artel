import {Button} from "@vervstack/chures"

import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {cn} from "@/app/utils/cn.ts"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"
import cls from "@/pages/tract-canvas/components/TractBlockPicker/TractBlockPicker.module.css"

interface ConnectionFilterRowProps {
    candidates: MomCandidate[]
    connFilter: string | null
    onSelect: (name: string | null) => void
}

export default function ConnectionFilterRow({candidates, connFilter, onSelect}: ConnectionFilterRowProps) {
    return (
        <div className={cls.FilterRow}>
            <Button
                variant="ghost"
                className={cn(cls.FilterPill, !connFilter && cls.FilterPillActive)}
                onClick={() => onSelect(null)}
            >
                All
            </Button>
            {candidates.map(c => (
                <Button
                    variant="ghost"
                    key={c.name}
                    className={cn(cls.FilterPill, connFilter === c.name && cls.FilterPillActive)}
                    onClick={() => onSelect(c.name ?? null)}
                >
                    <span className={cls.FilterPillIcon}>
                        <ProviderIcon provider={c.connections?.[0]?.provider}/></span>
                    {c.name}
                </Button>
            ))}
        </div>
    )
}
