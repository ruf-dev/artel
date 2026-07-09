import cls from "@/components/TractRunTimeline/TractRunTimeline.module.css"

interface Props {
    label: string
    value: unknown
}

export default function JsonBlock({label, value}: Props) {
    return (
        <div className={cls.JsonBlock}>
            <span className={cls.JsonLabel}>{label}</span>
            <pre className={cls.Json}>{JSON.stringify(value, null, 2)}</pre>
        </div>
    )
}
