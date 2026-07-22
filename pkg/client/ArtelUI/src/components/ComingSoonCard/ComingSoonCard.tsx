import {ReactNode} from "react"

import cls from "@/components/ComingSoonCard/ComingSoonCard.module.css"

interface ComingSoonCardProps {
    icon: ReactNode
    name: string
    hint: string
}

export default function ComingSoonCard({icon, name, hint}: ComingSoonCardProps) {
    return (
        <div
            className={cls.ComingSoonCardContainer}
            aria-disabled="true"
            data-tooltip-id="root-tooltip"
            data-tooltip-content={hint}
        >
            <div className={cls.CardHeader}>
                <div className={cls.CardIcon}>{icon}</div>
                <div className={cls.CardTitles}>
                    <div className={cls.CardName}>{name}</div>
                </div>
            </div>
            <div className={cls.CardFooter}>
                <span className={cls.ComingSoonBadge}>Coming soon</span>
            </div>
        </div>
    )
}
