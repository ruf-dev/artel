import cls from "@/components/McpKeyCard/CardHeader.module.css"

interface Props {
    name?: string
    keyPreview?: string
}

export default function CardHeader({name, keyPreview}: Props) {
    return (
        <div className={cls.CardHeaderContainer}>
            <span className={cls.CardName}>{name}</span>
            <span className={cls.CardPreview}>{keyPreview}…</span>
        </div>
    )
}
