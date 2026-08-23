import {motion} from "framer-motion"

import cls from "@/pages/workbench/components/Chat/components/UserMessageBubble/UserMessageBubble.module.css"

interface Props {
    text: string
}

export default function UserMessageBubble({text}: Props) {
    return (
        <motion.div
            className={cls.UserMessageBubbleContainer}
            initial={{opacity: 0, y: 28, scale: 0.9}}
            animate={{opacity: 1, y: 0, scale: 1}}
            exit={{opacity: 0, scale: 0.9}}
            transition={{type: "spring", stiffness: 420, damping: 32}}
        >
            <p className={cls.Text}>{text}</p>
        </motion.div>
    )
}
