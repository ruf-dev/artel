import cls from "@/components/McpKeyCard/CardMeta.module.css"

interface Props {
    lastAccessedAt?: string
}

export default function CardMeta({lastAccessedAt}: Props) {
    return (
        <div className={cls.CardMetaContainer}>
            Last accessed: {formatDate(lastAccessedAt)}
        </div>
    )
}

function formatDate(iso: string | undefined): string {
    if (!iso) return "Never"
    const d = new Date(iso)
    if (isNaN(d.getTime())) return "Never"
    return d.toLocaleDateString(undefined, {year: "numeric", month: "short", day: "numeric"})
}
