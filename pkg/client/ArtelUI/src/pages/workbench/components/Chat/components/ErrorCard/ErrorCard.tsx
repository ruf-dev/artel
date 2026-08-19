import cls from "@/pages/workbench/components/Chat/components/ErrorCard/ErrorCard.module.css"

interface Props {
    text: string
}

export default function ErrorCard({text}: Props) {
    return (
        <div className={cls.ErrorCardContainer}>
            <span className={cls.Icon}>⨯</span>
            <p className={cls.Text}>{text}</p>
        </div>
    )
}
