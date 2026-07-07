import cls from "@/pages/tract-canvas/components/Section/Section.module.css"

interface Props {
    title: string
    children: React.ReactNode
}

export default function Section({title, children}: Props) {
    return (
        <div className={cls.SectionContainer}>
            <div className={cls.SectionTitle}>{title}</div>
            {children}
        </div>
    )
}
