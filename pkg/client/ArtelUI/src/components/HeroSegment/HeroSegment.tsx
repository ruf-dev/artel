import {ReactNode} from "react"

import {cn} from "@/app/utils/cn.ts"
import cls from "@/components/HeroSegment/HeroSegment.module.css"

interface HeroSegmentProps {
    eyebrow: ReactNode
    title: string
    subtitle: ReactNode
    action?: ReactNode
    compact?: boolean
}

export default function HeroSegment({eyebrow, title, subtitle, action, compact}: HeroSegmentProps) {
    return (
        <div className={cn(cls.HeroSegmentContainer, compact && cls.Compact)}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>{eyebrow}</div>
                <h1 className={cls.HeroTitle}>{title}</h1>
                <p className={cls.HeroSub}>{subtitle}</p>
            </div>
            {action}
        </div>
    )
}
