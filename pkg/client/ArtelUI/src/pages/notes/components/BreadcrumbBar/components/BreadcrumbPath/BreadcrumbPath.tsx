import cls from "@/pages/notes/components/BreadcrumbBar/components/BreadcrumbPath/BreadcrumbPath.module.css"

interface BreadcrumbPathProps {
    path: string
}

export default function BreadcrumbPath({ path }: BreadcrumbPathProps) {
    const segments = path.split('/').filter(Boolean)
    return (
        <div className={cls.BreadcrumbPathContainer}>
            {segments.map((seg, i) => (
                <span key={i} className={cls.SegmentWrapper}>
                    {i > 0 && <span className={cls.Separator}>›</span>}
                    <span className={i === segments.length - 1 ? cls.SegmentActive : cls.Segment}>
                        {seg}
                    </span>
                </span>
            ))}
        </div>
    )
}
