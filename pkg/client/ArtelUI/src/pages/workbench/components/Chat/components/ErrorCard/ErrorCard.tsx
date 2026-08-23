import {motion} from "framer-motion"

import cls from "@/pages/workbench/components/Chat/components/ErrorCard/ErrorCard.module.css"

interface Props {
    text: string
}

export default function ErrorCard({text}: Props) {
    return (
        <motion.div
            className={cls.ErrorCardContainer}
            initial={{opacity: 0, y: 14}}
            animate={{opacity: 1, y: 0}}
            exit={{opacity: 0}}
            transition={{duration: 0.22, ease: "easeOut"}}
        >
            <span className={cls.Icon}>⨯</span>
            <p className={cls.Text}>{text}</p>
        </motion.div>
    )
}
