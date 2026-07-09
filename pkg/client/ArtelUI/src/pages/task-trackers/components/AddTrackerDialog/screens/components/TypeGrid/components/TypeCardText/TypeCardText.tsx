import cls from "@/pages/task-trackers/components/AddTrackerDialog/screens/components/TypeGrid/components/TypeCardText/TypeCardText.module.css"

export default function TypeCardText({label, desc}: {label: string; desc: string}) {
    return (
        <div className={cls.TypeCardTextContainer}>
            <div className={cls.TypeCardLabel}>{label}</div>
            <div className={cls.TypeCardDesc}>{desc}</div>
        </div>
    )
}
