import cls
    from "@/pages/tract-canvas/components/TractBlockPicker/components/OptionCell/components/OptText/OptText.module.css"

interface OptTextProps {
    name: string
    fn: string
}

export default function OptText({name, fn}: OptTextProps) {
    return (
        <div className={cls.OptTextContainer}>
            <div className={cls.OptName}>{name}</div>
            <div className={cls.OptFn}>{fn}</div>
        </div>
    )
}
