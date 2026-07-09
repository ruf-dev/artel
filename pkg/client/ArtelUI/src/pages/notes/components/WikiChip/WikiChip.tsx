import { useId } from "react"

import cls from "@/pages/notes/components/WikiChip/WikiChip.module.css"
import ArtelMark from "@/pages/notes/components/icons/ArtelMark.tsx"

interface WikiChipProps {
    name: string
    onClick?: () => void
}

export default function WikiChip({ name, onClick }: WikiChipProps) {
    const maskId = useId()

    return (
        <span className={cls.WikiChipContainer} onClick={e => { e.stopPropagation(); onClick?.() }}>
            <ArtelMark size={11} maskId={maskId} />
            {name}
        </span>
    )
}
