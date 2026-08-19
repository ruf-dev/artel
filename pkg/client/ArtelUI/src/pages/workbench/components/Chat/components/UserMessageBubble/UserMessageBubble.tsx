import cls from "@/pages/workbench/components/Chat/components/UserMessageBubble/UserMessageBubble.module.css"

interface Props {
    text: string
}

export default function UserMessageBubble({text}: Props) {
    return (
        <div className={cls.UserMessageBubbleContainer}>
            <p className={cls.Text}>{text}</p>
        </div>
    )
}
